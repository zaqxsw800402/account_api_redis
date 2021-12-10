package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/cmd/api/errs"
)

type UserRepositoryDb struct {
	client *gorm.DB
}

func NewUserRepositoryDb(dbClient *gorm.DB) UserRepositoryDb {
	return UserRepositoryDb{dbClient}
}

func (u UserRepositoryDb) SaveUser(user User) (*User, *errs.AppError) {
	result := u.client.Create(&user)
	err := result.Error
	if err != nil {
		//logger.Error("Error while creating new user")
		return nil, errs.NewUnexpectedError("unexpected error from database when save user")
	}

	return &user, nil
}

func (u UserRepositoryDb) ByID(id string) (*User, *errs.AppError) {
	var user User
	// 在account表格裡預載入交易紀錄的資料，並且讀取特定id的資料
	result := u.client.Table("users").Where("user_id = ?", id).Find(&user)
	if err := result.Error; err != nil {
		//logger.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("user not found by user_id")
		}
		return nil, errs.NewUnexpectedError("unexpected database error when get user by user_id")
	}
	return &user, nil
}

func (u UserRepositoryDb) ByEmail(id string) (*User, *errs.AppError) {
	var user User
	// 在account表格裡預載入交易紀錄的資料，並且讀取特定id的資料
	result := u.client.Table("users").Where("email = ?", id).Find(&user)
	if err := result.Error; err != nil {
		//logger.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("user not found by email")
		}
		return nil, errs.NewUnexpectedError("unexpected database error when get user by email")
	}
	return &user, nil
}

func (u UserRepositoryDb) FindAll() ([]User, *errs.AppError) {
	var users []User

	result := u.client.Find(&users)
	err := result.Error

	if err != nil {
		//logger.Error("Error while querying users table" + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error when get all users, " + err.Error())
	}

	return users, nil
}

func (u UserRepositoryDb) Update(user User) (*User, *errs.AppError) {
	result := u.client.Model(&user).Updates(User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
	})
	err := result.Error
	if err != nil {
		//logger.Error("Error while creating new user")
		return nil, errs.NewUnexpectedError("unexpected error from database when update user")
	}

	return &user, nil
}
