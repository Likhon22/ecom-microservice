package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"product_service/internal/api/handlers/product"
	client "product_service/internal/client/product"
	"product_service/internal/config"
	"product_service/internal/infra/db"
	"product_service/internal/interceptors"
	"product_service/internal/migrations"
	productrepo "product_service/internal/repo/productRepo"
	productservice "product_service/internal/services/productService"
	productpb "product_service/proto/gen"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/grpc"
)

type App struct {
	server   *grpc.Server
	listener net.Listener
	cfg      *config.Config
}

func InitializeApp(ctx context.Context, cfg *config.Config) (*App, error) {

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.ErrorInterCeptor()))
	userclient, closeUserClient, err := client.NewClient(ctx, cfg.UserServiceAddress)

	if err != nil {
		return nil, fmt.Errorf("dial user service: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err

	}
	dynamoDBConfig := db.GetDBConfig()
	client := dynamodb.NewFromConfig(dynamoDBConfig)
	log.Println("dynamo db connected")
	migrations.InitProductTable(client)

	go func() {
		<-ctx.Done()
		closeUserClient()

	}()

	productRepo := productrepo.NewRepo(client, "Products")
	productService := productservice.NewService(userclient, productRepo)
	productHandler := product.NewProductHandler(productService)
	productpb.RegisterProductServiceServer(server, productHandler)
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
