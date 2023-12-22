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
	Consumer(ctx context.Context) error
}

type messageHandler struct {
	msgService service.Message
	log        *logger.Logger
	cfg        *config.Config

	kafkaConsumer *kafka.Reader
}

func NewMessageConsumerHandler(msgService service.Message, log *logger.Logger, cfg *config.Config) *messageHandler {
	return &messageHandler{
		msgService:    msgService,
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
				//m.log.Error("workerID: %v, err: %v", workerID, err)
				m.log.Error("FetchMessage: %v", err)
				continue
			}

			m.log.Info("received: %v", msg)

			messageModel := &entity.Message{}
			err = json.Unmarshal(msg.Value, messageModel)
			if err != nil {
				if err := m.kafkaConsumer.CommitMessages(ctx, msg); err != nil {
					m.log.Error("CommitMessages: %v", err)
				}
			}

			err = m.msgService.Create(ctx, messageModel)
			if err != nil {
				if err := m.kafkaConsumer.CommitMessages(ctx, msg); err != nil {
					m.log.Error("msgService.Create: %v", err)
				}
			}
		}
	}
}
