package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. Will be returned from the models when looking up a record that doesn't exist in the database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models struct wraps all the individual data models so it can be easily passed around our application. 
type Models struct {
	Books BookModel
	Fines FineModel
	Loans LoanModel
	Members MemberModel
}

// NewModels returns a fully populated Models struct, initializing each individual model with the database connection pool.
func NewModels(db *sql.DB) Models {
	return Models{
		Books: BookModel{DB: db},
		Fines: FineModel{DB: db},
		Loans: LoanModel{DB: db},
		Members: MemberModel{DB: db},
	}
}