package pg

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

// New creates and returns a new PostgreSQL client connection using the provided DSN.
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	return &pgClient{
		masterDBC: NewDB(dbc),
	}, nil
}

// DB returns the underlying database connection interface.
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close closes the database connection and releases resources.
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
