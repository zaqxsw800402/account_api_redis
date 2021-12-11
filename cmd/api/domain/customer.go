package domain

import (
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

type Customer struct {
	Id          uint `gorm:"column:customer_id;primaryKey;autoIncrement"`
	UserID      uint `gorm:"column:user_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string
	Status      string    `gorm:"default:1"`
	Accounts    []Account `gorm:"foreignKey:CustomerId;references:Id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

//go:generate mockgen -destination=../mocks/domain/mockCustomerRepository.go -package=domain red/domain CustomerRepository
type CustomerRepository interface {
	FindAll(userID int) ([]Customer, *errs.AppError)
	ByID(string) (*Customer, *errs.AppError)
	Save(customer Customer) (*Customer, *errs.AppError)
	Update(customer Customer) (*Customer, *errs.AppError)
	Delete(string) *errs.AppError
}

// ToDto 將資料更改成需要回傳的格式
func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateOfBirth: c.DateOfBirth,
		Status:      c.statusAsText(),
		//Accounts:    c.Accounts,
	}
}

// statusAsText 轉換資料格式
func (c Customer) statusAsText() string {
	statusAsText := "active"
	if c.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}
