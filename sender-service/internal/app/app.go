package app

import (
	"context"
	"fmt"
	"github.com/Entreeka/sender/internal/config"
	msgGrpc "github.com/Entreeka/sender/internal/controller/grpc"
	"github.com/Entreeka/sender/internal/controller/grpc/interceptor"
	kafkaHandler "github.com/Entreeka/sender/internal/controller/tcp/kafka"
	"github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	"google.golang.org/grpc"
	"net"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	kafkaProducer := kafka.NewProducer(log, cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaProducer.Close()

	conn, err := kafka.New(context.Background())
	if err != nil {
		log.Fatal("failed to dial leader: %v", err)
	}

	brokers, err := conn.Brokers()
	if err != nil {
		return err
	}
	log.Info("kafka connected to brokers: %+v", brokers)

	errorHandler := kafkaHandler.NewMessageConsumerHandler(log, cfg)

	go errorHandler.ReadError(context.Background())

	msgHandler := kafkaHandler.NewMessageProducerHandler(log, cfg, kafkaProducer)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerUnaryInterceptorServer(log),
		),
	}

	s := grpc.NewServer(opts...)
	urlGRPCServer := msgGrpc.NewMessageHandler(log, msgHandler)

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
