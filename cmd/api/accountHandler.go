package main

import "C"
import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"red/cmd/api/dto"
	"red/cmd/api/service"
	"strconv"
)

type AccountHandler struct {
	service service.AccountService
	//redisDB Redis.Database
}

// newAccount 申請新帳戶
func (app *application) newAccount(c *gin.Context) {
	// 讀取id的值
	customerId := c.Param("id")
	var request dto.AccountRequest
	// 讀取BODY裡的資料
	err := c.ShouldBindJSON(&request)
	if err != nil {
		badRequest(c, http.StatusBadRequest, errors.New("amount must be a number"))
		return
	}

	// 轉換id成uint格式
	id, err := strconv.ParseUint(customerId, 10, 64)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
		return
	}
	request.CustomerId = uint(id)

	// 建立新帳戶
	_, appError := app.ah.service.NewAccount(request)
	if appError != nil {
		//c.JSON(appError.Code, appError.AsMessage())
		badRequest(c, appError.Code, appError)
		return
	}

	jsonResp(c, http.StatusOK)
}

func (app *application) deleteAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	err := app.ah.service.DeleteAccount(accountID)
	if err != nil {
		badRequest(c, http.StatusBadRequest, err)
	}
}

// makeTransaction 申請交易
func (app *application) makeTransaction(c *gin.Context) {
	// 紀錄交易的次數
	//appError := h.redisDB.TransactionTimes(accountId)
	//if appError != nil {
	//	c.JSON(appError.Code, appError.AsMessage())
	//	return
	//}

	// 讀取body裡的資料
	var request dto.TransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 申請交易
	_, appError := app.ah.service.MakeTransaction(request)
	if appError != nil {
		//c.JSON(appError.Code, appError.AsMessage())
		badRequest(c, appError.Code, appError)
	} else {
		//c.JSON(http.StatusOK, account)
		jsonResp(c, http.StatusOK)
	}
}

func (app *application) transfer(c *gin.Context) {
	var request dto.TransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, appError := app.ah.service.Transfer(request)
	if appError != nil {
		//c.JSON(appError.Code, appError.AsMessage())
		badRequest(c, appError.Code, appError)
	} else {
		//c.JSON(http.StatusOK, account)
		jsonResp(c, http.StatusOK)
	}
}

// getAllTransactions 讀取帳戶資料
func (app *application) getAllTransactions(c *gin.Context) {
	accountId := c.Param("account_id")
	customerId := c.Param("id")

	// 檢查Redis裡是否已經有資料
	//if account := h.redisDB.GetAccount(accountId); account != nil {
	//	c.JSON(http.StatusOK, account)
	//	return
	//}

	// 轉換accountId資料格式
	id, err := strconv.ParseUint(accountId, 10, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cId, err := strconv.ParseUint(customerId, 10, 64)
	if err != nil {
		//c.String(http.StatusBadRequest, err.Error())
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	// 讀取db裡的資料
	transactions, appError := app.ah.service.GetALlTransactions(uint(cId), uint(id))
	if appError != nil {
		//c.JSON(appError.Code, appError.AsMessage())
		badRequest(c, appError.Code, appError)
		return
	} else {
		c.JSON(http.StatusOK, transactions)
	}

	// 將資料存進Redis
	//h.redisDB.SaveAccount(account)
}

// getAllAccounts 讀取帳戶資料
func (app *application) getAllAccounts(c *gin.Context) {
	customerId := c.Param("id")

	// 檢查Redis裡是否已經有資料
	//if account := h.redisDB.GetAccount(accountId); account != nil {
	//	c.JSON(http.StatusOK, account)
	//	return
	//}

	// 轉換accountId資料格式
	id, err := strconv.ParseUint(customerId, 10, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, err.Error())
		badRequest(c, http.StatusBadRequest, err)
		return
	}

	// 讀取db裡的資料
	accounts, appError := app.ah.service.GetAllAccounts(uint(id))
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	} else {
		c.JSON(http.StatusOK, accounts)
	}

	// 將資料存進Redis
	//h.redisDB.SaveAccount(account)
}

func (app *application) getAllAccountWithUserID(c *gin.Context) {
	userID := c.GetInt("userID")

	customers, appError := app.ch.service.GetAllCustomer(userID)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	}
	// 檢查Redis裡是否已經有資料
	//if account := h.redisDB.GetAccount(accountId); account != nil {
	//	c.JSON(http.StatusOK, account)
	//	return
	//}

	type Response struct {
		AccountId   uint    ` json:"account_id"`
		AccountType string  ` json:"account_type"`
		Amount      float64 ` json:"amount"`
		Status      string  ` json:"status"`
		CustomerId  uint    ` json:"customer_id"`
	}

	resp := make([]Response, 0)
	for _, customer := range customers {
		accounts, appError := app.ah.service.GetAllAccounts(customer.Id)
		if appError != nil {
			c.JSON(appError.Code, appError.AsMessage())
			return
		}
		for _, a := range accounts {
			r := Response{
				AccountId:   a.AccountId,
				AccountType: a.AccountType,
				Amount:      a.Amount,
				Status:      a.Status,
				CustomerId:  customer.Id,
			}
			resp = append(resp, r)
		}

	}
	c.JSON(http.StatusOK, resp)

	// 將資料存進Redis
	//h.redisDB.SaveAccount(account)
}
