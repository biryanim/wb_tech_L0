package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler func(ctx context.Context) error

type SQLExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}

type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReadCommited(cxt context.Context, f Handler) error
}
