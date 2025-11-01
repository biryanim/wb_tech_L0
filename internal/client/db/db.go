package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Handler is a function type that performs database operations within a transaction context.
type Handler func(ctx context.Context) error

// SQLExecer defines methods for executing SQL queries and commands.
type SQLExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// Transactor defines methods for managing database transactions.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// Pinger defines a method for checking database connectivity.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB combines SQL execution, transaction management, and connection health checking.
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}

// Client provides access to the database and manages its lifecycle.
type Client interface {
	DB() DB
	Close() error
}

// TxManager handles transaction management with specified isolation levels.
type TxManager interface {
	ReadCommited(cxt context.Context, f Handler) error
}
