package kafka

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/internal/service"
	kafkaClient "github.com/Entreeka/receiver/pkg/kafka"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/segmentio/kafka-go"
	"sync"
)

type Consumer interface {
	Consumer(ctx context.Context, wg *sync.WaitGroup)
}

type messageHandler struct {
	msgService service.Message
	log        *logger.Logger
	cfg        *config.Config

	kafkaConsumer *kafka.Reader
	kafkaProducer ProducerError
}

func NewMessageConsumerHandler(msgService service.Message, kafkaProducer ProducerError, log *logger.Logger, cfg *config.Config) *messageHandler {
	return &messageHandler{
		msgService:    msgService,
		kafkaProducer: kafkaProducer,
		log:           log,
		cfg:           cfg,
		kafkaConsumer: kafkaClient.NewKafkaReader(cfg.Kafka.Brokers, cfg.Kafka.Topic),
	}
}

func (m *messageHandler) Consumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {

		case <-ctx.Done():
			return

		default:
			msg, err := m.kafkaConsumer.FetchMessage(ctx)
			if err != nil {
				m.log.Error("FetchMessage: %v", err)
				continue
			}

			m.log.Info("received: %s, %s", msg.Topic, string(msg.Value))
			m.createMessage(ctx, &msg)
		}
	}
}

func (m *messageHandler) createMessage(ctx context.Context, msg *kafka.Message) {
	messageModel := &entity.Message{}
	err := json.Unmarshal(msg.Value, messageModel)
	if err != nil {
		m.log.Error("Unmarshal: %v", err)
		if err := m.kafkaConsumer.CommitMessages(ctx, *msg); err != nil {
			m.log.Error("CommitMessages: %v", err)
		}

		m.kafkaProducer.WriteError(ctx, map[string]string{"error": err.Error()})
		return
	}

	err = m.msgService.Create(ctx, messageModel)
	if err != nil {
		m.log.Error("msgService.Create: %v", err)
		if err := m.kafkaConsumer.CommitMessages(ctx, *msg); err != nil {
			m.log.Error("CommitMessages: %v", err)
		}

		m.kafkaProducer.WriteError(ctx, map[string]string{"error": err.Error()})
		return
	}
}
