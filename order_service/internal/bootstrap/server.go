package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	orderHandler "order_service/internal/api/handlers/order"
	"order_service/internal/config"
	"order_service/internal/infra"
	"order_service/internal/kafka"
	orderRepo "order_service/internal/repo/order"
	orderService "order_service/internal/service/order"
	orderpb "order_service/proto/gen"

	"google.golang.org/grpc"
)

type App struct {
	server    *grpc.Server
	listener  net.Listener
	cnf       *config.Config
	prodClose func() error
	consClose func() error
	dbClose   func() error
}

func InitializeApp(ctx context.Context, cnf *config.Config) (*App, error) {
	db, err := infra.ConnectDb(cnf.DBCnf)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	kfInfra := infra.NewKafkaInfra([]string{"localhost:9092"})
	writer := kfInfra.Writer(kafka.OrderEventsTopic)
	reader := kfInfra.Reader(kafka.OrderResultsTopic, "order_service_group")

	producer, prodClose := kafka.NewProducer(writer)
	consumer, consClose := kafka.NewConsumer(reader)

	repo := orderRepo.NewRepo(db)
	service := orderService.NewService(producer, consumer, repo)
	handler := orderHandler.NewHandler(service)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", cnf.Addr)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("lis connection failed: %w", err)
	}
	orderpb.RegisterOrderServiceServer(grpcServer, handler)

	return &App{
		server:    grpcServer,
		listener:  lis,
		cnf:       cnf,
		prodClose: prodClose,
		consClose: consClose,
		dbClose:   db.Close,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	go func() {
		if err := a.server.Serve(a.listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()

	a.server.GracefulStop()
	a.prodClose()
	a.consClose()
	a.dbClose()
}
