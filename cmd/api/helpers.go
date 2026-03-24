// Filename: helpers.go

package main

import (
	"encoding/json"
	"net/http"
)

// envelope is a generic wrapper for JSON responses. It is used to ensure that all responses have a consistent structure.
type envelope map[string]any

// writeJSON takes the data, formats it as JSON, and writes it to the http.ResponseWriter along with the provided status code and headers.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// MarshalIndent makes the JSON easy for humans to read by adding indentation and newlines.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Adds a new line at the end for terminal readability
	js = append(js, '\n')

	// Add any additional headers to the response
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the status code and the JSON response
	w.WriteHeader(status)
	w.Write(js)

	return nil
}