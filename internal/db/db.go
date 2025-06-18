package db

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Config represents database connection parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New establishes a database connection and runs migrations
func New(cfg Config) (*pgx.Conn, error) {
	// Construct database connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Create a context with timeout for connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to database
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create a stdlib database handle from pgx for goose
	// (since goose requires a *sql.DB connection)
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// Create a stdlib database connection for migrations
	stddbConn := stdlib.OpenDB(*connConfig)
	defer stddbConn.Close()

	// Configure connection pool parameters
	stddbConn.SetMaxOpenConns(25)
	stddbConn.SetMaxIdleConns(5)
	stddbConn.SetConnMaxLifetime(5 * time.Minute)

	// Initialize goose and run migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(stddbConn, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database connected and migrations applied successfully")
	return conn, nil
}
