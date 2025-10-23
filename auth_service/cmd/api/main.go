package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/config"
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

	
	srv.StartServer(ctx)

}
