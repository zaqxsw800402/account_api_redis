package main

import (
	"errors"
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
func jsonResp(c *gin.Context, code int) {
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = false

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

func (app *application) checkCustomerInUser(userID int, customerID int) bool {
	customers, err := app.ch.service.GetAllCustomers(userID)
	if err != nil {
		return false
	}

	for _, c := range customers {
		i := int(c.Id)
		if i == customerID {
			return true
		}
	}

	return false
}

type checkUser struct {
	CustomerID int `json:"customer_id" binding:"required"`
	AccountId  int `json:"account_id" binding:"required"`
	Amount     int `json:"amount" binding:"required"`
}

func (app *application) checkUserID(c *gin.Context) {
	var req checkUser
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	userID := c.GetInt("userID")
	err := c.ShouldBindJSON(&req)
	if err != nil {
		badRequest(c, http.StatusOK, errors.New("must be numbers"))
		return
	}

	ok := app.checkCustomerInUser(userID, req.CustomerID)

	switch {
	case !ok:
		payload.Error = true
	case ok:
		payload.Error = false
	}

	c.JSON(http.StatusOK, &payload)
}
