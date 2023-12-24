package grpc

import (
	"fmt"
	"github.com/Entreeka/sender/internal/config"
	"github.com/Entreeka/sender/internal/controller/grpc/interceptor"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GRPCServer pb.MessageServiceServer
}

func NewServer(cfg *config.Config, log *logger.Logger, server Server) (func() error, *grpc.Server) {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerUnaryInterceptorServer(log),
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterMessageServiceServer(s, server.GRPCServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	go func() {
		log.Info("Starting gRPC listener on port :%s", cfg.GRPCServer.Port)
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()

	return lis.Close, s
}
