package Redis

import (
	"fmt"
	"red/domain"
	"red/logger"
	"time"
)

type AccountResponse struct {
	AccountId   uint    `redis:"AccountId" json:"account_id"`
	CustomerId  uint    `redis:"CustomerId" json:"customer_id"`
	OpeningDate string  `redis:"OpeningDate" json:"opening_date"`
	AccountType string  `redis:"AccountType" json:"account_type"`
	Amount      float64 `redis:"Amount" json:"amount"`
	Status      string  `redis:"Status" json:"status"`
}

// GetAccount 查詢帳戶資料
func GetAccount(id string) *AccountResponse {
	// 製作redis裡的專屬key
	key := fmt.Sprintf("%s:%s:response", Account, id)
	// 讀取redis裡的資料
	res := RC.HGetAll(ctx, key)
	// 判斷res裡面是否有值
	if _, ok := res.Val()["AccountId"]; !ok {
		return nil
	}

	// 將資料儲存到結構體裡
	var account AccountResponse
	if err := res.Scan(&account); err != nil {
		logger.Error("failed to get account data from redis, error: " + err.Error())
	}

	return &account
}

// SaveAccount 儲存資料到redis裡面
func SaveAccount(c *domain.Account) {
	// 製作專屬的key
	key := fmt.Sprintf("%s:%d:response", Account, c.AccountId)
	RC.HSet(ctx, key,
		"AccountId", c.AccountId,
		"CustomerId", c.CustomerId,
		"OpeningDate", c.OpeningDate,
		"AccountType", c.AccountType,
		"Amount", c.Amount,
		"Status", c.Status,
	)
	// 設定過期時間
	RC.Expire(ctx, key, time.Minute)
}
