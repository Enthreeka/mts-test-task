package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func New(ctx context.Context) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", "localhost:9092")
}
