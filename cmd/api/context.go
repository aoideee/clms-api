// Filename: cmd/api/context.go

package main

import (
	"context"
	"net/http"

	"github.com/aoideee/clms-api/internal/data"
)

// Define a custom type to prevent collissions with other packages that might also store data in the request context
type contextKey string

const userContextKey = contextKey("user")

// contextSetUser creates a new copy of the request with the provided User struct added to the context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrievs the User struct from the request context
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}