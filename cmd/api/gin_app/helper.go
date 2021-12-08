package gin_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

func invalidCredentials(c *gin.Context) {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	c.JSON(http.StatusUnauthorized, payload)
}
