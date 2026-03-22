// Filename: middleware.go

package main

import (
	"net"
	"net/http"
	"sync"

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