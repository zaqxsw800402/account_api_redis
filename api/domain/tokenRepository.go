package domain

import (
	"crypto/sha256"
	"database/sql"
	"red/errs"
	"time"
)

func (u UserRepositoryDb) UpdateToken(token Token) (*Token, *errs.AppError) {
	result := u.client.Model(&token).Updates(Token{
		UserID:    token.UserID,
		Name:      token.Name,
		Email:     token.Email,
		Hash:      token.Hash,
		Expiry:    token.Expiry,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err := result.Error; err != nil {
		//logger.Error("Error while creating new user")
		return nil, errs.NewUnexpectedError("unexpected error from database when update token")
	}
	return &token, nil
}

func (u UserRepositoryDb) GetUserWithToken(hash string) (*Token, *errs.AppError) {
	tokenHash := sha256.Sum256([]byte(hash))

	//var user User
	//result := u.client.Table("users").Select("users.user_id, users.first_name, users.last_name, users.email").
	//	Joins("inner join tokens on users.user_id = tokens.user_id").Where("tokens.token_hash = ? and tokens.expiry > ?", tokenHash[:], time.Now()).Scan(&user)

	var token Token
	result := u.client.Where("token_hash = ? and expiry > ?", tokenHash[:], time.Now()).First(&token)

	if err := result.Error; err != nil {
		//logger.Error("Error while creating new user")
		return nil, errs.NewUnexpectedError("unexpected error from database when get user with token")
	}
	return &token, nil
}

func (u UserRepositoryDb) SaveToken(token Token) (*Token, *errs.AppError) {
	result := u.client.Delete(&token)
	if result.Error != nil && result.Error != sql.ErrNoRows {
		return nil, errs.NewUnexpectedError("unexpected error from database when save token")
	}

	result = u.client.Create(&token)
	if err := result.Error; err != nil {
		//logger.Error("Error while creating new user")
		return nil, errs.NewUnexpectedError("unexpected error from database when save token")
	}

	return &token, nil
}
