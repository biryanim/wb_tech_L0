package pg

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"

	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

// New creates and returns a new PostgreSQL client connection using the provided DSN.
func New(ctx context.Context, config *pgxpool.Config) (db.Client, error) {
	dbc, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	if err = otelpgx.RecordStats(dbc); err != nil {
		return nil, fmt.Errorf("failed to record database stats: %v", err)
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
