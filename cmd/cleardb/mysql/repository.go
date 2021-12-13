package mysql

import (
	"gorm.io/gorm"
	"red/cmd/api/errs"
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

func (d DB) DeleteCustomers() *errs.AppError {
	result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Delete(&Customer{})
	if err := result.Error; err != nil {
		return errs.NewUnexpectedError("Unexpected database error when delete customer")
	}

	return nil
}

func (d DB) DeleteAccounts() ([]Account, *errs.AppError) {
	//result := d.client.Where("status = 0 and timediff(delete_at, ?) > ? ",time.Now(), time.Duration(time.Second)).Delete(&Account{}, accountID)
	//result := d.client.Where("status = 0 and delete_at <  ? ", time.Now().Add(-1*time.Second)).Delete(&Account{}, accountID)
	var accounts []Account

	result := d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Find(&accounts)
	if err := result.Error; err != nil {
		return nil, errs.NewUnexpectedError("Unexpected database error when delete account")
	}
	result = d.client.Where("status = 0 and deleted_at <  ? ", time.Now().Add(-1*time.Second)).Delete(&Account{})
	if err := result.Error; err != nil {
		return nil, errs.NewUnexpectedError("Unexpected database error when delete account")
	}

	return accounts, nil
}

func (d DB) DeleteTransactions(accountID uint) *errs.AppError {
	result := d.client.Where("account_id = ? ", accountID).Delete(&Transaction{})
	if err := result.Error; err != nil {
		return errs.NewUnexpectedError("Unexpected database error when delete transactions")
	}

	return nil
}
