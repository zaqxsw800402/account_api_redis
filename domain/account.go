package domain

import (
	"red/dto"
	"red/errs"
)

type Account struct {
	AccountId    uint          `gorm:"column:account_id;primaryKey;autoIncrement"`
	CustomerId   uint          `gorm:"column:customer_id"`
	OpeningDate  string        `gorm:"column:opening_date"`
	AccountType  string        `gorm:"column:account_type"`
	Amount       float64       `gorm:"column:amount"`
	Status       string        `gorm:"column:status"`
	Transactions []Transaction `gorm:"foreignKey:AccountId;references:AccountId"`
	//Customers   []Customer `gorm:"foreignKey:Id;references:CustomerId"`
}

func (a Account) ToNewAccountResponseDto() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountId: a.AccountId}
}

type AccountRepository interface {
	Save(account Account) (*Account, *errs.AppError)
	FindBy(id uint) (*Account, *errs.AppError)
	SaveTransaction(t Transaction) (*Transaction, *errs.AppError)
}

func (a Account) CanWithdraw(amount float64) bool {
	if a.Amount < amount {
		return false
	}
	return true
}
