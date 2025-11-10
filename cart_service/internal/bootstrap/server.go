package bootstrap

import (
	"cart_service/internal/config"
	"cart_service/internal/infra"
	cartRepo "cart_service/internal/repo/cart"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Application struct {
	server   *grpc.Server
	listener net.Listener
	cnf      *config.Config
}

func InitializeApp(ctx context.Context, cnf *config.Config) (*Application, error) {

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", cnf.Addr)

	if err != nil {
		return nil, err

	}
	rdb := infra.ConnectRedis("localhost:6379", 0)
	repo := cartRepo.NewRepo(rdb)
	log.Println("redis is connected")

	go func() {
		<-ctx.Done()
		lis.Close()
	}()
	return &Application{
		server:   grpcServer,
		listener: lis,
		cnf:      cnf,
	}, nil

}

func (a *Application) Run(ctx context.Context) {

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
