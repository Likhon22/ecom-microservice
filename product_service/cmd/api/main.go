package main

import (
	"context"
	"log"
	"os/signal"
	"product_service/internal/bootstrap"
	"product_service/internal/config"
	"syscall"
)

func main() {

	cfg := config.GetConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app, err := bootstrap.InitializeApp(ctx, cfg)
	if err != nil {
		log.Fatal("failed to start the server", err)

	}
	app.Run(ctx)

}
