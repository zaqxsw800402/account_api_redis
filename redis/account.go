package redis

import (
	"context"
	"fmt"
	"log"
	"red/cmd/api/dto"
	"sort"
	"strconv"
	"time"
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

func (d Database) getUserValueForAccount(userID int, accountID string) string {
	return fmt.Sprintf("%d:account:%s", userID, accountID)
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

func (d Database) GetAllAccounts(ctx context.Context, userID int) ([]dto.AccountResponse, error) {
	userKey := d.getUserKeyForAccount(userID)
	members := d.RC.SMembers(ctx, userKey)

	accounts := make([]dto.AccountResponse, 0)

	for _, member := range members.Val() {
		account, err := d.GetAccount(ctx, member)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, dto.AccountResponse{
			AccountId:   account.AccountId,
			CustomerId:  account.CustomerId,
			AccountType: account.AccountType,
			Amount:      account.Amount,
			Status:      account.Status,
		})
	}
	sort.SliceStable(accounts, func(i int, j int) bool {
		return accounts[i].CustomerId < accounts[j].CustomerId
	})

	return accounts, nil
}

// SaveAccount 儲存資料到redis裡面
func (d Database) SaveAccount(ctx context.Context, userID int, a dto.AccountResponse) {
	userKey := d.getUserKeyForAccount(userID)
	userValue := d.getUserValueForAccount(userID, strconv.Itoa(int(a.AccountId)))
	// 製作專屬的key
	d.RC.SAdd(ctx, userKey, userValue)
	d.RC.Expire(ctx, userKey, time.Hour*24)

	d.RC.HSet(ctx, userValue,
		"AccountId", a.AccountId,
		"CustomerId", a.CustomerId,
		"AccountType", a.AccountType,
		"Amount", a.Amount,
		"Status", a.Status,
	)
	// 設定過期時間
	d.RC.Expire(ctx, userValue, time.Hour*24)
}

func (d Database) SaveAllAccounts(ctx context.Context, userID int, accounts []dto.AccountResponse) {
	for _, account := range accounts {
		d.SaveAccount(ctx, userID, account)
	}
}

func (d Database) DeleteAccount(ctx context.Context, userID int, accountID string) error {
	userValue := d.getUserValueForAccount(userID, accountID)
	result := d.RC.HSet(ctx, userValue, "Status", 0)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
