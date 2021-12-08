package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"red/cmd/api/domain"
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

type TokenService interface {
	UpdateToken(string) (*dto.TokenResponse, *errs.AppError)
	SaveToken(dto.TokenRequest) (*dto.TokenResponse, *errs.AppError)
	GetUser(string) (*dto.TokenResponse, *errs.AppError)
}

type DefaultTokenService struct {
	repo domain.TokenRepository
}

func (t DefaultTokenService) UpdateToken(req dto.TokenRequest) (*dto.TokenResponse, *errs.AppError) {
	t.repo.UpdateToken()
}

func GenerateToken(UserID int, ttl time.Duration, scope string) (*Token, error) {
	token := &dto.TokenRequest{
		UserID: int64(UserID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlanText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlanText))
	token.Hash = hash[:]

	return token, nil
}
