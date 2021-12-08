package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

func GenerateToken(UserID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
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

type Token struct {
	PlanText string    `json:"token"`
	UserID   int64     `json:"-"`
	Hash     []byte    `json:"-"`
	Expiry   time.Time `json:"expiry"`
	Scope    string    `json:"-"`
}

type TokenDB struct {
	UserID    int64     `gorm:"column:user_id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Hash      string    `gorm:"column:token_hash"`
	Expiry    time.Time `gorm:"expiry"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx := m.DB.WithContext(ctx)

	stmt := `delete from tokens where user_id = ?`
	result := tx.Exec(stmt, u.ID)
	if result.Error != nil {
		return result.Error
	}

	stmt = `insert into tokens (user_id, name, email, token_hash, expiry, created_at, updated_at)
			values (?,?,?,?,?,?,?)`

	result = tx.Exec(stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		t.Expiry,
		time.Now(),
		time.Now(),
	)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx := m.DB.WithContext(ctx)

	tokenHash := sha256.Sum256([]byte(token))
	var user User

	query := `
		select
			u.id, u.first_name, u.last_name, u.email
		from
			users u
			inner join tokens t on (u.id = t.user_id)
		where
			t.token_hash = ?
			and t.expiry > ?
	`

	result := tx.Exec(query, tokenHash[:], time.Now()).Scan(&user)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return &user, nil
}
