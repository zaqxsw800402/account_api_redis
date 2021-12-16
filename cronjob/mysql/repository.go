package mysql

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	Id          uint `gorm:"column:customer_id;primaryKey;autoIncrement"`
	UserID      uint `gorm:"column:user_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string
	Status      string     `gorm:"default:1"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeleteAt    *time.Time `gorm:"column:deleted_at;index"`
}

type Account struct {
	AccountId   uint       `gorm:"column:account_id;primaryKey;autoIncrement" `
	CustomerId  uint       `gorm:"column:customer_id" `
	OpeningDate string     `gorm:"column:opening_date" `
	AccountType string     `gorm:"column:account_type" `
	Amount      float64    `gorm:"column:amount" `
	Status      string     `gorm:"column:status;default:1" `
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeleteAt    *time.Time `gorm:"column:deleted_at;index"`
}

type Transaction struct {
	TransactionId   uint `gorm:"primaryKey;autoIncrement"`
	AccountId       uint
	Amount          float64
	NewBalance      float64
	TransactionType string
	TransactionDate string
}

func NewDB(dbClient *gorm.DB) DB {
	return DB{dbClient}
}

type DB struct {
	client *gorm.DB
}

func (d DB) DeleteCustomers() ([]Customer, error) {
	var customers []Customer
	result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Find(&customers)
	//result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-7*24*time.Hour)).Find(&customers)
	if err := result.Error; err != nil {
		return nil, errors.New("unexpected database error when delete customer " + err.Error())
	}

	result = d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Delete(&Customer{})
	//result = d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-7*24*time.Hour)).Delete(&Customer{})
	if err := result.Error; err != nil {
		return nil, errors.New("unexpected database error when delete customer " + err.Error())
	}

	return customers, nil
}

func (d DB) DeleteAccounts() ([]Account, error) {
	var accounts []Account

	result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Find(&accounts)
	//result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-7*24*time.Hour)).Find(&accounts)
	if err := result.Error; err != nil {
		return nil, errors.New("unexpected database error when delete account " + err.Error())
	}

	result = d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Delete(&Account{})
	//result = d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-7*24*time.Hour)).Delete(&Account{})
	if err := result.Error; err != nil {
		return nil, errors.New("unexpected database error when delete account " + err.Error())
	}

	return accounts, nil
}

func (d DB) DeleteTransactions(accountID uint) error {
	result := d.client.Where("account_id = ? ", accountID).Delete(&Transaction{})
	if err := result.Error; err != nil {
		return errors.New("unexpected database error when delete transactions " + err.Error())
	}

	return nil
}
