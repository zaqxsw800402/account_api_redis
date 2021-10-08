package domain

import "red/dto"

type Transaction struct {
	TransactionId   string `gorm:"primaryKey;autoIncrement"`
	AccountId       string
	Amount          float64
	TransactionType string
	TransactionDate string
	Accounts        []Account `gorm:"foreignKey:AccountId;references:AccountId"`
}

func (t Transaction) IsWithdrawal() bool {
	if t.TransactionType == "withdrawal" {
		return true
	}
	return false
}

func (t Transaction) ToDto() dto.TransactionResponse {
	return dto.TransactionResponse{
		TransactionId:   t.TransactionId,
		AccountId:       t.AccountId,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionDate: t.TransactionDate,
	}
}
