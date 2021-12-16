package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

type Database struct {
	RC *redis.Client
}

func New(client *redis.Client) Database {
	return Database{client}
}

func GetClient(host string) (*redis.Client, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "",
		DB:       0,
	})

	Pong, err := client.Ping(ctx).Result()
	log.Println("redis connection " + Pong)
	return client, err
}

type Customer struct {
	Id          uint   `redis:"Id" json:"customer_id"`
	Name        string `redis:"Name" json:"name"`
	City        string `redis:"City" json:"city"`
	DateOfBirth string `redis:"DateOfBirth" json:"date_of_birth"`
	Status      string `redis:"Status" json:"status"`
}

func (d Database) getUserValueForCustomer(customerID string) string {
	return fmt.Sprintf("customer:%s", customerID)
}

func (d Database) getUserKeyForCustomer(userID int) string {
	return fmt.Sprintf("%d:customer", userID)
}

// GetCustomer 讀取使用者資料
func (d Database) GetCustomer(ctx context.Context, customerKey string) (*Customer, error) {
	var customer Customer
	err := d.RC.HGetAll(ctx, customerKey).Scan(&customer)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &customer, nil
}

func (d Database) DeleteCustomer(ctx context.Context, customerID string) error {
	userValue := d.getUserValueForCustomer(customerID)
	result := d.RC.Del(ctx, userValue)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
