package app

import (
	"context"
	"github.com/Entreeka/sender/internal/config"
	"github.com/Entreeka/sender/internal/controller/grpc"
	kafkaHandler "github.com/Entreeka/sender/internal/controller/tcp/kafka"
	kafkaService "github.com/Entreeka/sender/internal/service/kafka"
	"github.com/Entreeka/sender/pkg/kafka"
	"github.com/Entreeka/sender/pkg/logger"
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

	msgService := kafkaService.NewMessageProducerService(log, cfg, kafkaProducer)
	msgGRPCServer := grpc.NewMessageHandler(log, msgService)
	errorHandler := kafkaHandler.NewMessageConsumerHandler(log, cfg)

	go errorHandler.ReadError(ctx)

	closeGrpcServer, grpcServer := grpc.NewServer(cfg, log, grpc.Server{GRPCServer: msgGRPCServer})
	defer closeGrpcServer()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
