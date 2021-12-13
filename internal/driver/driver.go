package driver

import (
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
		//log.Error("failed to connect mysql " + err.Error())
		//log.Println("failed to connect mysql " + err.Error())
		return nil, err
	}

	//建立表格，如果沒有表格
	err = db.AutoMigrate(&models.Session{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
