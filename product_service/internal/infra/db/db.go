package db

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var (
	once sync.Once
	cfg  aws.Config
)

func loadDBConfig(dbUrl string) {
	loadedCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: dbUrl}, nil
			}),
		),
	)
	if err != nil {
		log.Fatal(err)

	}
	cfg = loadedCfg
}

func GetDBConfig(dbUrl string) aws.Config {
	once.Do(func() {
		loadDBConfig(dbUrl)
	})
	return cfg
}
