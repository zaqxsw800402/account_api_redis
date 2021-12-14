package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/dto"
	"red/service"
	"strconv"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (app *application) getAllCustomers(c *gin.Context) {
	userID := c.GetInt("userID")

	allCustomers, _ := app.redis.GetAllCustomers(c, userID)
	if len(allCustomers) != 0 {
		c.JSON(http.StatusOK, allCustomers)
		return
	}

	customers, err2 := app.ch.service.GetAllCustomer(userID)
	if err2 != nil {
		badRequest(c, err2.Code, err2)
	}

	app.redis.SaveAllCustomers(c, userID, customers)

	c.JSON(http.StatusOK, customers)

}

func (app *application) getCustomer(c *gin.Context) {
	id := c.Param("id")

	//if customer := app.ch.redis.GetCustomer(id); customer != nil {
	//	c.JSON(http.StatusOK, customer)
	//	fmt.Println("Using redis")
	//	return
	//}

	customer, appError := app.ch.service.GetCustomer(id)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	} else {
		c.JSON(http.StatusOK, customer)
	}

	//app.ch.redis.SaveCustomer(customer)

}

func (app *application) newCustomer(c *gin.Context) {
	userID := c.GetInt("userID")

	var req dto.CustomerRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	customer, appError := app.ch.service.SaveCustomer(userID, req)
	if appError != nil {
		badRequest(c, appError.Code, appError)
		return
	}

	app.redis.SaveCustomer(c, userID, *customer)

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false

	c.JSON(http.StatusOK, resp)
}

func (app *application) deleteCustomer(c *gin.Context) {
	userID := c.GetInt("userID")
	customerId := c.Param("id")
	err := app.ch.service.DeleteCustomer(customerId)
	if err != nil {
		badRequest(c, err.Code, err)
	}

	_ = app.redis.DeleteCustomer(userID, customerId)

	id, appError := strconv.ParseUint(customerId, 10, 64)
	if appError != nil {
		//c.JSON(http.StatusBadRequest, err.Error())
		badRequest(c, http.StatusBadRequest, appError)
		return
	}

	accounts, appError := app.ah.service.GetAllAccounts(uint(id))
	for _, account := range accounts {
		appError := app.ah.service.DeleteAccount(strconv.Itoa(int(account.AccountId)))
		if appError != nil {
			//c.JSON(http.StatusBadRequest, err.Error())
			badRequest(c, http.StatusBadRequest, appError)
			return
		}
	}
}
