// Filename: middleware.go

package main

import (
	"net"
	"net/http"
	"sync"
	"compress/gzip"
	"io"
	"strings"


	"golang.org/x/time/rate"
)

// rateLimitMiddleware is a middleware function that limits the number of requests a client can make to the API within a certain time period.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Creates a map to store the rate limiters for each client IP address, and a mutex to protect access to the map.
	var (
		mu sync.Mutex
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

// enableCORS is a middleware function that adds the necessary CORS headers to the response to allow cross-origin requests from trusted origins
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Origin" header to tell browsers that the response
		// may vary based on the Origin header, preventing improper caching
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		// If the Origin header is present and matches one of the trusted origins, add the appropriate CORS headers
		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// If it's a preflight request, add the necessary headers and return a 200 OK response
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
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

// Gzip response writer that wraps the standard http.ResponseWriter and provides gzip compression for the response body.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// compressResponse is a middleware function that compresses the response body using gzip if the client supports it.
func (app *application) compressResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the client doesn't support gzip compression, call the next handler in the chain without modifying the response
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Create the gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Set the header so the browser knows the content is compressed.
		w.Header().Set("Content-Encoding", "gzip")

		// Wrap the response writer and pass it to the next handler
		next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}