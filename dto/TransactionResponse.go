package dto

type TransactionResponse struct {
	TransactionId   uint    `json:"transaction_id"`
	AccountId       uint    `json:"account_id"`
	Amount          float64 `json:"new_balance"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}
