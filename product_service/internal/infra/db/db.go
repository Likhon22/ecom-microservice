package db

import (
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

func loadDBConfig() {
	loadedCfg, err := config.LoadDefaultConfig(nil,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:9000"}, nil
			}),
		),
	)
	if err != nil {
		log.Fatal(err)

	}
	cfg = loadedCfg
}

func GetDBConfig() aws.Config {
	once.Do(loadDBConfig)
	return cfg
}
