// Filename: internal/data/loans.go

package data

import (
	"context"
	"database/sql"
	"time"
)

// Loan maps perfectly to your transactional Loans table.
type Loan struct {
	ID           int64      `json:"loan_id"`
	CopyID       int64      `json:"copy_id"`
	MemberID     int64      `json:"member_id"`
	CheckoutDate time.Time  `json:"checkout_date"`
	DueDate      time.Time  `json:"due_date"`
	ReturnDate   *time.Time `json:"return_date,omitempty"` // Nullable in DB, so we use a pointer
}

// LoanModel is the bridge to the database for checkout operations.
type LoanModel struct {
	DB *sql.DB
}

// Insert creates a new checkout record in the database.
func (m LoanModel) Insert(loan *Loan) error {
	// We only insert CopyID, MemberID, and DueDate. 
	// CheckoutDate defaults to NOW() in your PostgreSQL schema!
	query := `
		INSERT INTO Loans (CopyID, MemberID, DueDate) 
		VALUES ($1, $2, $3) 
		RETURNING LoanID, CheckoutDate`

	args := []any{loan.CopyID, loan.MemberID, loan.DueDate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&loan.ID, &loan.CheckoutDate)
}