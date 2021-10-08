package gin_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/dto"
	"red/service"
)

type AccountHandler struct {
	service service.AccountService
}

func (h AccountHandler) NewAccount(c *gin.Context) {
	customerId := c.Param("id")
	var request dto.NewAccountRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = customerId
		account, appError := h.service.NewAccount(request)
		if appError != nil {
			c.JSON(appError.Code, appError.AsMessage())
		} else {
			c.JSON(http.StatusOK, account)
		}
	}
}

func (h AccountHandler) MakeTransaction(c *gin.Context) {
	customerId := c.Param("id")
	accountId := c.Param("account_id")
	var request dto.TransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		request.AccountId = accountId
		request.CustomerId = customerId
		//fmt.Println(request)
		account, appError := h.service.MakeTransaction(request)

		if appError != nil {
			c.JSON(appError.Code, appError.AsMessage())
		} else {
			c.JSON(http.StatusOK, account)
		}
	}

}
