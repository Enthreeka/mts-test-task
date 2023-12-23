package kafka

import (
	"context"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/kafka"
	kafkaclient "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMessageProducerService_CreateHandler(t *testing.T) {
	t.Parallel()

	p := kafka.Preset(
		kafka.WithTopics("main", "error"))

	container, err := gnomock.Start(p,
		gnomock.WithDebugMode(),
		gnomock.WithLogWriter(os.Stdout),
		gnomock.WithContainerName("kafka"))
	require.NoError(t, err)

	defer func() {
		require.NoError(t, gnomock.Stop(container))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	alertsReader := kafkaclient.NewReader(kafkaclient.ReaderConfig{
		Brokers: []string{container.Address(kafka.BrokerPort)},
		Topic:   "main",
	})

	m, err := alertsReader.ReadMessage(ctx)
	require.NoError(t, err)
	require.NoError(t, alertsReader.Close())

	require.Equal(t, "CPU", string(m.Key))
	require.Equal(t, "92", string(m.Value))
}
