package Redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"red/dto"
	"red/errs"
	"time"
)

const (
	Customer = "customer"
	Account  = "account"
)

var ctx = context.Background()
var RC *redis.Client

func init() {
	RC = newClient()

	RC.Del(ctx, Customer)
	RC.Del(ctx, Account)

}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(ctx).Result()
	log.Println(pong)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func CustomerTimes(customerId string) *errs.AppError {
	return addKeys(Customer, customerId, 5)
}

func TransactionTimes(accountId string) *errs.AppError {
	return addKeys(Account, accountId, 5)
}

func addKeys(keyPrefix string, keyId string, times int) *errs.AppError {
	key := fmt.Sprintf("%s:%s:times", keyPrefix, keyId)
	RC.Incr(ctx, key)
	result := RC.Get(ctx, key)
	value, _ := result.Int()
	RC.Expire(ctx, key, time.Minute)
	if value > times {
		return errs.TooManyTimes("Logging too many times")
	}
	return nil
}

func GetCustomer(id string) *CustomerResponse {
	key := fmt.Sprintf("%s:%s:response", Customer, id)
	res := RC.HGetAll(ctx, key)
	result, _ := res.Result()
	//fmt.Println(result["Id"],err)
	if _, err := result["Id"]; !err {
		return nil
	}

	var customer CustomerResponse
	if err := res.Scan(&customer); err != nil {
		panic(err)
	}
	//fmt.Println("all", customer)
	return &customer
}

func SaveCustomer(c *dto.CustomerResponse) {
	key := fmt.Sprintf("%s:%d:response", Customer, c.Id)
	RC.HSet(ctx, key,
		"Id", c.Id,
		"Name", c.Name,
		"City", c.City,
		"Zipcode", c.Zipcode,
		"DateOfBirth", c.DateOfBirth,
		"Status", c.Status)
	RC.Expire(ctx, key, time.Minute)
	//fmt.Println(CustomerResponse)
}

type CustomerResponse struct {
	Id          uint   `redis:"Id"`
	Name        string `redis:"Name"`
	City        string `redis:"City"`
	Zipcode     string `redis:"Zipcode"`
	DateOfBirth string `redis:"DateOfBirth"`
	Status      string `redis:"Status"`
}
