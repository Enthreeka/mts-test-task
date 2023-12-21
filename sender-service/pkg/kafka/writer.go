package kafka

import "github.com/segmentio/kafka-go"

func NewWriter(brokers []string, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	return w
}
