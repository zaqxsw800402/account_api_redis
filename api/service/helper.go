package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"red/errs"
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

func passwordMatches(hash, password string) (bool, *errs.AppError) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, errs.NewNotFoundError("password mismatch" + err.Error())
		default:
			return false, errs.NewNotFoundError("password mismatch" + err.Error())
		}
	}
	return true, nil
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
