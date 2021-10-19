package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"red/domain"
	"red/logger"
	"red/service"
	"time"
)

func Start() {
	// 讀取.env檔案
	loadEnv()

	//建立連線池
	dbClient := getDBClient()

	server := gin.Default()

	// 使用zap來記錄api的使用
	server.Use(logger.GinLogger(), logger.GinRecovery(true))

	//建立各個Repository
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)

	//建立各個Handlers
	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}

	//查詢全部顧客
	server.GET("/customers", ch.getAllCustomers)
	//建立顧客
	server.POST("/customers", ch.newCustomers)
	//查詢特定顧客的資訊
	server.GET("/customers/:id", ch.getCustomer)

	// 在特地的顧客下創立帳戶
	server.POST("/customers/:id/account", ah.NewAccount)
	// 查詢該帳戶的資料
	server.GET("/customers/:id/account/:account_id", ah.getAccount)
	// 提供儲存或領取金錢
	server.POST("/customers/:id/account/:account_id", ah.MakeTransaction)
	err := server.Run(":8080")
	if err != nil {
		logger.Error("failed to run server, error: " + err.Error())
	}
}

// loadEnv 讀取.env檔案
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file, error: " + err.Error())
	}
}

// getDBClient 建立與db的連接
func getDBClient() *gorm.DB {
	// 讀取環境變數
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 建立與資料庫的連結
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbAddress, dbPort, dbName),
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		logger.Error("failed to connect db " + err.Error())
	}

	sqlDB, _ := db.DB()
	// 最多閒置數量
	sqlDB.SetMaxIdleConns(10)
	// 最多連接數量
	sqlDB.SetMaxOpenConns(10)
	// 等待醉酒時間
	sqlDB.SetConnMaxIdleTime(time.Hour)

	//建立表格，如果沒有表格
	err = db.AutoMigrate(&domain.Customer{}, &domain.Account{}, &domain.Transaction{})
	if err != nil {
		logger.Error("Failed to create tables" + err.Error())
		return nil
	}

	return db
}
