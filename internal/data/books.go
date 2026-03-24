package data

import (
	"context"
	"database/sql"
	"time"
	"errors"
)

// Book represents a single book record in our library database.
type Book struct {
	ID              int64  `json:"book_id"` // Maps to BookID
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear int32  `json:"publication_year,omitempty"`
	MinimumAge      int32  `json:"minimum_age,omitempty"`
	Description     string `json:"description,omitempty"`
}

// BookModel is the bridge to our database.
type BookModel struct {
	DB *sql.DB // We have handed the model a connection to the database!
}

// Insert adds a new book record into the database.
func (m BookModel) Insert(book *Book) error {
	// The SQL query returning the auto-generated ID, CreatedAt, and Version
	query := `
		INSERT INTO Books (Title, ISBN, Publisher, PublicationYear, MinimumAge, Description) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING BookID`

	// We use an array adapter for the PostgreSQL text[] array
	args := []any{
		book.Title, 
		book.ISBN, 
		book.Publisher, 
		book.PublicationYear, 
		book.MinimumAge, 
		book.Description,
	}

	// Create a context with a 3-second timeout to prevent the system from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, passing in the context, the query, and the arguments.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID)
}

// Get fetches a specific book from the database by its BookID.
func (m BookModel) Get(id int64) (*Book, error) {
	// If the ID is less than 1, we know it's invalid right away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// The exact SELECT query matching your professional ERD
	query := `
		SELECT BookID, Title, ISBN, Publisher, PublicationYear, MinimumAge, Description
		FROM Books
		WHERE BookID = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the result into our Book struct
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.ISBN,
		&book.Publisher,
		&book.PublicationYear,
		&book.MinimumAge,
		&book.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

// Update applies changes to a specific book record in the database.
func (m BookModel) Update(book *Book) error {
	// The SQL query to update the record. We use the BookID to ensure 
	// we only update that specific row.
	query := `
		UPDATE Books 
		SET Title = $1, ISBN = $2, Publisher = $3, PublicationYear = $4, MinimumAge = $5, Description = $6
		WHERE BookID = $7`

	// The arguments must perfectly match the $1, $2, etc. placeholders
	args := []any{
		book.Title,
		book.ISBN,
		book.Publisher,
		book.PublicationYear,
		book.MinimumAge,
		book.Description,
		book.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query. We use ExecContext because we don't need to read any data back,
	// we just need to know if the update was successful.
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// Check if the record actually existed to be updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes a specific book record from the database.
func (m BookModel) Delete(id int64) error {
	// Guard against invalid IDs
	if id < 1 {
		return ErrRecordNotFound
	}

	// The SQL query to delete the record
	query := `
		DELETE FROM Books 
		WHERE BookID = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check if the record actually existed before we tried to delete it
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}