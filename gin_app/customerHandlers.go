package gin_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/dto"
	"red/service"
)

type CustomerHandlers struct {
	service service.CustomerService
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
	customer, appError := ch.service.GetCustomer(id)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, customer)
	}
}

func (ch *CustomerHandlers) newCustomers(c *gin.Context) {
	var customer dto.CustomerRequest
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	response, appError := ch.service.SaveCustomer(customer)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, response)
	}

}
