package app

import (
	"fmt"
	"github.com/Entreeka/sender/internal/config"
	msgGrpc "github.com/Entreeka/sender/internal/controller/grpc"
	"github.com/Entreeka/sender/internal/controller/grpc/interceptor"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	"google.golang.org/grpc"
	"net"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerUnaryInterceptorServer(log),
		),
	}

	s := grpc.NewServer(opts...)
	urlGRPCServer := msgGrpc.NewMessageHandler(log)

	pb.RegisterMessageServiceServer(s, urlGRPCServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return err
	}

	log.Info("Starting gRPC listener on port :%s", cfg.GRPCServer.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
		return err
	}

	return nil
}
