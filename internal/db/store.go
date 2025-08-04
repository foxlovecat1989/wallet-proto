package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// Store interface for database operations
type Store interface {
	Close() error
	DB() *sqlx.DB
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	// SQLx specific methods
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

// store implements Store
type store struct {
	db *sqlx.DB
}

// NewStore creates a new store
func NewStore(cfg *DatabaseConfig) (Store, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &store{db: db}, nil
}

// Close closes the database connection
func (d *store) Close() error {
	return d.db.Close()
}

// DB returns the underlying sqlx.DB instance
func (d *store) DB() *sqlx.DB {
	return d.db
}

// QueryRowContext executes a query that returns a single row
func (d *store) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// QueryContext executes a query that returns multiple rows
func (d *store) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

// ExecContext executes a query that doesn't return rows
func (d *store) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// BeginTx starts a new transaction
func (d *store) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return d.db.BeginTxx(ctx, opts)
}

// GetContext executes a query that returns a single row and scans it into dest
func (d *store) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.GetContext(ctx, dest, query, args...)
}

// SelectContext executes a query that returns multiple rows and scans them into dest
func (d *store) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.SelectContext(ctx, dest, query, args...)
}

// NamedExecContext executes a named query that doesn't return rows
func (d *store) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return d.db.NamedExecContext(ctx, query, arg)
}

// NamedQueryContext executes a named query that returns rows
func (d *store) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return d.db.NamedQueryContext(ctx, query, arg)
}
