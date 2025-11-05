package productservice

import (
	"context"
	"product_service/internal/domain"
	productrepo "product_service/internal/repo/productRepo"
	"product_service/internal/utils"
	productpb "product_service/proto/gen"

	"github.com/google/uuid"
)

type service struct {
	repo productrepo.ProductRepo
}
type Service interface {
	Create(ctx context.Context, payload *productpb.CreateProductRequest, email string) (*productpb.CreateProductResponse, error)
}

func NewService(repo productrepo.ProductRepo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, payload *productpb.CreateProductRequest, email string) (*productpb.CreateProductResponse, error) {
	uid := uuid.New().String()

	// Map gRPC payload to domain.Product
	productData := &domain.Product{
		ProductID:   uid,
		Name:        payload.Name,
		Description: payload.Description,
		Category:    payload.Category,
		Price:       payload.Price,
		CreatedBy:   email,
		ImageURLs:   payload.ImageUrls,
		Status:      payload.Status,
		IsFeatured:  payload.IsFeatured,
		Tags:        payload.Tags,
	}

	if err := s.repo.Create(ctx, productData); err != nil {
		return nil, err
	}
	return utils.ProductResponse(productData), nil
}
