package dto

import (
	"red/errs"
)

type TransactionRequest struct {
	AccountId       uint    `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
	CustomerId      string  `json:"-"`
}

const WITHDRAWAL = "withdrawal"
const DEPOSIT = "deposit"

func (r TransactionRequest) Validate() *errs.AppError {
	if !r.IsTransactionTypeWithdrawal() && !r.IsTransactionTypeDeposit() {
		return errs.NewValidationError("transaction type can only be deposit or withdrawal")
	}

	if r.Amount < 0 {
		return errs.NewValidationError("Amount cannot be less than zero")
	}
	return nil
}

func (r TransactionRequest) IsTransactionTypeWithdrawal() bool {
	return r.TransactionType == WITHDRAWAL
}

func (r TransactionRequest) IsTransactionTypeDeposit() bool {
	return r.TransactionType == DEPOSIT
}
