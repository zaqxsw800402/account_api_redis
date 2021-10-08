package domain

import (
	"red/dto"
	"red/errs"
)

type Customer struct {
	Id          uint `gorm:"column:customer_id;primaryKey;autoIncrement"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string
	Status      string    `gorm:"default:1"`
	Accounts    []Account `gorm:"foreignKey:CustomerId;references:Id"`
}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateOfBirth: c.DateOfBirth,
		Status:      c.statusAsText(),
		Accounts:    c.Accounts,
	}
}

func (c Customer) statusAsText() string {
	statusAsText := "active"
	if c.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
	Save(customer Customer) (*Customer, *errs.AppError)
}
