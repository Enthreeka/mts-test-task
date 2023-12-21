package app

import (
	"context"
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/segmentio/kafka-go"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		log.Info("received: %s", string(msg.Value))
	}

	return nil
}
