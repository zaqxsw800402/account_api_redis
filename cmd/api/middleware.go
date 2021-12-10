package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := app.uh.authenticateToken(c)
		if err != nil {
			app.invalidCredentials(c)
			return
		}
	}
}
