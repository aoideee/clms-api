// Filename: internal/data/users.go

package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/aoideee/clms-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// ErrDuplicateEmail is a custom error raised when the email address already exists in the database.
var ErrDuplicateEmail = errors.New("duplicate email")

// ErrEditConflict is a custom error raised when the version number doesn't match
var ErrEditConflict = errors.New("edit conflict")

type User struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  password  `json:"-"` // The "-" prevents the hash from EVER appearing in JSON responses
	Role      string    `json:"role"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

//*********************//
// Password methods    //
//*********************//

// password struct to hold the plaintext and hash of the password
type password struct {
	plaintext *string
	hash      []byte
}

// Set calculates the bcrypt hash of a plaintext password and updates the struct
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Matches checks whether the provided plaintext password matches the hashed password
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

//*********************//
// Validation methods  //
//*********************//

// ValidateEmail checks that the email address is not empty and matches the regex pattern
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext ensures the password meets the strict length requirements
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(len(password) <= 72, "password", "must be less than 72 characters long")
}

// ValidateUser acts as a wrapper to validate all individual fields of a User struct
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(len(user.FirstName) <= 100, "first_name", "must not be more than 100 characters long")

	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.LastName) <= 100, "last_name", "must not be more than 100 characters long")

	// Validae the email using the helper function
	ValidateEmail(v, user.Email)

	// Validate the plain text password
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
}

//*********************//
// Database methods    //
//*********************//

// UserModel wraps the database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert a new user into the database
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role, Activated)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING UserID, CreatedAt, Version
	`

	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password.hash,
		user.Role,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UserID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		// If the email already exists, PostgreSQL throws a unique constraint error
		case err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// GetByEmail retrives a user's details based on their email address
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT UserID, CreatedAt, FirstName, LastName, Email, PasswordHash, Role, Activated, Version
		FROM Users
		WHERE Email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update updates the details for a specific user
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE Users
		SET FirstName = $1, LastName = $2, Email = $3, PasswordHash = $4, Role = $5, Activated = $6, Version = Version +1
		WHERE UserID = $7 AND Version = $8
		RETURNING Version`

	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password.hash,
		user.Role,
		user.Activated,
		user.UserID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		// Checks if the updated email conflicts with another user
		case err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

//*********************//
// Token methods       //
//*********************//

// GetForToken retieves a user from the database using a specific token and scope
func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the hash for the plaintext token provided by the client
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Uses an INNER JOIN to fetch the User details based on a match in the Tokens table
	query := `
			SELECT Users.UserID, Users.CreatedAt, Useres.FirstName, Users.LastName, Users.Email, Users.PasswordHash, Users.Role, Users.Activated, Users.Version
			FROM Users
			INNER JOIN Tokens ON Users.UserID = Tokens.UserID
			WHERE Tokens.Hash = $1
			AND Tokens.Scope = $2
			AND Tokens.Expirey > $3`

	// Pass tokenHash[:] to convert the array to a slice
	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}
