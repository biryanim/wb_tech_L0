package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig creates and returns a new PostgreSQL configuration by reading the DSN from the PG_DSN environment variable.
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{dsn: dsn}, nil
}

// DSN returns the database source name (connection string) for PostgreSQL.
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
