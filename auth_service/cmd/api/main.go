package main

import (
	authhandler "auth_service/internal/api/handler/auth"
	userhanlder "auth_service/internal/api/handler/user"
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
	userHandler := userhanlder.NewHandler(userService)
	// auth handler
	authHandler := authhandler.NewHandler(userService)
	userpb.RegisterUserServiceServer(srv.GRPCServer(), userHandler)
	userpb.RegisterAuthServiceServer(srv.GRPCServer(), authHandler)
	srv.StartServer(ctx)

}
