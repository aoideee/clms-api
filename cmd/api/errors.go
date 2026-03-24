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

// serverErrorResponse handles sudden library collases (500 errors)
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