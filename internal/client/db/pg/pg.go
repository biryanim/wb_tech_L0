package pg

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type key string

const (
	// TxKey is the context key used to store transaction information.
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB creates and returns a new PostgreSQL database handler using the provided connection pool.
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ExecContext executes a SQL command using a transaction if present in the context, otherwise uses the connection pool.
func (p *pg) ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query, args...)
	}

	return p.dbc.Exec(ctx, query, args...)
}

// QueryContext executes a SQL query using a transaction if present in the context, otherwise uses the connection pool.
func (p *pg) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query, args...)
	}

	return p.dbc.Query(ctx, query, args...)
}

// QueryRowContext executes a SQL query that returns a single row using a transaction if present in the context, otherwise uses the connection pool.
func (p *pg) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.dbc.QueryRow(ctx, query, args...)
}

// BeginTx starts a new transaction with the specified options.
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping verifies the database connection is active.
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes the connection pool and releases all resources.
func (p *pg) Close() {
	p.dbc.Close()
}

// MakeContextTx returns a new context with the transaction embedded as a value for use in query operations.
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
