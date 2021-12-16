package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type DBModel struct {
	DB *gorm.DB
}

type Models struct {
	DB DBModel
}

type User struct {
	ID        int       `gorm:"column:user_id;primaryKey;autoIncrement"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Session struct {
	Token string `gorm:"column:token;primaryKey"`
	Data  []byte `gorm:"data;not null"`
	//Expiry int64 `gorm:"autoCreateTime"`
	Expiry time.Time `gorm:"column:expiry;not null"`
}

func (m *DBModel) Authenticate(email, password string) (int, error) {

	var user User
	result := m.DB.Table("users").Where("email = ?", email).Find(&user)
	if err := result.Error; err != nil {
		return user.ID, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return user.ID, nil

}
