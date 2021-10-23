package dto

import (
	"red/errs"
	"reflect"
	"testing"
)

func TestTransactionRequest_Validate(t *testing.T) {
	type fields struct {
		AccountId       uint
		Amount          float64
		TransactionType string
		TransactionDate string
		CustomerId      string
	}
	tests := []struct {
		name   string
		fields fields
		want   *errs.AppError
	}{
		// TODO: Add test cases.
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
