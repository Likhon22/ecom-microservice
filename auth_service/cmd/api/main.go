package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/config"
	"log"
)

func main() {
	config := config.GetConfig()

	grpcServer, err := bootstrap.NewServer(config.Addr)
	if err != nil {
		log.Fatal(err)

	}
	grpcServer.StartServer()

}
