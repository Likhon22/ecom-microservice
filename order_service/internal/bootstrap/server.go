package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"order_service/internal/config"
	"order_service/internal/infra"

	"google.golang.org/grpc"
)

type App struct {
	server   *grpc.Server
	listener net.Listener
	cnf      *config.Config
}

func InitializeApp(ctx context.Context, cnf *config.Config) (*App, error) {

	db, err := infra.ConnectDb(cnf.DBCnf)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)

	}

	log.Println("Db connection successful")
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", cnf.Addr)
	if err != nil {
		return nil, fmt.Errorf("lis connection failed: %w", err)

	}
	go func() {
		<-ctx.Done()
		db.Close()
	}()
	return &App{
		server:   grpcServer,
		listener: lis,
		cnf:      cnf,
	}, nil

}

func (a *App) Run(ctx context.Context) {
	go func() {
		log.Println("grpc server running on:", a.cnf.Addr)
		if err := a.server.Serve(a.listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("shutting down server gracefully...")
	a.server.GracefulStop()
}
