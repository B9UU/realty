package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/b9uu/realty/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
	ScopeActivation     = "activation"
)

type Token struct {
	Plaintext string    `json:"plaintext"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

type TokenInterface interface {
	New(userId int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
}

func generateNewToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	// generate random slice of bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	// encode the byte slice to a base-32-encoded string
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// generate sha-256 we will save in db
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

// generates new token and insert to tokens table
func (t TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateNewToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	if err := t.Insert(token); err != nil {
		return nil, err
	}
	return token, nil

}

func (t TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	args := []interface{}{
		&token.Hash, &token.UserID, &token.Expiry, &token.Scope,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil

}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

// validate token
func ValidateTokenPlainText(v *validator.Validator, plainText string) {
	v.Check(plainText != "", "token", "token must be provided")
	v.Check(len(plainText) == 26, "token", "must be 26 bytes long")
}
