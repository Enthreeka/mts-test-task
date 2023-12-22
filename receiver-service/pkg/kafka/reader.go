package kafka

import "github.com/segmentio/kafka-go"

func NewKafkaReader(brokers []string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MaxBytes: 10e6,
	})
}
