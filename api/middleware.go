package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (app *application) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := app.authenticateToken(c)

		if err != nil {
			log.Println("auth failed ", err.Error())
			app.invalidCredentials(c)
			c.Abort()
			return
		}

		// UserTimes 設定1 min內最多查詢次數
		err = app.redis.UserTimes(c, *userID, 60)
		if err != nil {
			//c.String(err.Code, err.Message)
			app.invalidCredentials(c)
			c.Abort()
			return
		}

		c.Set("userID", *userID)
	}
}
