package domain

import (
	"red/cmd/api/dto"
)

type Transaction struct {
	TransactionId   uint `gorm:"primaryKey;autoIncrement"`
	AccountId       uint
	Amount          float64
	TransactionType string
	TransactionDate string
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
