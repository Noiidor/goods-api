package postgres

import (
	"goods-api/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewConnection(config *config.AppConfig) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("pgx", config.Data.Postgres.URL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
