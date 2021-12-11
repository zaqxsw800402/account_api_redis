package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := app.authenticateToken(c)
		if err != nil {
			app.invalidCredentials(c)
			return
		}

		c.Set("userID", user.ID)
	}
}
