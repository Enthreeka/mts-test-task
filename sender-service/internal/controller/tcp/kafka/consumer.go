package kafka

import (
	"context"
	"github.com/Entreeka/sender/internal/config"
	kafkaClient "github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	ReadError(ctx context.Context)
}

type consumerErrorHandler struct {
	log *logger.Logger
	cfg *config.Config

	kafkaConsumer *kafka.Reader
}

func NewMessageConsumerHandler(log *logger.Logger, cfg *config.Config) *consumerErrorHandler {
	kafkaConsumer := kafkaClient.NewKafkaReader(cfg.Kafka.Brokers, cfg.Kafka.TopicError)
	err := kafkaConsumer.SetOffset(-1)
	if err != nil {
		log.Error("SetOffset: %v", err)
	}

	return &consumerErrorHandler{
		log:           log,
		cfg:           cfg,
		kafkaConsumer: kafkaConsumer,
	}
}

func (c *consumerErrorHandler) ReadError(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.kafkaConsumer.FetchMessage(ctx)
			if err != nil {
				c.log.Error("FetchMessage: %v", err)
				continue
			}

			c.log.Error("Topic:%s, msg:%s", msg.Topic, string(msg.Value))
		}
	}
}
