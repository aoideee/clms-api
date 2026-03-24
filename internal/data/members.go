// Filename: internal/data/members.go

package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/aoideee/clms-api/internal/validator"
)

// Member perfectly mirrors your professional SQL table design.
type Member struct {
	ID            int64  `json:"member_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	DOB           string `json:"dob"` // Stored as a string in YYYY-MM-DD format for easy JSON parsing
	PhoneNumber   string `json:"phone_number,omitempty"`
	Email         string `json:"email"`
	Address       string `json:"address,omitempty"`
	AccountStatus string `json:"account_status"` // e.g., 'Active', 'Suspended'
}

// ValidateMember runs the strict BNLSIS checks on a new application.
func ValidateMember(v *validator.Validator, member *Member) {
	v.Check(member.FirstName != "", "first_name", "must be provided")
	v.Check(len(member.FirstName) <= 100, "first_name", "must not be more than 100 bytes long")

	v.Check(member.LastName != "", "last_name", "must be provided")
	v.Check(len(member.LastName) <= 100, "last_name", "must not be more than 100 bytes long")

	v.Check(member.DOB != "", "dob", "must be provided")
	// Note: In a production system, we would also parse the DOB to ensure they meet the minimum age 
	// required to act as their own Guarantor!

	v.Check(member.Email != "", "email", "must be provided")
	v.Check(validator.Matches(member.Email, validator.EmailRX), "email", "must be a valid email address")

	v.Check(member.AccountStatus != "", "account_status", "must be provided")
	v.Check(validator.PermittedValue(member.AccountStatus, "Active", "Suspended", "Pending Guarantor"), "account_status", "invalid account status")
}

// MemberModel is the bridge to our database.
type MemberModel struct {
	DB *sql.DB
}

// Insert adds a new member record into the database.
func (m MemberModel) Insert(member *Member) error {
	query := `
		INSERT INTO Member (FirstName, LastName, DOB, PhoneNumber, Email, Address, AccountStatus) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING MemberID`

	args := []any{
		member.FirstName,
		member.LastName,
		member.DOB,
		member.PhoneNumber,
		member.Email,
		member.Address,
		member.AccountStatus,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&member.ID)
}

// Get fetches a specific member from the database.
func (m MemberModel) Get(id int64) (*Member, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT MemberID, FirstName, LastName, DOB, PhoneNumber, Email, Address, AccountStatus
		FROM Member
		WHERE MemberID = $1`

	var member Member
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&member.ID,
		&member.FirstName,
		&member.LastName,
		&member.DOB,
		&member.PhoneNumber,
		&member.Email,
		&member.Address,
		&member.AccountStatus,
	)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &member, nil
}

// Update applies changes to a specific member record.
func (m MemberModel) Update(member *Member) error {
	query := `
		UPDATE Member 
		SET FirstName = $1, LastName = $2, DOB = $3, PhoneNumber = $4, Email = $5, Address = $6, AccountStatus = $7
		WHERE MemberID = $8`

	args := []any{
		member.FirstName,
		member.LastName,
		member.DOB,
		member.PhoneNumber,
		member.Email,
		member.Address,
		member.AccountStatus,
		member.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Delete removes a member record from the database.
func (m MemberModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM Member WHERE MemberID = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}