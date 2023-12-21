package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Kafka Kafka `json:"kafka"`
	}

	Kafka struct {
		Topic   string   `json:"topic"`
		Brokers []string `json:"brokers"`
	}
)

func New() (*Config, error) {
	cfgKafka, err := KafkaConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
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
