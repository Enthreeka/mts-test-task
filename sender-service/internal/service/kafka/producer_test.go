package kafka

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/sender/internal/config"
	"github.com/Entreeka/sender/internal/entity"
	"github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMessageProducerService_CreateHandler(t *testing.T) {
	t.Parallel()
	log := logger.New()

	cfg := &config.Config{
		Kafka: config.Kafka{
			Topic:      "Test",
			TopicError: "TopicError",
			Brokers:    []string{"localhost:9092"},
		},
	}

	kafkaProducer := kafka.NewProducer(log, cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaProducer.Close()

	msgService := NewMessageProducerService(log, cfg, kafkaProducer)

	tests := []struct {
		name  string
		topic string
		msg   *entity.Message
	}{
		{
			name:  "ok1",
			topic: "Test",
			msg: &entity.Message{
				Msg:         "test msg",
				CreatedTime: time.Now(),
				MsgUUID:     uuid.New(),
			},
		},
		{
			name:  "ok2",
			topic: "Test",
			msg: &entity.Message{
				Msg:         "test msg14124124",
				CreatedTime: time.Now(),
				MsgUUID:     uuid.New(),
			},
		},
		{
			name:  "ok3",
			topic: "Test",
			msg: &entity.Message{
				Msg:         "test msg141246456456124",
				CreatedTime: time.Now(),
				MsgUUID:     uuid.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Log(tt.msg)
		msgService.CreateHandler(context.Background(), tt.msg)
	}

	reader := kafka.NewKafkaReader(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer reader.Close()

	timeout := time.After(10 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for {
				select {
				case <-timeout:
					t.Fatal("Timeout waiting for messages")
				default:
					msg, err := reader.FetchMessage(context.Background())
					if err != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}

					messageModel := &entity.Message{}
					err = json.Unmarshal(msg.Value, messageModel)
					if err != nil {
						t.Error(err)
					}

					if tt.topic == msg.Topic && tt.msg.Msg == messageModel.Msg {
						return
					}

					time.Sleep(100 * time.Millisecond)
				}
			}
		})
	}
}
