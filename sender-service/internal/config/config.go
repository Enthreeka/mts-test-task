package config

import (
	"encoding/json"
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

	cfgKafka, err := KafkaConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		GRPCServer: GRPCServer{
			Port: os.Getenv("GRPC_SERVER_PORT"),
		},
		Kafka: *cfgKafka,
	}

	return config, nil
}

func KafkaConfig() (*Kafka, error) {
	file, err := os.Open("configs/config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Kafka{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
