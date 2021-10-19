package gin_app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"red/Redis"
	"red/dto"
	"red/logger"
	"red/service"
	"strconv"
)

type AccountHandler struct {
	service service.AccountService
	redisDB Redis.Database
}

// NewAccount 申請新帳戶
func (h AccountHandler) NewAccount(c *gin.Context) {
	// 讀取id的值
	customerId := c.Param("id")
	var request dto.NewAccountRequest
	// 讀取BODY裡的資料
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 轉換id成uint格式
	id, err := strconv.ParseUint(customerId, 10, 64)
	if err != nil {
		logger.Error("Failed to parse customerId " + err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	request.CustomerId = uint(id)

	// 建立新帳戶
	account, appError := h.service.NewAccount(request)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, account)
	}

}

// MakeTransaction 申請交易
func (h AccountHandler) MakeTransaction(c *gin.Context) {
	// 讀取資料
	customerId := c.Param("id")
	accountId := c.Param("account_id")

	// 紀錄交易的次數
	appError := h.redisDB.TransactionTimes(accountId)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	}

	// 讀取body裡的資料
	var request dto.TransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 轉換accountId的資料格式
	id, _ := strconv.ParseUint(accountId, 10, 64)
	// 輸入資料到request裡
	request.AccountId = uint(id)
	request.CustomerId = customerId

	// 申請交易
	account, appError := h.service.MakeTransaction(request)
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
	} else {
		c.JSON(http.StatusOK, account)
	}

}

// getAccount 讀取帳戶資料
func (h AccountHandler) getAccount(c *gin.Context) {
	accountId := c.Param("account_id")

	// 檢查Redis裡是否已經有資料
	if account := h.redisDB.GetAccount(accountId); account != nil {
		c.JSON(http.StatusOK, account)
		return
	}

	// 轉換accountId資料格式
	id, _ := strconv.ParseUint(accountId, 10, 64)

	// 讀取db裡的資料
	account, appError := h.service.GetAccount(uint(id))
	if appError != nil {
		c.JSON(appError.Code, appError.AsMessage())
		return
	} else {
		c.JSON(http.StatusOK, account)
	}

	// 將資料存進Redis
	h.redisDB.SaveAccount(account)
}
