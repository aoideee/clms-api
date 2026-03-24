// Filename: internal/data/filters.go

package data

import (
	"time"
	"database/sql"
)

// Fine represents a monetary penalty in our library system.
type Fine struct {
	ID        int64     `json:"id"`
	LoanID    int64     `json:"loan_id"`
	MemberID  int64     `json:"member_id"`
	Amount    float64   `json:"amount"` 
	Status    string    `json:"status"` // e.g., "unpaid" or "paid"
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

// FineModel is the bridge to the fines table in our database.
type FineModel struct {
	DB *sql.DB
}