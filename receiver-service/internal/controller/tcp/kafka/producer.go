package kafka

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/receiver/internal/config"
	kafkaClient "github.com/Entreeka/receiver/pkg/kafka"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

type ProducerError interface {
	WriteError(ctx context.Context, msgErr map[string]interface{}) error
}

type errorHandler struct {
	log           *logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewErrorProducerHandler(log *logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *errorHandler {
	return &errorHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (m *errorHandler) WriteError(ctx context.Context, msgErr map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgBytes, err := json.Marshal(&msgErr)
	if err != nil {
		m.log.Error("json.Marshal: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(m.cfg.Kafka.TopicError),
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	return m.kafkaProducer.PublishMessage(ctx, msg)
}
