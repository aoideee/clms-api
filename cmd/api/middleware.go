// Filename: middleware.go

package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/aoideee/clms-api/internal/data"
	"golang.org/x/time/rate"
)

//*********************//
// Rate limiting       //
//*********************//

// rateLimitMiddleware is a middleware function that limits the number of requests a client can make to the API within a certain time period.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Creates a map to store the rate limiters for each client IP address, and a mutex to protect access to the map.
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforces the limit if enabled in the configuration (main.go)
		if app.config.limiter.enabled {

			// Get the client's IP address from the request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			// Lock the mutex before accessing the clients map
			mu.Lock()

			// Check if a rate limiter already exists for this IP address, and if not, create a new one
			if _, found := clients[ip]; !found {
				clients[ip] = rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)
			}

			// Check if the request is allowed by the rate limiter, and if not, return a 429 Too Many Requests response
			if !clients[ip].Allow() {
				mu.Unlock()

				w.Header().Set("Retry-After", "1")
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		// If the request is allowed, call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

//*********************//
// CORS                //
//*********************//

// enableCORS is a middleware function that adds the necessary CORS headers to the response to allow cross-origin requests from trusted origins
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Origin" header to tell browsers that the response
		// may vary based on the Origin header, preventing improper caching
		vary := w.Header().Get("Vary")
		if vary != "" {
			w.Header().Set("Vary", "Origin, Access-Control-Request-Method, "+vary)
		} else {
			w.Header().Set("Vary", "Origin, Access-Control-Request-Method")
		}

		origin := r.Header.Get("Origin")

		// If the Origin header is present and matches one of the trusted origins, add the appropriate CORS headers
		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// If it's a preflight request, add the necessary headers and return a 200 OK response
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

//*********************//
// Gzip                //
//*********************//

// Gzip response writer that wraps the standard http.ResponseWriter and provides gzip compression for the response body.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.ResponseWriter.Header().Del("Content-Length")
		w.ResponseWriter.WriteHeader(statusCode)
		w.wroteHeader = true
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

//*********************//
// Gzip                //
//*********************//

// compressResponse is a middleware function that compresses the response body using gzip if the client supports it.
func (app *application) compressResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check for gzip support (case-insensitive)
		if !strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Set headers immediately before the handler runs
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		// 3. Wrap the response in our gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

//*********************//
// Panic recovery      //
//*********************//

// recoverPanic is a middleware function that recovers from any panics that occur during the handling of a request and returns a 500 Internal Server Error response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//*********************//
// Authentication      //
//*********************//

// authenticate is a middleware function that authenticates the user based on the Authorization header.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to tell browsers that the response
		// may vary based on the Authorization header, preventing improper caching
		w.Header().Add("Vary", "Authorization")

		// Get the Authorization header from the request
		authorizationHeader := r.Header.Get("Authorization")

		// If there is no Authorization header, treat the user as Anonymous
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Split the Authorization header into the token type and the token itself
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get the token from the Authorization header
		token := headerParts[1]

		// Validate the token to ensure it is in a sensible format
		if len(token) != 26 {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get the user from the database based on the token
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Set the user in the request context
		r = app.contextSetUser(r, user)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// requireAuthenticatedUser checks that a user is not anonymous
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requireActivatedUser checks that a user is both activated and authenticated
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		// Checks if the user is activated
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)

}

// requirePermission checks that a user has the necessary permissions to access a resource
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context.
		user := app.contextGetUser(r)

		// Get the slice of permissions for the user.
		permissions, err := app.models.Permissions.GetAllForUser(user.UserID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Check if the slice includes the required permission. If it doesn't, then
		// return a 403 Forbidden response.
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		// Otherwise they have the necessary permission so we call the next handler in
		// the chain.
		next.ServeHTTP(w, r)
	}

	// Wrap this with the requireActivatedUser middleware before returning it.
	return app.requireActivatedUser(fn)
}
	