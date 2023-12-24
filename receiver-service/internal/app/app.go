package app

import (
	"context"
	"github.com/Entreeka/receiver/internal/config"
	kafkaClient "github.com/Entreeka/receiver/internal/controller/tcp/kafka"
	"github.com/Entreeka/receiver/internal/repo"
	"github.com/Entreeka/receiver/internal/service"
	kafkaService "github.com/Entreeka/receiver/internal/service/kafka"
	"github.com/Entreeka/receiver/pkg/kafka"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/Entreeka/receiver/pkg/postgres"
	"sync"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}

	defer psql.Close()

	producer := kafka.NewProducer(log, cfg.Kafka.Brokers, cfg.Kafka.TopicError)

	defer producer.Close()

	conn, err := kafka.New(context.Background())
	if err != nil {
		log.Fatal("failed to dial leader: %v", err)
	}

	brokers, err := conn.Brokers()
	if err != nil {
		return err
	}
	log.Info("kafka connected to brokers: %+v", brokers)

	messageRepo := repo.NewMessageRepo(psql)
	messageService := service.NewMessageService(messageRepo)
	errProducerService := kafkaService.NewErrorProducerService(log, cfg, producer)
	handler := kafkaClient.NewMessageConsumerHandler(messageService, errProducerService, log, cfg)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go handler.Consumer(context.Background(), wg)
	wg.Wait()

	return nil
}
