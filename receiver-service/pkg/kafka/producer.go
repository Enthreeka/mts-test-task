package kafka

import (
	"context"
	"github.com/Entreeka/receiver/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msg ...kafka.Message) error
	Close() error
}

type producer struct {
	log     *logger.Logger
	brokers []string
	w       *kafka.Writer
}

// TODO make options
func NewProducer(log *logger.Logger, brokers []string, topic string) *producer {
	return &producer{
		log:     log,
		brokers: brokers,
		w:       NewWriter(brokers, topic),
	}
}

func (p *producer) PublishMessage(ctx context.Context, msg ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msg...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
