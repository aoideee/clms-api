// Filename: internal/data/tokens.go

package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

// Define the different scopes/purposes a token can have in the CLMS
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"` // Never expose the hash to the client
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// generateToken create a cryptographically secure token
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Create a byte slice and fill it with 16 cryptographically secure random bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the 16 bytes into a base-32 string (this results in exactly 26 characters)
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Hash the plaintext token using SHA-256 to store safely in the database
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

// TokenModel wraps the database connection pool.
type TokenModel struct {
	DB *sql.DB
}

// New generates a new token and immediately inserts it into the database.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

// Insert adds the new token hash and details to your PostgreSQL Tokens table.
func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO Tokens (Hash, UserID, Expiry, Scope)
		VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser deletes tokens for a specific user and scope.
// This is used to clean up old tokens once a user logs out or activates their account.
func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM Tokens
		WHERE Scope = $1 AND UserID = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
