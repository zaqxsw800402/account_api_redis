package domain

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func GetDBClient(dsn string) *gorm.DB {
	// 建立與資料庫的連結
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Println("failed to connect db " + err.Error())
	}

	//sqlDB, _ := db.DB()
	//// 最多閒置數量
	//sqlDB.SetMaxIdleConns(10)
	//// 最多連接數量
	//sqlDB.SetMaxOpenConns(10)
	//// 等待醉酒時間
	//sqlDB.SetConnMaxIdleTime(time.Hour)

	//建立表格，如果沒有表格
	err = db.AutoMigrate(&Customer{}, &Account{}, &Transaction{}, &User{}, &Token{})
	if err != nil {
		log.Println("Failed to create tables" + err.Error())
		return nil
	}

	return db
}
