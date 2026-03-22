// Filename: server.go

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// Configire the HTTP server with robust timeouts
	srv := &http.Server{
		Addr: 	   fmt.Sprintf(":%d", app.config.port),
		// Temporrily use the default HTTP mux until we set up our own
		Handler:   http.DefaultServeMux,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a channel to receive any errors returned by the graceful Shutdown() function
	shutdownError := make(chan error)

	// Start a background goroutine to listen for OS signals for shutdown
	go func() {
		// Create a channel to listen for shutdown signals
		quit := make(chan os.Signal, 1)

		// Listen for SIGINT (Ctrl+C) and SIGTERM (standard termination signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// The goroutine will block here until a signal is received
		s := <-quit

		// Log that a shutdown signal has been received
		app.logger.Info("Shutting down server", "signal", s.String())

		// Add a 3-second delay before initiating the actual shutdown process
		time.Sleep(3 * time.Second)

		// Create a context with a 5-second timeout
		// This gives in-flight requests a maximum of 5 seconds to complete before the server is forcefully closed
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the Shutdown() method on the server, passing in the context
		shutdownError <- srv.Shutdown(ctx)
	}()

	// Start the server in the background
	app.logger.Info("Starting server", "addr", srv.Addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait to receive the return value from Shutdown() 
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("Server stopped gracefully")

	return nil
}