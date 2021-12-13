package domain

import (
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"strconv"
	"strings"
	"time"
)

type Customer struct {
	Id          uint `gorm:"column:customer_id;primaryKey;autoIncrement"`
	UserID      uint `gorm:"column:user_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string
	Status      string `gorm:"default:1"`
	//Accounts    []Account `gorm:"foreignKey:CustomerId;references:Id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeleteAt  *time.Time `gorm:"column:delete_at;index"`
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

func (c Customer) IsValid() bool {
	date := strings.Split(c.DateOfBirth, "/")
	if len(date) != 3 {
		return false
	}
	year, err := strconv.Atoi(date[0])
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(date[1])
	if err != nil {
		return false
	}
	day, err := strconv.Atoi(date[2])
	if err != nil {
		return false
	}

	switch {
	case len(date[0]) != 4:
		return false
	case len(date[1]) != 2:
		return false
	case len(date[2]) != 2:
		return false
	case year > time.Now().Year():
		return false
	case time.Now().Year()-year > 150:
		return false
	case month > 12:
		return false
	case month < 0:
		return false
	case day > 31:
		return false
	case day < 0:
		return false
	default:
		return true
	}

}
