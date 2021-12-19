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

	server.POST("/api/authenticate", app.createAuthToken)
	server.POST("/api/all-users", app.newUser)
	server.POST("/api/forgot-password", app.SendPasswordResetEmail)
	server.POST("/api/reset-password", app.ResetPassword)

	// 使用zap來記錄api的使用
	//server.Use(logger.GinLogger(), logger.GinRecovery(true))

	// admin group
	admin := server.Group("/api/admin", app.Auth())

	//admin.GET("/all-users", app.getAllUsers)
	//建立顧客
	admin.POST("/all-customers/0", app.newCustomer)
	//查詢特定顧客的資訊
	admin.GET("/all-customers/:id", app.getCustomer)
	//查詢全部顧客
	admin.GET("/all-customers", app.getAllCustomers)
	//查詢全部顧客
	admin.POST("/all-customers/delete/:id", app.deleteCustomer)

	// 在特地的顧客下創立帳戶
	admin.POST("/all-customers/:id/accounts/0", app.newAccount)
	// 在特地的顧客下創立帳戶
	admin.POST("/all-customers/accounts/0", app.newAccount)
	// 查詢的已知customer id 顧客下查詢帳戶
	admin.GET("/all-customers/:id/accounts", app.getAllAccounts)
	// 刪除帳戶
	admin.POST("/all-customers/:id/accounts/delete/:account_id", app.deleteAccount)
	// 查詢的已知user id id 顧客下查詢帳戶
	admin.GET("/all-customers/accounts", app.getAllAccountWithUserID)
	// 查詢該帳戶的資料
	admin.GET("/all-customers/:id/accounts/:account_id/transactions", app.getAllTransactions)
	// 提供儲存或領取金錢
	admin.POST("/withdrawal", app.makeTransaction)
	admin.POST("/transfer", app.transfer)

	// check whether customer id is in user id
	admin.POST("/check-customer_id", app.checkUserID)
	//admin.POST("/all-customers/:id/accounts/:account_id", app.ah.makeTransaction)

	return server
}
