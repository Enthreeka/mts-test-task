package grpc

import (
	"context"
	"fmt"
	"github.com/Entreeka/sender/internal/apperror"
	"github.com/Entreeka/sender/internal/entity"
	"github.com/Entreeka/sender/internal/service/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type messageHandler struct {
	log      *logger.Logger
	producer kafka.MessageProducerService

	pb.UnimplementedMessageServiceServer
}

func NewMessageHandler(log *logger.Logger, producer kafka.MessageProducerService) pb.MessageServiceServer {
	return &messageHandler{
		log:      log,
		producer: producer,
	}
}

func (m *messageHandler) CreateMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	if req.Message == "" {
		m.log.Error("req.Message == \"\"")
		return nil, status.Errorf(apperror.ParseGRPCErrStatusCode(apperror.ErrEmptyMessage), "CreateMessage: %v", apperror.ErrEmptyMessage)
	}

	msg := m.messageProtoToModel(req)

	msgUUID := uuid.New()
	msg.MsgUUID = msgUUID

	err := m.producer.CreateHandler(ctx, msg)
	if err != nil {
		m.log.Error("producer.CreateHandler: %v", err)
		return nil, status.Errorf(apperror.ParseGRPCErrStatusCode(err), "CreateMessage: %v", err)
	}

	return m.messageModelToProto(msg), nil
}

func (m *messageHandler) messageModelToProto(message *entity.Message) *pb.MessageResponse {
	fmt.Println(message.MsgUUID.String())
	return &pb.MessageResponse{
		Message:     message.Msg,
		CreatedTime: timestamppb.New(message.CreatedTime),
		ID:          message.MsgUUID.String(),
	}
}

func (m *messageHandler) messageProtoToModel(req *pb.MessageRequest) *entity.Message {
	return &entity.Message{
		Msg:         req.Message,
		CreatedTime: time.Now(),
	}
}
