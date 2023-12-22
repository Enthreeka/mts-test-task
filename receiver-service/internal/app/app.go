package app

import (
	"context"
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/internal/repo"
	"github.com/Entreeka/receiver/internal/service"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/Entreeka/receiver/pkg/postgres"
	"github.com/segmentio/kafka-go"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}

	defer psql.Close()

	messageRepo := repo.NewMessageRepo(psql)
	messageService := service.NewMessageService(messageRepo)

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
