package clickhouse

import (
	"goods-api/internal/config"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
)

func NewConnection(config *config.AppConfig) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("clickhouse", config.Data.Clickhouse.URL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
