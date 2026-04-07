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

	mux.HandleFunc("POST /v1/books", app.requirePermission("books:write", app.createBookHandler))
	mux.HandleFunc("GET /v1/books/{id}", app.requirePermission("books:read", app.showBookHandler))
	mux.HandleFunc("PATCH /v1/books/{id}", app.requirePermission("books:write", app.updateBookHandler))
	mux.HandleFunc("DELETE /v1/books/{id}", app.requirePermission("books:write", app.deleteBookHandler))
	mux.HandleFunc("GET /v1/books", app.requirePermission("books:read", app.listBooksHandler))

	//*********************//
	// Users endpoints     //
	//*********************//

	// PROTECTED: Only staff can register a new user (librarian creating a citizen profile)
	mux.HandleFunc("POST /v1/users", app.requirePermission("members:write", app.registerUserHandler))
	// PUBLIC: Anyone with an email link can activate their account
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	//*********************//
	// Tokens endpoints    //
	//*********************//
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	//*********************//
	// Members endpoints   //
	//*********************//

	mux.HandleFunc("POST /v1/members", app.requirePermission("members:write", app.createMemberHandler))
	mux.HandleFunc("GET /v1/members/{id}", app.requirePermission("members:read", app.showMemberHandler))
	mux.HandleFunc("PATCH /v1/members/{id}", app.requirePermission("members:write", app.updateMemberHandler))
	mux.HandleFunc("DELETE /v1/members/{id}", app.requirePermission("members:write", app.deleteMemberHandler))

	//*********************//
	// Loans endpoints     //
	//*********************//

	mux.HandleFunc("POST /v1/loans", app.requirePermission("loans:write", app.createLoanHandler))

	//*********************//
	// Fines endpoints     //
	//*********************//

	mux.HandleFunc("POST /v1/fines", app.requirePermission("fines:write", app.createFineHandler))

	return app.authenticate(app.compressResponse(app.recoverPanic(app.enableCORS(app.rateLimit(mux)))))
}