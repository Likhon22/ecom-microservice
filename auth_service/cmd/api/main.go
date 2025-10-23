package main

import (
	authhanlder "auth_service/internal/api/handler/auth"
	"auth_service/internal/bootstrap"
	"auth_service/internal/clients/usersvc"
	"auth_service/internal/config"
	"auth_service/internal/services/auth"
	userpb "auth_service/proto/gen"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	config := config.GetConfig()

	srv, err := bootstrap.NewServer(config.Addr)
	if err != nil {
		log.Fatal(err)

	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	userClient, closeUserClient, err := usersvc.NewClient(ctx, config.User_Service_Addr)
	if err != nil {
		log.Fatalf("dial user services: %v", err)

	}

	defer func() {
		if err := closeUserClient(); err != nil {
			log.Printf("closing user client: %v", err)
		}
	}()
	userService := auth.NewService(userClient)
	userHandler := authhanlder.NewHandler(userService)
	userpb.RegisterUserServiceServer(srv.GRPCServer(), userHandler)
	srv.StartServer(ctx)

}
