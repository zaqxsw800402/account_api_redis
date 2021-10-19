package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"red/Redis"
	"red/dto"
	"red/service"
)

type CustomerHandlers struct {
	service service.CustomerService
	redisDB Redis.Database
}

func (ch *CustomerHandlers) getAllCustomers(c *gin.Context) {
	status := c.Query("status")
	customers, err := ch.service.GetAllCustomer(status)
	if err != nil {
		c.JSON(err.Code, err.AsMessage())
	} else {
		c.JSON(http.StatusOK, customers)
	}
}

func (ch *CustomerHandlers) getCustomer(c *gin.Context) {
	id := c.Param("id")

	// 檢查時間內讀取的次數，太多次則拒絕提供資料
	appError := ch.redisDB.CustomerTimes(id)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	}

	if customer := ch.redisDB.GetCustomer(id); customer != nil {
		c.JSON(http.StatusOK, customer)
		fmt.Println("Using redis")
		return
	}

	customer, appError := ch.service.GetCustomer(id)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	} else {
		c.JSON(http.StatusOK, customer)
	}

	ch.redisDB.SaveCustomer(customer)

}

func (ch *CustomerHandlers) newCustomers(c *gin.Context) {
	var customer dto.CustomerRequest
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	response, appError := ch.service.SaveCustomer(customer)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, response)
	}

}
