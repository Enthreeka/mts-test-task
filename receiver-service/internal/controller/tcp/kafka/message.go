package kafka

import (
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/pkg/logger"
)

type Consumer interface {
}

type messageHandler struct {
	log *logger.Logger
	cfg *config.Config
}

func NewMessageConsumerHandler(log *logger.Logger, cfg *config.Config) *messageHandler {
	return &messageHandler{
		log: log,
		cfg: cfg,
	}
}
