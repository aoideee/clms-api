// Filename: cmd/api/handlers_users.go

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/aoideee/clms-api/internal/data"
	"github.com/aoideee/clms-api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Expected JSON paylod from the librarian (No password field)
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy data into the User struct. Activated defaults to false.
	user := &data.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Role:      input.Role,
		Activated: false,
	}

	// Validate the input data
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user into the database
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate a secure activation token valid for 3 days
	token, err := app.models.Token.New(user.UserID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Log the token to the terminal (for testing purposes)
	app.logger.Info("user registered and activation token generated",
		"activation_token", token.Plaintext,
		"email", user.Email,
	)

	// Return a 201 Created response to the user
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Expected JSON input from the user
	var input struct {
		TokenPlaintext string `json:"token"`
		Password       string `json:"password"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the token and the new password
	v := validator.New()
	v.Check(input.TokenPlaintext != "", "token", "must be provided")
	v.Check(len(input.TokenPlaintext) == 26, "token", "must be exactly 26 characters long")
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Find the user associated with the token
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Set the new password and activate the user's account
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Activated = true

	// Save the updated user record to the database
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Delete the activation token from teh database so it can never be used again
	err = app.models.Token.DeleteAllForUser(data.ScopeActivation, user.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the user a welcome JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
