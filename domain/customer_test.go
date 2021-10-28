package domain

import (
	"red/dto"
	"reflect"
	"testing"
)

func TestCustomer_ToDto(t *testing.T) {
	type fields struct {
		Id          uint
		Name        string
		City        string
		Zipcode     string
		DateOfBirth string
		Status      string
		Accounts    []Account
	}
	a := []Account{}

	tests := []struct {
		name   string
		fields fields
		want   dto.CustomerResponse
	}{
		{name: "Success", fields: fields{1, "Ivy", "TW", "23", "2012-10-18", "1", a},
			want: dto.CustomerResponse{Id: 1, Name: "Ivy", City: "TW", Zipcode: "23", DateOfBirth: "2012-10-18", Status: "active", Accounts: a}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Customer{
				Id:          tt.fields.Id,
				Name:        tt.fields.Name,
				City:        tt.fields.City,
				Zipcode:     tt.fields.Zipcode,
				DateOfBirth: tt.fields.DateOfBirth,
				Status:      tt.fields.Status,
				Accounts:    tt.fields.Accounts,
			}
			if got := c.ToDto(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomer_statusAsText(t *testing.T) {

	tests := []struct {
		name   string
		status string
		want   string
	}{
		{"SuccessActive", "1", "active"},
		{"SuccessInactive", "0", "inactive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Customer{
				Status: tt.status,
			}
			if got := c.statusAsText(); got != tt.want {
				t.Errorf("statusAsText() = %v, want %v", got, tt.want)
			}
		})
	}
}
