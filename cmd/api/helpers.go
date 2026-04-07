// Filename: helpers.go

package main

import (
	"encoding/json"
	"fmt"
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

// readJSON securely decodes the incoming JSON request body into a target destination.
func (app *application) readJSON(r *http.Request, dst any) error {
	// Initialize the json.Decoder, and call DisallowUnknownFields() on it before decoding.
	// This ensures that if the client sends fields that do not exist in our target struct,
	// the request will be rejected rather than silently ignoring the extra data.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

// background runs an arbitrary function in a background goroutine.
func (app *application) background(fn func()) {
	// Increment the WaitGroup counter.
	app.wg.Add(1)

	// Launch a background goroutine.
	go func() {
		// Decrement the WaitGroup counter when the goroutine completes.
		defer app.wg.Done()

		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("background task panic", "error", fmt.Errorf("%s", err))
			}
		}()

		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}