// Filename: errors.go

package main

import (
	"net/http"
)

// logError securely records the problem in the backend ledger (the JSON logger)
func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
}

// errorResponse is the foundational helper for sending any JSON error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse handles sudden library collapses (500 errors)
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// rateLimitExceededResponse handles 429 Too Many Requests errors
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "You have exceeded your rate limit. Please try again later."
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// notFoundResponse handles 404 Not Found errors
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// failedValidationResponse sends a 422 Unprocessable Entity status code and the map of validation errors.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// editConflictResponse handles 409 Conflict errors
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "Unable to update the record due to a conflict with another user's changes. Please try again."
	app.errorResponse(w, r, http.StatusConflict, message)
}

// badRequestResponse handles 400 Bad Request errors
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	message := "The request could not be understood by the server due to malformed syntax."
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

// invalidCredentialsResponse handles 401 Unauthorized errors
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid authentication credentials."
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// invalidAuthenticationTokenResponse handles 401 Unauthorized errors
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid or missing authentication token."
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// authenticationRequiredResponse handles user log in errors
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// inactiveAccountResponse handles unactivated accounts
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}