package domain

import (
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

type User struct {
	ID        int       `gorm:"column:user_id;primaryKey;autoIncrement"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email;unique"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeleteAt  time.Time `gorm:"column:delete_at"`
}

type UserRepository interface {
	SaveUser(User) (*User, *errs.AppError)
	ByID(string) (*User, *errs.AppError)
	ByEmail(string) (*User, *errs.AppError)
	FindAll() ([]User, *errs.AppError)
	Update(User) (*User, *errs.AppError)

	UpdateToken(Token) (*Token, *errs.AppError)
	SaveToken(Token) (*Token, *errs.AppError)
	GetUserWithToken(string) (*User, *errs.AppError)
}

func (u User) ToDto() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		//Password:  u.Password,
		//CreatedAt: u.CreatedAt,
		//UpdatedAt: u.UpdatedAt,
	}
}
