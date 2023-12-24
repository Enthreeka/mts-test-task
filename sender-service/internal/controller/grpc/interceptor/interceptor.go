package interceptor

import (
	"context"
	"github.com/Entreeka/sender/pkg/logger"
	"google.golang.org/grpc"
)

func LoggerUnaryInterceptorServer(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log.Info("[SENDER-SERVICE] Received request: %v, %s", req, info.FullMethod)
		resp, err := handler(ctx, req)
		return resp, err
	}
}
