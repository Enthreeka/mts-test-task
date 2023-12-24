package app

import (
	"context"
	"fmt"
	"github.com/Entreeka/sender/internal/config"
	msgGrpc "github.com/Entreeka/sender/internal/controller/grpc"
	"github.com/Entreeka/sender/internal/controller/grpc/interceptor"
	kafkaHandler "github.com/Entreeka/sender/internal/controller/tcp/kafka"
	kafkaService "github.com/Entreeka/sender/internal/service/kafka"
	"github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
	pb "github.com/Entreeka/sender/proto/v1"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, log *logger.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	kafkaProducer := kafka.NewProducer(log, cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaProducer.Close()

	conn, err := kafka.New(ctx)
	if err != nil {
		log.Fatal("failed to dial leader: %v", err)
	}

	brokers, err := conn.Brokers()
	if err != nil {
		return err
	}
	log.Info("kafka connected to brokers: %+v", brokers)

	errorHandler := kafkaHandler.NewMessageConsumerHandler(log, cfg)

	go errorHandler.ReadError(ctx)

	msgService := kafkaService.NewMessageProducerService(log, cfg, kafkaProducer)
	
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerUnaryInterceptorServer(log),
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	s := grpc.NewServer(opts...)
	urlGRPCServer := msgGrpc.NewMessageHandler(log, msgService)

	pb.RegisterMessageServiceServer(s, urlGRPCServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return err
	}

	go func() {
		log.Info("Starting gRPC listener on port :%s", cfg.GRPCServer.Port)
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	s.GracefulStop()
	return nil
}
