package domain

import (
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

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
	//Transactions []Transaction `gorm:"foreignKey:AccountId;references:AccountId" `
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain red/domain AccountRepository
type AccountRepository interface {
	// Save Creates a new account
	Save(account Account) (*Account, *errs.AppError)
	// ByID Get Account by customer id and account id
	ByID(accountID uint) (*Account, *errs.AppError)
	// TransactionsByID Get all transactions by account id
	TransactionsByID(accountID uint) ([]Transaction, *errs.AppError)
	SaveTransaction(t Transaction) (*Transaction, *errs.AppError)
	Update(account Account) (*Account, *errs.AppError)
	// ByCustomerID Get all accounts by customer id
	ByCustomerID(id uint) ([]Account, *errs.AppError)
	Delete(accountID string) *errs.AppError
	// ByUserID  Get all accounts by user id
	//ByUserID(id int) ([]Account, *errs.AppError)
}

// ToNewAccountResponseDto 轉換成回傳的json格式
func (a Account) ToNewAccountResponseDto() dto.AccountResponse {
	return dto.AccountResponse{AccountId: a.AccountId}
}

// CanWithdraw 判斷是否能取出金錢
func (a Account) CanWithdraw(amount float64) bool {
	if a.Amount < amount {
		return false
	}
	return true
}

func (a Account) ToDto() dto.AccountResponse {
	return dto.AccountResponse{
		AccountId:   a.AccountId,
		AccountType: a.AccountType,
		Amount:      a.Amount,
		Status:      a.statusAsText(),
	}
}

// statusAsText 轉換資料格式
func (a Account) statusAsText() string {
	statusAsText := "active"
	if a.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

func (a Account) InactiveAccount() bool {
	if a.Status == "0" {
		return false
	}
	return true
}
