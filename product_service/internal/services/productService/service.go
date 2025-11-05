package productservice

import (
	"context"
	"errors"
	client "product_service/internal/client/product"
	"product_service/internal/domain"
	productrepo "product_service/internal/repo/productRepo"
	"product_service/internal/utils"
	productpb "product_service/proto/gen"

	"github.com/google/uuid"
)

type service struct {
	client client.Client
	repo   productrepo.ProductRepo
}
type Service interface {
	Create(ctx context.Context, payload *productpb.CreateProductRequest, email string) (*productpb.CreateProductResponse, error)
	GetAll(ctx context.Context) (*productpb.GetProductsResponse, error)
}

func NewService(client client.Client, repo productrepo.ProductRepo) Service {
	return &service{
		repo:   repo,
		client: client,
	}
}

func (s *service) Create(ctx context.Context, payload *productpb.CreateProductRequest, email string) (*productpb.CreateProductResponse, error) {

	customer, err := s.client.GetCustomerByEmail(ctx, &productpb.GetCustomerByEmailRequest{Email: email})
	uid := uuid.New().String()
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("email is invalid")
	}

	// Map gRPC payload to domain.Product
	productData := &domain.Product{
		ProductID:   uid,
		Name:        payload.Name,
		Description: payload.Description,
		Category:    payload.Category,
		Price:       payload.Price,
		CreatedBy:   customer.Email,
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

func (s *service) GetAll(ctx context.Context) (*productpb.GetProductsResponse, error) {

	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	pbProducts := make([]*productpb.Product, 0, len(products))
	for _, p := range products {
		pbProducts = append(pbProducts, &productpb.Product{
			ProductId:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			Price:       p.Price,
			ImageUrls:   p.ImageURLs,
			Status:      p.Status,
			IsFeatured:  p.IsFeatured,
			Tags:        p.Tags,
			CreatedBy:   p.CreatedBy,
		})
	}
	return &productpb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
