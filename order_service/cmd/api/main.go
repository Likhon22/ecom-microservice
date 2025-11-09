package main

import (
	"context"
	"log"
	"order_service/internal/bootstrap"
	"order_service/internal/config"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	cnf := config.GetConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app, err := bootstrap.InitializeApp(ctx, cnf)

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	app.Run(ctx)
}
