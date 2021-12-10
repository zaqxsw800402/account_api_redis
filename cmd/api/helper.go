package main

import (
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func badRequest(c *gin.Context, code int, err error) {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	c.JSON(code, &payload)
}

func (app *application) invalidCredentials(c *gin.Context) {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	c.JSON(http.StatusUnauthorized, payload)
}

func newSession(client *gorm.DB) *scs.SessionManager {
	sqlDB, _ := client.DB()

	// 最多閒置數量
	sqlDB.SetMaxIdleConns(10)
	// 最多連接數量
	sqlDB.SetMaxOpenConns(10)
	// 等待醉酒時間
	sqlDB.SetConnMaxIdleTime(time.Hour)

	// set up session
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(sqlDB)
	return session
}
