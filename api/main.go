package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
	"log"
	"net/http"
	"os"
	"red/domain"
	"red/logger_logrus"
	"red/mongo"
	"red/redis"
	"red/service"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}

	frontend  string
	secretKey string
}

type application struct {
	config config
	log    logger_logrus.Logger
	//errorLog *log.Logger
	ch    CustomerHandler
	ah    AccountHandler
	uh    UserHandlers
	mongo mongo.Collection
	redis redis.Database
	mail  *nsq.Producer
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.log.Infof("Starting Back end server in %s mode on port %d", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	log.Println("Api Starting...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	var cfg config

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	cfg.db.dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false", dbUser, dbPassword, dbHost, dbPort, dbName)

	cfg.secretKey = os.Getenv("Secret")

	redisHost := os.Getenv("REDIS_HOST")

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	flag.Parse()

	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbClient, err := domain.GetDBClient(cfg.db.dsn)
	if err != nil {
		log.Println("failed to connect mysql " + err.Error())
	}

	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	userRepositoryDb := domain.NewUserRepositoryDb(dbClient)

	//建立各個Handlers
	//ch := CustomerHandler{service.NewCustomerService(customerRepositoryDb)}
	ch := CustomerHandler{service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}
	uh := UserHandlers{service.NewUserService(userRepositoryDb)}

	// connect to mongodb
	mongodbUser := os.Getenv("MONGODB_USER")
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")
	mongodbHost := os.Getenv("MONGODB_HOST")

	mongodb := mongo.New(mongodbHost, mongodbUser, mongodbPassword)

	//建立redis
	redisCli, err := redis.GetClient(redisHost)
	if err != nil {
		log.Println("failed to connect redis " + err.Error())
	}
	redisDB := redis.New(redisCli)

	nsqHost := os.Getenv("NSQ_HOST")
	nsqAddr := fmt.Sprintf("%s:4150", nsqHost)
	nsqConfig := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqAddr, nsqConfig)
	if err != nil {
		log.Println("connection to nsq failed: " + err.Error())
	}
	defer producer.Stop()

	err = producer.Publish("mailer", []byte("start"))
	if err != nil {
		log.Println(err.Error())
	}

	logrus := logger_logrus.New("./log", "gin.log")

	app := &application{
		config: cfg,
		log:    logrus,
		//errorLog: errorLog,
		ch:    ch,
		ah:    ah,
		uh:    uh,
		redis: redisDB,
		mongo: mongodb,
		mail:  producer,
	}

	err = app.serve()
	if err != nil {
		app.log.Debug(err.Error())
		log.Fatal(err)
	}

}
