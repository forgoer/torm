package torm

import (
	"database/sql"
)

type Config struct {
	Driver string
	Prefix string
	Dsn    string
}

type Manager struct {
}

func Open(config Config) (*Connection, error) {
	db, err := sql.Open(config.Driver, config.Dsn)

	if err != nil {
		return nil, err
	}

	return &Connection{
		DB: db,
	}, nil
}
