// Filename: main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aoideee/clms-api/internal/data"
	_ "github.com/lib/pq"
)

// config holds all the rules and settings for our application
type config struct {
	port int
	env string 
	db struct {
		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
	limiter struct {
		rps float64
		burst int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
}

// application holds the tools needed to run our library system properly
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	// Server settings
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Database settings
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CLMS_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Rate limiting settings
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Rate limiter maximum burst size")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// CORS settings
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	// Parse the command-line flags
	flag.Parse()

	// Initialize a new logger that writes JSON-formatted logs to the standard output stream
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Log the parsed configuration settings at the INFO level
	logger.Info("Starting application with the following configuration settings",
		"port", cfg.port,
		"env", cfg.env,
		"db_dsn", cfg.db.dsn,
		"limiter_enabled", cfg.limiter.enabled,
		"limiter_rps", cfg.limiter.rps,
		"limiter_burst", cfg.limiter.burst,
		"cors_trusted_origins", cfg.cors.trustedOrigins,
	)

	// Establish the database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Database connection pool established")

	// Initialize the application struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Start the HTTP server
	err = app.serve()
	if err != nil {
		logger.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}

// openDB establishes a connection pool to the PostgreSQL database using the provided configuration settings
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections to the database
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of idle connections in the pool
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Parse the max idle time duration string
	maxIdleTime, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(maxIdleTime)

	// Create a context with a 5-second timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to verify that the connection is successful
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}