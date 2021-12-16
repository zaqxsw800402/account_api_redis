package mysql

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func GetDBClient() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false", dbUser, dbPassword, dbHost, dbPort, dbName)
	// 建立與資料庫的連結
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Println("failed to connect mysql " + err.Error())
	}

	sqlDB, _ := db.DB()
	// 最多閒置數量
	sqlDB.SetMaxIdleConns(10)
	// 最多連接數量
	sqlDB.SetMaxOpenConns(10)
	// 等待醉酒時間
	sqlDB.SetConnMaxIdleTime(time.Hour)

	return db
}
