package domain

import (
	"red/dto"
	"time"
)

type Token struct {
	UserID    int64     `gorm:"column:user_id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Hash      []byte    `gorm:"column:token_hash"`
	Expiry    time.Time `gorm:"column:expiry"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	PlanText  string    `gorm:"-"`
}

func (t Token) ToDto() dto.TokenResponse {
	return dto.TokenResponse{
		PlanText: t.PlanText,
		UserID:   t.UserID,
		Hash:     t.Hash,
		Expiry:   t.Expiry,
		Scope:    dto.ScopeAuthentication,
	}
}
