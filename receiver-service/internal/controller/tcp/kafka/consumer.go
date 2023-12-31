package kafka

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/internal/entity"
	"github.com/Entreeka/receiver/internal/service"
	kafkaService "github.com/Entreeka/receiver/internal/service/kafka"
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
	kafkaProducer kafkaService.ErrorProducerService
}

func NewMessageConsumerHandler(msgService service.Message, kafkaProducer kafkaService.ErrorProducerService, log *logger.Logger, cfg *config.Config) *messageHandler {
	kafkaConsumer := kafkaClient.NewKafkaReader(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	err := kafkaConsumer.SetOffset(0)
	if err != nil {
		log.Error("SetOffset: %v", err)
	}

	return &messageHandler{
		msgService:    msgService,
		kafkaProducer: kafkaProducer,
		log:           log,
		cfg:           cfg,
		kafkaConsumer: kafkaConsumer,
	}
}

// Как вариант можно реализовать consumer_group, в которой будет реализован пул воркеров для параллельной обработки данных.

func (m *messageHandler) Consumer(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if err := m.kafkaConsumer.Close(); err != nil {
			m.log.Error("kafkaConsumer.Close: %v", err)
		}
	}()

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

			m.log.Info("[RECEIVER-SERVICE] get message at topic/partition/offset"+
				" %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

			m.createMessage(ctx, &msg)
		}
	}
}

func (m *messageHandler) createMessage(ctx context.Context, msg *kafka.Message) {
	messageModel := &entity.Message{}
	err := json.Unmarshal(msg.Value, messageModel)
	if err != nil {
		m.log.Error("Unmarshal: %v", err)

		if err = m.kafkaProducer.WriteError(ctx, map[string]interface{}{
			"error": err.Error(),
		}); err != nil {
			m.log.Error("WriteError: %v", err)
		}

		return
	}

	err = m.msgService.Create(ctx, messageModel)
	if err != nil {
		m.log.Error("msgService.Create: %v", err)

		if err = m.kafkaProducer.WriteError(ctx, map[string]interface{}{
			"error":    err.Error(),
			"msg_uuid": messageModel.MsgUUID,
		}); err != nil {
			m.log.Error("WriteError: %v", err)
		}

		return
	}
}
