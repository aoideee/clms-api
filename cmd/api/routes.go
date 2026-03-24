// Filename: routes.go

package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a new servemux and register the handler functions for the different routes
	mux := http.NewServeMux()

	// Register the handler function for the "/v1/healthcheck" endpoint
	mux.HandleFunc("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The Community Library Management System (CLMS) is up and running!"))
	})

	mux.Handle("/v1/observability/metrics", expvar.Handler())

	// Register the handler functions for the "/v1/books" endpoints
	mux.HandleFunc("POST /v1/books", app.createBookHandler)
	mux.HandleFunc("GET /v1/books/{id}", app.showBookHandler)
	mux.HandleFunc("PATCH /v1/books/{id}", app.updateBookHandler)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBookHandler)

	return app.enableCORS(app.rateLimit(app.compressResponse(mux)))

}