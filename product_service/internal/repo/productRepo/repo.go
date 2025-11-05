package productrepo

import (
	"context"
	"fmt"
	"product_service/internal/domain"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type productRepo struct {
	client    *dynamodb.Client
	tableName string
}
type FilterOptions struct {
	Category string
	Search   string
}

type ProductRepo interface {
	Create(ctx context.Context, product *domain.Product) error
	GetAll(ctx context.Context, filters *FilterOptions) ([]*domain.Product, int, error)
	GetById(ctx context.Context, productId, category string) (*domain.Product, error)
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

func (r *productRepo) GetAll(ctx context.Context, filters *FilterOptions) ([]*domain.Product, int, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	// Build filter expressions
	var filterExpressions []string
	expressionValues := make(map[string]types.AttributeValue)
	expressionNames := make(map[string]string)

	// Filter by category (exact match)
	if filters.Category != "" {
		filterExpressions = append(filterExpressions, "Category = :category")
		expressionValues[":category"] = &types.AttributeValueMemberS{Value: filters.Category}
	}

	// Search by name (contains - case sensitive in DynamoDB)
	if filters.Search != "" {
		filterExpressions = append(filterExpressions, "contains(#name, :search)")
		expressionNames["#name"] = "name" // 'name' might be reserved keyword
		expressionValues[":search"] = &types.AttributeValueMemberS{Value: filters.Search}
	}

	// Apply filter expressions
	if len(filterExpressions) > 0 {
		input.FilterExpression = aws.String(strings.Join(filterExpressions, " AND "))
		input.ExpressionAttributeValues = expressionValues
	}
	if len(expressionNames) > 0 {
		input.ExpressionAttributeNames = expressionNames
	}

	// Execute scan (gets ALL matching items)
	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan table: %w", err)
	}

	// Convert DynamoDB items to domain.Product
	allProducts := make([]*domain.Product, 0, len(result.Items))
	for _, item := range result.Items {
		var product domain.Product
		if err := attributevalue.UnmarshalMap(item, &product); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal product: %w", err)
		}
		allProducts = append(allProducts, &product)
	}

	totalCount := len(allProducts)
	return allProducts, totalCount, nil
}

func (r *productRepo) GetById(ctx context.Context, productId, category string) (*domain.Product, error) {

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{

		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Category":  &types.AttributeValueMemberS{Value: category},
			"ProductID": &types.AttributeValueMemberS{Value: productId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if result.Item == nil {
		return nil, fmt.Errorf("product not found")
	}
	var product domain.Product
	if err := attributevalue.UnmarshalMap(result.Item, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}
	return &product, nil

}
