// File: internal/data/permissions.go

package data

import (
	"context"
	"database/sql"
	"time"
)

// Permissions is a custom slice type to hold the permission codes
type Permissions []string

// Include is a helper methed to check whether the Permissions slice contains a specific code
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

// PermissionModel wraps the database connection pool
type PermissionModel struct {
	DB *sql.DB
}

// GetAllForUser returns all permission codes for a specific user in a Permission slice
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	// Uses an INNER JOIN to fetch the actual permission codes based on the linking table
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		WHERE users_permissions.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
