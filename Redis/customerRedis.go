package Redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

const (
	Customer = "customer"
	Account  = "account"
)

var ctx = context.Background()

//var RC *redis.Client

type Database struct {
	RC *redis.Client
}

func NewRedisDb(RC *redis.Client) Database {
	return Database{RC}
}

//func init() {
//	RC = newClient()
//}

// newClient 建立連接

// CustomerTimes 設定時間內最多查詢次數
func (d Database) CustomerTimes(customerId string) *errs.AppError {
	return d.addKeys(Customer, customerId, 5)
}

// TransactionTimes 設定時間內最多查詢次數
func (d Database) TransactionTimes(accountId string) *errs.AppError {
	return d.addKeys(Account, accountId, 5)
}

// addKeys 設定限制次數及時間
func (d Database) addKeys(keyPrefix string, keyId string, times int) *errs.AppError {
	// 製造專屬的key
	key := fmt.Sprintf("%s:%s:times", keyPrefix, keyId)
	// 每查詢一次，值增加一
	d.RC.Incr(ctx, key)
	// 設定過期時間
	d.RC.Expire(ctx, key, time.Minute)
	// 讀取次數，並判斷是否達到次數限制
	result := d.RC.Get(ctx, key)
	value, _ := result.Int()
	if value > times {
		return errs.TooManyTimes("Logging too many times")
	}
	return nil
}

// GetCustomer 讀取使用者資料
func (d Database) GetCustomer(id string) *CustomerResponse {
	// 設定專屬的key
	key := fmt.Sprintf("%s:%s:response", Customer, id)
	// 讀取資料
	res := d.RC.HGetAll(ctx, key)
	if _, ok := res.Val()["Id"]; !ok {
		return nil
	}

	// 將資料輸入進結構體哩，並且回傳
	var customer CustomerResponse
	if err := res.Scan(&customer); err != nil {
		panic(err)
	}
	return &customer
}

// SaveCustomer 存資料到redis
func (d Database) SaveCustomer(c *dto.CustomerResponse) {
	// 製造專屬的key
	key := fmt.Sprintf("%s:%d:response", Customer, c.Id)
	// 存進資料
	d.RC.HSet(ctx, key,
		"Id", c.Id,
		"Name", c.Name,
		"City", c.City,
		"Zipcode", c.Zipcode,
		"DateOfBirth", c.DateOfBirth,
		"Status", c.Status)
	// 設定過期時間
	d.RC.Expire(ctx, key, time.Minute)
}

type CustomerResponse struct {
	Id          uint   `redis:"Id"`
	Name        string `redis:"Name"`
	City        string `redis:"City"`
	Zipcode     string `redis:"Zipcode"`
	DateOfBirth string `redis:"DateOfBirth"`
	Status      string `redis:"Status"`
}
