// Filename: internal/data/loans.go

package data

import (
	"time"
)

// Loan tracks the borrowing of a book by a member.
type Loan struct {
	ID           int64      `json:"id"`
	BookID       int64      `json:"book_id"`
	MemberID     int64      `json:"member_id"`
	CheckoutDate time.Time  `json:"checkout_date"`
	DueDate      time.Time  `json:"due_date"`
	// We use a pointer for ReturnDate (*time.Time) because it might be empty (nil) 
	// if the book has not been returned yet!
	ReturnDate   *time.Time `json:"return_date,omitempty"` 
	Status       string     `json:"status"` // e.g., "active", "returned", "overdue"
	Version      int32      `json:"version"`
}

// LoanModel is the bridge to the loans table in our database.
type LoanModel struct {
	// The database connection pool will be added here
}