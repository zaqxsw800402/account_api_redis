package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/cmd/api/dto"
	"red/cmd/api/service"
	"strconv"
)

type CustomerHandler struct {
	service service.CustomerService
	//redisDB Redis.Database
}

func (app *application) getAllCustomers(c *gin.Context) {
	userID := c.GetInt("userID")

	//status := c.Query("status")
	customers, err := app.ch.service.GetAllCustomer(userID)
	if err != nil {
		badRequest(c, err.Code, err)
	} else {
		c.JSON(http.StatusOK, customers)
	}
}

func (app *application) getCustomer(c *gin.Context) {
	id := c.Param("id")

	// 檢查時間內讀取的次數，太多次則拒絕提供資料
	//appError := app.ch.redisDB.CustomerTimes(id)
	//if appError != nil {
	//	c.JSON(appError.Code, appError.AsMessage())
	//	return
	//}

	//if customer := app.ch.redisDB.GetCustomer(id); customer != nil {
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

	//app.ch.redisDB.SaveCustomer(customer)

}

func (app *application) newCustomer(c *gin.Context) {
	userID := c.GetInt("userID")

	var customer dto.CustomerRequest
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	_, appError := app.ch.service.SaveCustomer(userID, customer)
	if appError != nil {
		badRequest(c, appError.Code, appError)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false

	c.JSON(http.StatusOK, resp)
}

func (app *application) deleteCustomer(c *gin.Context) {
	customerId := c.Param("id")
	err := app.ch.service.DeleteCustomer(customerId)
	if err != nil {
		badRequest(c, err.Code, err)
	}

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
