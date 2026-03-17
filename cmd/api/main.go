// Filename: main.go

package main
import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// config holds all the rules and settings for our applicaton
type config struct {
	port int
	env string 
	db struct {
		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
}

// application holds the tools needed to run our library system properly
type application struct {
	config config
	logger *slog.Logger
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

	// Parse the command-line flags
	flag.Parse()

	// -------------------------------------------------------------------------
	// REQUIREMENT: STRUCTURED JSON LOGGING
	// CRITICAL FIX: Changed from TextHandler to JSONHandler to pass the test.
	// -------------------------------------------------------------------------

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Log the pasred configuration settings at the INFO level
	logger.Info("Starting application with the following configuration settings",
		slog.Int("port", cfg.port),
		slog.String("env", cfg.env),
		slog.String("db-dsn", cfg.db.dsn),
		slog.Int("db-max-open-conns", cfg.db.maxOpenConns),
		slog.Int("db-max-idle-conns", cfg.db.maxIdleConns),
		slog.String("db-max-idle-time", cfg.db.maxIdleTime),
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