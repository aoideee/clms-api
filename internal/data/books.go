package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq" // REQUIRED: For handling PostgreSQL arrays like our genres
)

// Book represents a single book record in our library database.
type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Pages     int32     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

// BookModel is the bridge to our database.
type BookModel struct {
	DB *sql.DB // We have handed the model a connection to the database!
}

// Insert adds a new book record into the database.
func (m BookModel) Insert(book *Book) error {
	// The SQL query returning the auto-generated ID, CreatedAt, and Version
	query := `
		INSERT INTO books (title, year, pages, genres) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, version`

	// We use an array adapter for the PostgreSQL text[] array
	args := []any{book.Title, book.Year, book.Pages, pq.Array(book.Genres)}

	// Create a context with a 3-second timeout to prevent the system from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, passing in the context, the query, and the arguments.
	// We use QueryRowContext because we expect exactly one row back (the RETURNING clause).
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}