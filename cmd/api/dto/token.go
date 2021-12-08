package dto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"red/cmd/api/errs"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type TokenRequest struct {
	PlanText string    `json:"token"`
	UserID   int64     `json:"-"`
	Hash     []byte    `json:"-"`
	Expiry   time.Time `json:"expiry"`
	Scope    string    `json:"-"`
}

type TokenResponse struct {
	PlanText string    `json:"token"`
	UserID   int64     `json:"-"`
	Hash     []byte    `json:"-"`
	Expiry   time.Time `json:"expiry"`
	Scope    string    `json:"-"`
}

func GenerateToken(UserID int, ttl time.Duration, scope string) (*TokenRequest, *errs.AppError) {
	token := &TokenRequest{
		UserID: int64(UserID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, errs.NewNotFoundError("generate token failed" + err.Error())
	}

	token.PlanText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlanText))
	token.Hash = hash[:]

	return token, nil
}
