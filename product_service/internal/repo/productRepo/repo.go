package productrepo

import (
	"context"
	"fmt"
	"product_service/internal/domain"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type productRepo struct {
	client    *dynamodb.Client
	tableName string
}

type ProductRepo interface {
	Create(ctx context.Context, product *domain.Product) error
}

func NewRepo(client *dynamodb.Client, tableName string) ProductRepo {
	return &productRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	av, err := attributevalue.MarshalMap(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}
	return nil
}
