package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // import pq drivers.

	"github.com/cedric-parisi/fizzbuzz-api/config"
)

const (
	driverName     = "postgres"
	dataSourceName = "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
)

// NewPostgres creates a new postgres connection from the config.
func NewPostgres(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, fmt.Sprintf(dataSourceName,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbName,
		cfg.DbPassword,
	))
	if err != nil {
		return nil, err
	}
	return db, nil
}
