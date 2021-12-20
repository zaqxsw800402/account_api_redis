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
		DateOfBirth string
		Status      string
	}

	tests := []struct {
		name   string
		fields fields
		want   dto.CustomerResponse
	}{
		{name: "Success", fields: fields{1, "Ivy", "TW", "2012-10-18", "1"},
			want: dto.CustomerResponse{Id: 1, Name: "Ivy", City: "TW", DateOfBirth: "2012-10-18", Status: "active"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Customer{
				Id:          tt.fields.Id,
				Name:        tt.fields.Name,
				City:        tt.fields.City,
				DateOfBirth: tt.fields.DateOfBirth,
				Status:      tt.fields.Status,
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

func TestCustomer_birthday_IsValid(t *testing.T) {

	tests := []struct {
		name string
		date string
		want bool
	}{
		{"SuccessActive", "2006/01/02", true},
		{"SuccessInactive", "2006/1/02", false},
		{"SuccessInactive", "2006/11/2", false},
		{"SuccessInactive", "206/11/02", false},
		{"SuccessInactive", "2006-11/02", false},
		{"SuccessInactive", "2206-11/02", false},
		{"SuccessInactive", "2006-31/02", false},
		{"SuccessInactive", "2006-11/32", false},
		{"SuccessInactive", "1106-11/32", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Customer{
				DateOfBirth: tt.date,
			}
			if got := c.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
