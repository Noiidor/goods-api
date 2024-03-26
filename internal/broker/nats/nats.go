package nats

import (
	"goods-api/internal/config"

	"github.com/nats-io/nats.go"
)

func NewConnection(config *config.AppConfig) (*nats.Conn, error) {
	conn, err := nats.Connect(config.Communication.NATS.URL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
