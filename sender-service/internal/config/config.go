package config

import (
	"github.com/joho/godotenv"
	"os"
)

type (
	Config struct {
		GRPCServer GRPCServer `json:"grpc_server"`
	}

	GRPCServer struct {
		Port string `json:"port"`
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
	}

	return config, nil
}
