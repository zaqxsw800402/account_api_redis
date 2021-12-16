package dto

import (
	"red/errs"
)

const (
	WITHDRAWAL = "withdrawal"
	DEPOSIT    = "deposit"
)

type TransactionRequest struct {
	CustomerId      int64  `json:"customer_id"`
	AccountId       int64  `json:"account_id"`
	TargetAccountId int64  `json:"target_account_id"`
	Amount          int64  `json:"amount"`
	TransactionType string `json:"transaction_type"`
	TransactionDate string `json:"-"`
}

type TransactionResponse struct {
	TransactionId     uint    `json:"transaction_id"`
	AccountId         uint    `json:"account_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	NewBalance        float64 `json:"new_balance"`
	TransactionType   string  `json:"transaction_type"`
	TransactionDate   string  `json:"transaction_date"`
}

// Validate 判斷內容是否有效
func (r TransactionRequest) Validate() *errs.AppError {
	// 判斷有無錯字
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
