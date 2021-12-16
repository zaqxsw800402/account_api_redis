package domain

import (
	"red/cmd/api/dto"
	"reflect"
	"testing"
)

func TestAccount_CanWithdraw(t *testing.T) {
	type fields struct {
		Amount float64
	}

	tests := []struct {
		name   string
		fields fields
		amount float64
		want   bool
	}{
		{"Success", fields{50000.00}, 7000, true},
		{"failed", fields{500}, 7000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Account{
				Amount: tt.fields.Amount,
			}
			if got := a.CanWithdraw(tt.amount); got != tt.want {
				t.Errorf("CanWithdraw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_ToNewAccountResponseDto(t *testing.T) {
	type fields struct {
		AccountId    uint
		CustomerId   uint
		OpeningDate  string
		AccountType  string
		Amount       float64
		Status       string
		Transactions []Transaction
	}
	tests := []struct {
		name   string
		fields fields
		want   dto.AccountResponse
	}{
		{"Success_1", fields{1, 1, "2012-10-18", "saving", 5000, "status", nil},
			dto.AccountResponse{AccountId: 1}},
		{"Success_2", fields{2, 2, "2012-10-18", "saving", 5000, "status", nil},
			dto.AccountResponse{AccountId: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Account{
				AccountId:    tt.fields.AccountId,
				CustomerId:   tt.fields.CustomerId,
				OpeningDate:  tt.fields.OpeningDate,
				AccountType:  tt.fields.AccountType,
				Amount:       tt.fields.Amount,
				Status:       tt.fields.Status,
				Transactions: tt.fields.Transactions,
			}
			if got := a.ToNewAccountResponseDto(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNewAccountResponseDto() = %v, want %v", got, tt.want)
			}
		})
	}
}
