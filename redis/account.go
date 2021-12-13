package redis

import (
	"context"
	"fmt"
	"red/cmd/api/domain"
	"red/cmd/api/dto"
	"red/logger"
	"time"
)

type Account struct {
	AccountId   uint    `redis:"AccountId" `
	CustomerId  uint    `redis:"CustomerId" `
	OpeningDate string  `redis:"OpeningDate" `
	AccountType string  `redis:"AccountType" `
	Amount      float64 `redis:"Amount" `
	Status      string  `redis:"Status" `
}

func (d Database) getUserValueForAccount(userID int, accountID string) string {
	return fmt.Sprintf("%d:account:%s", userID, accountID)
}

func (d Database) getUserKeyForAccount(userID int) string {
	return fmt.Sprintf("%d:customer", userID)
}

// GetAccount 查詢帳戶資料
func (d Database) GetAccount(id string) *Account {
	// 製作redis裡的專屬key
	key := fmt.Sprintf("%s:%s:response", "Account", id)
	// 讀取redis裡的資料
	res := d.RC.HGetAll(ctx, key)
	// 判斷res裡面是否有值
	if _, ok := res.Val()["AccountId"]; !ok {
		return nil
	}

	// 將資料儲存到結構體裡
	var account Account
	if err := res.Scan(&account); err != nil {
		logger.Error("failed to get account data from redis, error: " + err.Error())
	}

	return &account
}

// SaveAccount 儲存資料到redis裡面
func (d Database) SaveAccount(ctx context.Context, userID int, a dto.AccountResponse) {
	userKey := getUserValueForAccount(userID)
	// 製作專屬的key
	key := fmt.Sprintf("%s:%d:response", "Account", c.AccountId)
	d.RC.HSet(ctx, key,
		"AccountId", c.AccountId,
		"CustomerId", c.CustomerId,
		"OpeningDate", c.OpeningDate,
		"AccountType", c.AccountType,
		"Amount", c.Amount,
		"Status", c.Status,
	)
	// 設定過期時間
	d.RC.Expire(ctx, key, time.Minute)
}
