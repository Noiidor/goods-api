package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_PORT = "8080"
	DEFAULT_HOST = "0.0.0.0"
)

type AppConfig struct {
	Server struct {
		Host string
		Port string
	}
	Data struct {
		Postgres struct {
			URL string
		}
		Redis struct {
			URL string
		}
		Clickhouse struct {
			URL string
		}
	}
	Communication struct {
		NATS struct {
			URL string
		}
	}
}

func LoadConfig() (AppConfig, error) {
	err := godotenv.Load("../../local.env")
	if err != nil {
		log.Printf("Could not load .env file")
	}

	config := AppConfig{}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	config.Server.Port = port

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = DEFAULT_HOST
	}
	config.Server.Host = host

	pgUrl := os.Getenv("POSTGRES_URL")
	if pgUrl == "" {
		return AppConfig{}, fmt.Errorf("env variable POSTGRES_URL is not set")
	}
	config.Data.Postgres.URL = pgUrl

	chUrl := os.Getenv("CLICKHOUSE_URL")
	if chUrl == "" {
		return AppConfig{}, fmt.Errorf("env variable CLICKHOUSE_URL is not set")
	}
	config.Data.Clickhouse.URL = chUrl

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		return AppConfig{}, fmt.Errorf("env variable REDIS_URL is not set")
	}
	config.Data.Redis.URL = redisUrl

	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		return AppConfig{}, fmt.Errorf("env variable NATS_URL is not set")
	}
	config.Communication.NATS.URL = natsUrl

	return config, nil
}
