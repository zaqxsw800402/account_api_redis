package dto

import (
	"red/errs"
	"reflect"
	"testing"
)

func TestNewAccountRequest_Validate(t *testing.T) {
	type fields struct {
		CustomerId  int64
		AccountType string
		Amount      int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *errs.AppError
	}{
		{"testSaving", fields{1, "saving", 6000}, nil},
		{"testSaving", fields{1, "SAVING", 6000}, nil},
		{"testSaving", fields{1, "SAV", 6000},
			errs.NewValidationError("Account type must be saving or checking")},
		{"testChecking", fields{1, "checking", 6000}, nil},
		{"testChecking", fields{1, "CHecking", 6000}, nil},
		{"testAmountLessThan5000", fields{1, "checking", 4000},
			errs.NewValidationError("To open a new account you need to deposit at least 5000")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := AccountRequest{
				CustomerId:  tt.fields.CustomerId,
				AccountType: tt.fields.AccountType,
				Amount:      tt.fields.Amount,
			}
			if got := r.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
