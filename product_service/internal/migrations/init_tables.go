package migrations

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func InitProductTable(client *dynamodb.Client) {
	tableName := "Products"

	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: &tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("Category"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("ProductID"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("Category"), KeyType: types.KeyTypeHash},   // Partition Key
			{AttributeName: aws.String("ProductID"), KeyType: types.KeyTypeRange}, // Sort Key
		},
		BillingMode: types.BillingModePayPerRequest, // Local mode ignores provisioned throughput
	})

	if err != nil {
		var exists *types.ResourceInUseException
		if errors.As(err, &exists) {
			log.Println("Table already exists:", tableName)
			return
		}
		log.Fatal("Failed to create table:", err)
	}

	log.Println("Table created successfully:", tableName)
}
