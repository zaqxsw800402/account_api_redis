package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"red/cmd/api/domain"
	"red/internal/models"
)

func GetDBClient(dsn string) (*gorm.DB, error) {
	// 讀取環境變數
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 建立與資料庫的連結
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName),
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		//log.Error("failed to connect db " + err.Error())
		//log.Println("failed to connect db " + err.Error())
		return nil, err
	}

	//建立表格，如果沒有表格
	err = db.AutoMigrate(&domain.Customer{}, &domain.Account{}, &domain.Transaction{}, &models.Token{}, &models.User{})
	if err != nil {
		//logger.Error("Failed to create tables" + err.Error())
		//log.Println("Failed to create tables" + err.Error())
		return nil, err
	}

	return db, nil
}

func OpenDb(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}
