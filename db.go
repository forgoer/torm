package torm

import (
	"database/sql"
)

func Open(config Config) (*Connection, error) {
	db, err := sql.Open(config.Driver, config.Dsn)

	if err != nil {
		return nil, err
	}

	return &Connection{
		DB: db,
	}, nil
}
