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
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.InitializeApp(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	app.Run(ctx)
}
