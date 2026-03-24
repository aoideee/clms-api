package data

import (
	"time"
)

// Book represents a single book record in the library database.
type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // The hyphen tells Go to hide this from the public JSON
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Pages     int32     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

// BookModel is the structural bridge between Go code and the PostgreSQL database.
type BookModel struct {
	// Will add the database connection pool here very soon!
}