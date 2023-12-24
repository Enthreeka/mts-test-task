package grpc

import (
	"context"
	"github.com/Entreeka/sender/internal/controller/grpc/interceptor"
	"github.com/Entreeka/sender/internal/service/kafka/mock"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"testing"
	"time"
)

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	log := logger.New()

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerUnaryInterceptorServer(log),
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	s := grpc.NewServer(opts...)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgServiceMock := mock.NewMockMessageProducerService(ctrl)
	msgServiceMock.EXPECT().CreateHandler(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	urlGRPCServer := NewMessageHandler(log, msgServiceMock)

	pb.RegisterMessageServiceServer(s, urlGRPCServer)

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatal("%v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		msg     string
		res     *pb.MessageResponse
		errCode codes.Code
	}{
		{
			name: "ok",
			msg:  "test message",
			res: &pb.MessageResponse{
				Message:     "test message",
				CreatedTime: timestamppb.New(time.Now()),
				MsgUUID:     uuid.New().String(),
			},
			errCode: codes.OK,
		},
		{
			name:    "empty result",
			msg:     "",
			res:     &pb.MessageResponse{},
			errCode: codes.InvalidArgument,
		},
	}

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(t)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.MessageRequest{
				Message: tt.msg,
			}

			response, err := client.CreateMessage(ctx, request)
			if response != nil {
				assert.Equal(t, tt.res.Message, response.Message)
			}

			if err, ok := status.FromError(err); ok {
				if err.Code() != tt.errCode {
					t.Error("error code: expected", tt.errCode, "received", err.Code())
				}
			}
		})
	}
}
