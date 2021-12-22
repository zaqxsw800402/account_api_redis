package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type smtp struct {
	host     string
	port     int
	username string
	password string
}

type application struct {
	smtp      smtp
	secretKey string
	frontend  string
}

func main() {
	log.Println("Mailer Starting...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var s smtp
	s.host = "smtp.mailtrap.io"
	s.port = 587
	s.username = os.Getenv("Username")
	s.password = os.Getenv("Password")

	secretKey := os.Getenv("Secret")
	frontend := os.Getenv("Frontend")

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("mailer", "channel", config)
	if err != nil {
		log.Fatal(err)
	}
	app := &application{
		smtp:      s,
		secretKey: secretKey,
		frontend:  frontend,
	}

	//新增一個handler來處理收到訊息時動作
	consumer.AddHandler(app)

	//nsqHost := os.Getenv("NSQ_HOST")
	//nsqHost := os.Getenv("NSQ_HOST")

	//err = consumer.ConnectToNSQLookupd("nsqd:4161")
	err = consumer.ConnectToNSQLookupd("nsqlookupd:4161")

	//err = consumer.ConnectToNSQD(nsqHost)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}

func (app *application) HandleMessage(m *nsq.Message) error {
	var op struct {
		Op string
	}

	err := json.Unmarshal(m.Body, &op)
	if err != nil {
		log.Println(err)
	}

	switch op.Op {
	case "forgotPassword":
		app.resetPassword(m)
	case "newUser":
		app.newUser(m)
	default:
	}

	log.Println("op is ", op.Op)

	return nil
}
