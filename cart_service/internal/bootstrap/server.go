package bootstrap

import (
	"cart_service/internal/api/handlers"
	client "cart_service/internal/clients/product"
	"cart_service/internal/config"
	"cart_service/internal/infra"
	"cart_service/internal/interceptors"
	cartRepo "cart_service/internal/repo/cart"
	cartService "cart_service/internal/services/cart"
	cartpb "cart_service/proto/gen"
	"context"
	"fmt"
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.ErrorInterCeptor()))

	lis, err := net.Listen("tcp", cnf.Addr)

	if err != nil {
		return nil, err

	}
	rdb := infra.ConnectRedis("localhost:6379", 0)
	log.Println("redis is connected")

	productClient, closeProductClient, err := client.NewClient(ctx, cnf.User_Service_Addr)
	if err != nil {
		return nil, fmt.Errorf("dial user service: %w", err)
	}
	repo := cartRepo.NewRepo(rdb)
	service := cartService.NewService(repo, productClient)
	handler := handlers.NewHandler(service)
	cartpb.RegisterCartServiceServer(grpcServer, handler)

	go func() {
		<-ctx.Done()
		lis.Close()
		closeProductClient()
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
