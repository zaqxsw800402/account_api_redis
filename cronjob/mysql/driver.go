package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func GetDBClient(dsn string) *gorm.DB {

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
