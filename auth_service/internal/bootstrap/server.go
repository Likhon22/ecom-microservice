package bootstrap

import (
	authhandler "auth_service/internal/api/handler/auth"
	userhandler "auth_service/internal/api/handler/user"
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/config"
	"auth_service/internal/infra/db"
	repo "auth_service/internal/repo/auth"
	"auth_service/internal/services/auth"
	userpb "auth_service/proto/gen"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	server     *grpc.Server
	cfg        *config.Config
	listener   net.Listener
	userClient usersvc.Client
}

func InitializeApp(ctx context.Context, cfg *config.Config) (*App, error) {
	// MongoDB
	client, err := db.ConnectMongo(cfg.DBCnf)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}
	log.Println("MongoDB connected successfully")
	rdb := db.ConnectRedis("localhost:6379", 0)
	log.Println("redis connected")

	// External gRPC client
	userClient, closeUserClient, err := usersvc.NewClient(ctx, cfg.User_Service_Addr)
	if err != nil {
		return nil, fmt.Errorf("dial user service: %w", err)
	}
	go func() {
		<-ctx.Done()
		closeUserClient()
		client.Disconnect(ctx)
		rdb.Close()
	}()

	// Listener
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()

	// repo

	authRepo := repo.NewAuthRepo(client, rdb)

	// Services and Handlers
	userService := auth.NewService(userClient, cfg.AuthCnf, authRepo)
	userHandler := userhandler.NewHandler(userService)
	authHandler := authhandler.NewHandler(userService)

	userpb.RegisterUserServiceServer(s, userHandler)
	userpb.RegisterAuthServiceServer(s, authHandler)

	return &App{
		server:     s,
		cfg:        cfg,
		listener:   lis,
		userClient: userClient,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	go func() {
		log.Printf("gRPC server listening at %s\n", a.cfg.Addr)
		if err := a.server.Serve(a.listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server gracefully...")
	a.server.GracefulStop()
}
