package config

import (
	"github.com/joho/godotenv"
	"os"
)

type (
	Config struct {
		Postgres Postgres `json:"postgres"`
		Kafka    Kafka    `json:"kafka"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	Kafka struct {
	}
)

func New() (*Config, error) {
	err := godotenv.Load("configs/app.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
		},
		Kafka: Kafka{},
	}

	return config, nil
}
