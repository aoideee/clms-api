// Filenaeme: routes.go

package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a new servemux and register the handler functions for the different routes
	mux := http.NewServeMux()

	// Register the handler function for the "/v1/healthcheck" endpoint
	mux.HandleFunc("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The Community Library Management System (CLMS) is up and running!"))
	})

	return mux

}