package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"red/internal/models"
)

func GetDBClient(dsn string) (*gorm.DB, error) {
	// 讀取環境變數
	// 建立與資料庫的連結
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		//log.Error("failed to connect db " + err.Error())
		//log.Println("failed to connect db " + err.Error())
		return nil, err
	}

	//建立表格，如果沒有表格
	err = db.AutoMigrate(&models.Session{})
	if err != nil {
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
