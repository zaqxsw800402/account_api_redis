package redis

import (
	"context"
	"fmt"
	"log"
)

type Account struct {
	AccountId   uint    `redis:"AccountId" `
	CustomerId  uint    `redis:"CustomerId" `
	AccountType string  `redis:"AccountType" `
	Amount      float64 `redis:"Amount" `
	Status      string  `redis:"Status" `
}

func (d Database) getUserKeyForAccount(userID int) string {
	return fmt.Sprintf("%d:account", userID)
}

func (d Database) getUserValueForAccount(accountID string) string {
	return fmt.Sprintf("account:%s", accountID)
}

// GetAccount 查詢帳戶資料
func (d Database) GetAccount(ctx context.Context, accountKey string) (*Account, error) {
	// 將資料儲存到結構體裡
	var account Account
	err := d.RC.HGetAll(ctx, accountKey).Scan(&account)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &account, nil
}

func (d Database) DeleteAccount(ctx context.Context, accountID string) error {
	userValue := d.getUserValueForAccount(accountID)
	result := d.RC.Del(ctx, userValue)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
