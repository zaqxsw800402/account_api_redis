package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"red/Redis"
	"red/dto"
	"red/service"
	"strconv"
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
		id, _ := strconv.ParseUint(customerId, 10, 64)
		request.CustomerId = uint(id)
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
	appError := Redis.TransactionTimes(accountId)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	}
	var request dto.TransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		id, _ := strconv.ParseUint(accountId, 10, 64)
		request.AccountId = uint(id)
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

func (h AccountHandler) getAccount(c *gin.Context) {
	accountId := c.Param("account_id")

	if account := Redis.GetAccount(accountId); account != nil {
		c.JSON(http.StatusOK, account)
		fmt.Println("Using redis")
		return
	}

	id, _ := strconv.ParseUint(accountId, 10, 64)
	account, appError := h.service.GetAccount(uint(id))
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, account)
	}

	Redis.SaveAccount(account)
}
