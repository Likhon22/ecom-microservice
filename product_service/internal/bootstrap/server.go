package bootstrap

import (
	"context"
	"log"
	"net"
	"product_service/internal/config"

	"google.golang.org/grpc"
)

type App struct {
	server   *grpc.Server
	listener net.Listener
	cfg      *config.Config
}

func InitializeApp(ctx context.Context, cfg *config.Config) (*App, error) {

	server := grpc.NewServer()

	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err

	}
	return &App{
		server:   server,
		listener: listener,
		cfg:      cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) {

	go func() {
		log.Printf("gRPC server listening at %v", a.listener.Addr())
		if err := a.server.Serve(a.listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("shutting down server gracefully...")
	a.server.GracefulStop()
}
