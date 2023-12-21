package grpc

import (
	"context"
	"github.com/Entreeka/sender/internal/entity"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type messageHandler struct {
	log *logger.Logger

	pb.UnimplementedMessageServiceServer
}

func NewMessageHandler(log *logger.Logger) pb.MessageServiceServer {
	return &messageHandler{
		log: log,
	}
}

func (m *messageHandler) CreateMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	message := &entity.Message{
		Msg: req.Message,
	}

	return m.messageModelToProto(message), nil
}

func (m *messageHandler) messageModelToProto(message *entity.Message) *pb.MessageResponse {
	return &pb.MessageResponse{
		Message:     message.Msg,
		CreatedTime: timestamppb.New(message.CreatedTime),
	}
}
