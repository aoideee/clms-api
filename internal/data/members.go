// Filename: internal/data/members.go

package data

import (
	"time"
)

// Member represents an individual who holds a library card in our system.
type Member struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Version   int32     `json:"version"`
}

// MemberModel is the bridge to the members table in our database.
type MemberModel struct {
	// The database connection pool will be added here
}