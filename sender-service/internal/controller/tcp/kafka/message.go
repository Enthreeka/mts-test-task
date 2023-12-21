package kafka

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/sender/internal/config"
	"github.com/Entreeka/sender/internal/entity"
	kafkaClient "github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

type ProducerMessage interface {
	CreateHandler(ctx context.Context, message *entity.Message) error
}

type messageHandler struct {
	log           *logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewMessageProducerHandler(log *logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *messageHandler {
	return &messageHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (m *messageHandler) CreateHandler(ctx context.Context, message *entity.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgBytes, err := json.Marshal(&message)
	if err != nil {
		m.log.Error("json.Marshal: %v", err)
	}

	msg := kafka.Message{
		Topic: m.cfg.Kafka.Topic,
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	return m.kafkaProducer.PublishMessage(ctx, msg)
}
