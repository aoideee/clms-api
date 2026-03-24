// Filename: internal/data/filters.go

package data

import (
	"context"
	"database/sql"
	"time"
)

// Fine maps to your database and handles both overdue fines and BNLSIS membership deposits.
type Fine struct {
	ID         int64   `json:"fine_id"`
	LoanID     *int64  `json:"loan_id,omitempty"` // Nullable if this is a general membership fee!
	MemberID   int64   `json:"member_id"`
	FineType   string  `json:"fine_type"` // e.g., 'Overdue', 'Local Membership Fee ($3)', 'Foreigner Deposit ($43)'
	Amount     float64 `json:"amount"`
	PaidStatus bool    `json:"paid_status"`
}

type FineModel struct {
	DB *sql.DB
}

// Insert creates a new fee or fine on a patron's account.
func (m FineModel) Insert(fine *Fine) error {
	query := `
		INSERT INTO Fine (LoanID, MemberID, FineType, Amount, PaidStatus) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING FineID`

	args := []any{fine.LoanID, fine.MemberID, fine.FineType, fine.Amount, fine.PaidStatus}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&fine.ID)
}