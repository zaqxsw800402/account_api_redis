package dto

import (
	"red/errs"
	"reflect"
	"testing"
)

func TestTransactionRequest_Validate(t *testing.T) {
	type fields struct {
		AccountId       int64
		Amount          int64
		TransactionType string
		TransactionDate string
		CustomerId      int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *errs.AppError
	}{
		{"SuccessDeposit", fields{1, 6000, "deposit", "", 1}, nil},
		{"SuccessWithdraw", fields{1, 6000, "withdrawal", "", 1}, nil},
		{"FailedAmount", fields{1, -4000, "withdrawal", "", 1},
			errs.NewValidationError("Amount cannot be less than zero")},
		{"FailedTransactionType", fields{1, 6000, "wi", "", 1},
			errs.NewValidationError("transaction type can only be deposit or withdrawal")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := TransactionRequest{
				AccountId:       tt.fields.AccountId,
				Amount:          tt.fields.Amount,
				TransactionType: tt.fields.TransactionType,
				TransactionDate: tt.fields.TransactionDate,
				CustomerId:      tt.fields.CustomerId,
			}
			if got := r.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
