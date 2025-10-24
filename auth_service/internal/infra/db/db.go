package db

import (
	"auth_service/internal/config"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoInstance *mongo.Client

func ConnectMongo(cnf *config.DBConfig) (*mongo.Client, error) {
	if mongoInstance != nil {
		return mongoInstance, nil

	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(cnf.DBUrl)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err

	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	mongoInstance = client
	return client, nil
}
