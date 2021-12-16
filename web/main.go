package main

import (
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"red/internal/driver"
	"red/internal/models"
	"time"

	"github.com/alexedwards/scs/v2"
)

const version = "1.0.0"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	secretKey string
	frontend  string
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
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

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	cfg.db.dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false", dbUser, dbPassword, dbHost, dbPort, dbName)
	//flag.StringVar(&cfg.mysql.dsn, "dsn", "trevor:secret@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	flag.Parse()

	cfg.secretKey = os.Getenv("Secret")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//conn, err := mysql.OpenDb(cfg.mysql.dsn)
	conn, err := driver.GetDBClient(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	sqlDB, _ := conn.DB()

	// 最多閒置數量
	sqlDB.SetMaxIdleConns(10)
	// 最多連接數量
	sqlDB.SetMaxOpenConns(10)
	// 等待醉酒時間
	sqlDB.SetConnMaxIdleTime(time.Hour)

	// set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(sqlDB)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Session:       session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
