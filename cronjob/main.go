package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"red/mysql"
	"red/redis"
)

func main() {
	fmt.Println("Start cronjob")

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false", dbUser, dbPassword, dbHost, dbPort, dbName)

	redisHost := os.Getenv("REDIS_HOST")

	dbClient := mysql.GetDBClient(dsn)
	db := mysql.NewDB(dbClient)

	client, err := redis.GetClient(redisHost)
	if err != nil {
		log.Println("redis connection failed" + err.Error())
	}
	redisDB := redis.New(client)

	job := Job{
		redis: redisDB,
		db:    db,
		Shut:  make(chan int, 1),
	}

	// clear data at every 6 am
	go StartJob("*/2 * * * *", job)
	select {}
}
