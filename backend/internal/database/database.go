package database

import (
	"github.com/aaltgod/bezdna/internal/config"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

type DBAdapter struct {
	Pool *pgx.ConnPool
}

func New(cfg config.DBConfig) (*DBAdapter, error) {
	pool, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     cfg.Host,
				Port:     cfg.Port,
				Database: cfg.Database,
				User:     cfg.Username,
				Password: cfg.Password,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "pgx.NewConnPool")
	}

	return &DBAdapter{
		Pool: pool,
	}, nil
}
