package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"red/Redis"
	"red/cmd/api/domain"
	"red/cmd/api/service"
	"red/logger"
	"time"
)

func Start() {
	logger.ZapInit()
	// 讀取.env檔案
	loadEnv()

	//建立連線池
	dbClient := getDBClient()
	redisClient := getRedisClient()

	server := gin.Default()

	// 使用zap來記錄api的使用
	server.Use(logger.GinLogger(), logger.GinRecovery(true))

	//建立各個Repository
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	userRepositoryDb := domain.NewUserRepositoryDb(dbClient)

	//建立redis
	redisDB := Redis.NewRedisDb(redisClient)

	//建立各個Handlers
	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb), redisDB}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb), redisDB}
	uh := UserHandlers{service.NewUserService(userRepositoryDb)}

	server.POST("/all-users", uh.getAllUsers)
	server.POST("api/authenticate", uh.CreateAuthToken)

	//查詢全部顧客
	server.GET("/api/customers", ch.getAllCustomers)
	//建立顧客
	server.POST("/api/customers", ch.newCustomers)
	//查詢特定顧客的資訊
	server.GET("/api/customers/:id", ch.getCustomer)

	// 在特地的顧客下創立帳戶
	server.POST("/api/customers/:id/account", ah.newAccount)
	// 查詢該帳戶的資料
	server.GET("/api/customers/:id/account/:account_id", ah.getAccount)
	// 提供儲存或領取金錢
	server.POST("/api/customers/:id/account/:account_id", ah.makeTransaction)
	err := server.Run(":4001")
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

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//pong, err := client.Ping(ctx).Result()
	//log.Println(pong)
	//if err != nil {
	//	log.Fatal(err)
	//}
	return client
}
