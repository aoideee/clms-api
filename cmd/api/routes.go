// Filename: routes.go

package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a new servemux and register the handler functions for the different routes
	mux := http.NewServeMux()

	//*********************//
	// Default catch-all   //
	//*********************//

	mux.HandleFunc("/", app.notFoundResponse)

	//*********************//
	// Healthcheck endpoint//
	//*********************//

	mux.HandleFunc("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The Community Library Management System (CLMS) is up and running!"))
	})

	//*********************//
	// Observability       //
	//*********************//

	mux.Handle("GET /debug/vars", expvar.Handler())
	mux.Handle("/v1/observability/metrics", expvar.Handler())

	//*********************//
	// Books endpoints     //
	//*********************//

	mux.HandleFunc("POST /v1/books", app.createBookHandler)
	mux.HandleFunc("GET /v1/books/{id}", app.showBookHandler)
	mux.HandleFunc("PATCH /v1/books/{id}", app.updateBookHandler)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBookHandler)
	mux.HandleFunc("GET /v1/books", app.listBooksHandler)

	//*********************//
	// Users endpoints     //
	//*********************//

	// Register a new user (librarian)
	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	// Activate an account (user)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	//*********************//
	// Members endpoints   //
	//*********************//

	mux.HandleFunc("POST /v1/members", app.createMemberHandler)
	mux.HandleFunc("GET /v1/members/{id}", app.showMemberHandler)
	mux.HandleFunc("PATCH /v1/members/{id}", app.updateMemberHandler)
	mux.HandleFunc("DELETE /v1/members/{id}", app.deleteMemberHandler)

	//*********************//
	// Loans endpoints     //
	//*********************//

	mux.HandleFunc("POST /v1/loans", app.createLoanHandler)

	//*********************//
	// Fines endpoints     //
	//*********************//

	mux.HandleFunc("POST /v1/fines", app.createFineHandler)

	return app.compressResponse(app.recoverPanic(app.enableCORS(app.rateLimit(mux))))

}