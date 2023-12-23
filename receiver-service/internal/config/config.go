package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"os"
)

type (
	Config struct {
		Kafka    Kafka    `json:"kafka"`
		Postgres Postgres `json:"postgres"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	Kafka struct {
		Topic      string   `json:"topic"`
		TopicError string   `json:"topic_error"`
		Brokers    []string `json:"brokers"`
		GroupID    string   `json:"group_id"`
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
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
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
