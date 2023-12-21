package config

import (
	"github.com/joho/godotenv"
	"os"
)

type (
	Config struct {
		GRPCServer GRPCServer `json:"grpc_server"`
		Kafka      Kafka      `json:"kafka"`
	}

	GRPCServer struct {
		Port string `json:"port"`
	}

	Kafka struct {
		Topic   string   `json:"topic"`
		Brokers []string `json:"brokers"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load("configs/app.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		GRPCServer: GRPCServer{
			Port: os.Getenv("GRPC_SERVER_PORT"),
		},
		Kafka: Kafka{
			Topic: os.Getenv("KAFKA_TOPIC"),
		},
	}

	return config, nil
}
