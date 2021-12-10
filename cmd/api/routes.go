package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) routes() http.Handler {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		//AllowOrigins:   []string{"https://*", "http://*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowHeaders:   []string{"Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		AllowHeaders: []string{"Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 使用zap來記錄api的使用
	//server.Use(logger.GinLogger(), logger.GinRecovery(true))

	// admin group
	admin := server.Group("/api/admin", app.Auth())
	admin.POST("/all-users", app.uh.getAllUsers)
	//建立顧客
	admin.POST("/all-customers/edit/:id", app.editCustomers)
	//查詢特定顧客的資訊
	admin.GET("/all-customers/:id", app.getCustomer)
	//查詢全部顧客
	admin.GET("/all-customers", app.getAllCustomers)
	//查詢全部顧客
	admin.POST("/all-customers/delete/:id", app.deleteCustomers)

	//建立各個Handlers
	server.POST("api/authenticate", app.uh.createAuthToken)

	// 在特地的顧客下創立帳戶
	server.POST("/api/all-customers/:id/accounts", app.ah.newAccount)
	// 在特地的顧客下查詢帳戶
	server.GET("/api/all-customers/:id/accounts", app.ah.getAllAccount)
	// 查詢該帳戶的資料
	server.GET("/api/all-customers/:id/accounts/:account_id", app.ah.getAccount)
	// 提供儲存或領取金錢
	server.POST("/api/all-customers/:id/accounts/:account_id", app.ah.makeTransaction)

	return server
}
