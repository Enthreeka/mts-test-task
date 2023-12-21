package grpc

import (
	"context"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
)

type messageHandler struct {
	log *logger.Logger

	pb.UnimplementedMessageServiceServer
}

func NewMessageHandler(log *logger.Logger) *messageHandler {
	return &messageHandler{
		log: log,
	}
}

func (m *messageHandler) CreateMessageHandler(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	return nil, nil
}
