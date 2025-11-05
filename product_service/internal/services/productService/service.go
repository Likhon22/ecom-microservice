package productservice

import (
	"context"
	"errors"
	"log"
	client "product_service/internal/client/product"
	"product_service/internal/domain"
	productrepo "product_service/internal/repo/productRepo"
	"product_service/internal/utils"
	productpb "product_service/proto/gen"
	"time"

	"github.com/google/uuid"
)

type service struct {
	client client.Client
	repo   productrepo.ProductRepo
}
type Service interface {
	Create(ctx context.Context, payload *productpb.CreateProductRequest, email string) (*productpb.CreateProductResponse, error)
	GetAll(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.GetProductsResponse, error)
	GetById(ctx context.Context, req *productpb.GetProductByIdRequest) (*productpb.GetProductByIdResponse, error)
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
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, productData); err != nil {
		return nil, err
	}
	return utils.ProductResponse(productData), nil
}

func (s *service) GetAll(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.GetProductsResponse, error) {
	log.Println("category", req.Category)
	filters := &productrepo.FilterOptions{
		Category: req.Category,
		Search:   req.Search,
	}
	products, total, err := s.repo.GetAll(ctx, filters)
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
		Products:   pbProducts,
		TotalCount: int32(total),
	}, nil
}

func (s *service) GetById(ctx context.Context, req *productpb.GetProductByIdRequest) (*productpb.GetProductByIdResponse, error) {

	if req.ProductId == "" {
		return nil, errors.New("product id is required")

	}
	product, err := s.repo.GetById(ctx, req.ProductId, req.Category)
	if err != nil {
		return nil, err

	}
	pbProduct := &productpb.Product{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
		ImageUrls:   product.ImageURLs,
		Status:      product.Status,
		IsFeatured:  product.IsFeatured,
		Tags:        product.Tags,
		CreatedBy:   product.CreatedBy,
	}
	return &productpb.GetProductByIdResponse{
		Product: pbProduct,
	}, nil
}
