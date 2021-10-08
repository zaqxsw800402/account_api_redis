package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"red/domain"
	"red/service"
	"time"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := gin.Default()
	dbClient := getDBClient()
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)

	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}

	server.GET("/customers", ch.getAllCustomers)
	server.POST("/customers", ch.newCustomers)
	server.GET("/customers/:id", ch.getCustomer)

	server.POST("/customers/:id/account", ah.NewAccount)
	server.GET("/customers/:id/account/:account_id", ah.getAccount)
	server.POST("/customers/:id/account/:account_id", ah.MakeTransaction)
	server.Run(":8080")
}

func getDBClient() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbAddress, dbPort, dbName),
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	db.AutoMigrate(&domain.Customer{}, &domain.Account{}, &domain.Transaction{})

	return db
}
