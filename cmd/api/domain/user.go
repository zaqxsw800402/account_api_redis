package domain

import (
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"column:user_id;primaryKey;autoIncrement"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Email     string    `json:"email" gorm:"column:email"`
	Password  string    `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

type UserRepository interface {
	SaveUser(User) (*User, *errs.AppError)
	ByID(string) (*User, *errs.AppError)
	ByEmail(string) (*User, *errs.AppError)
	FindAll() ([]User, *errs.AppError)
	Update(User) (*User, *errs.AppError)

	UpdateToken(Token) (*Token, *errs.AppError)
	SaveToken(Token) (*Token, *errs.AppError)
	GetUserByToken(string) (*User, *errs.AppError)
}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
