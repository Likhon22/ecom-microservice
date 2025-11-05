package product

import (
	"context"
	"errors"
	"log"
	productservice "product_service/internal/services/productService"
	"product_service/internal/utils"
	productpb "product_service/proto/gen"

	"google.golang.org/grpc/metadata"
)

type handler struct {
	productpb.UnimplementedProductServiceServer
	service productservice.Service
}

func NewProductHandler(service productservice.Service) *handler {
	return &handler{
		service: service,
	}

}

func (h *handler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.StandardResponse, error) {

	if err := req.ValidateAll(); err != nil {
		return nil, utils.MapError(err)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapError(errors.New("missing authentication metadata"))
	}
	emails := md.Get("x-user-email")
	if len(emails) == 0 {
		return nil, utils.MapError(errors.New("user email not found in metadata"))
	}

	product, err := h.service.Create(ctx, req, emails[0])
	log.Println(err)
	if err != nil {
		return nil, utils.MapError(err)

	}
	log.Println("email", emails[0])
	return &productpb.StandardResponse{
		Success:    true,
		Message:    "Product created successfully",
		StatusCode: 201,
		Result: &productpb.StandardResponse_ProductData{
			ProductData: product,
		},
	}, nil
}

func (h *handler) GetProduct(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.StandardResponse, error) {

	products, err := h.service.GetAll(ctx, req)
	if err != nil {
		return nil, utils.MapError(err)

	}
	return &productpb.StandardResponse{
		Success:    true,
		Message:    "products fetched successfully",
		StatusCode: 200,
		Result: &productpb.StandardResponse_Products{
			Products: products,
		},
	}, nil
}

func (h *handler) GetProductById(ctx context.Context, req *productpb.GetProductByIdRequest) (*productpb.StandardResponse, error) {

	product, err := h.service.GetById(ctx, req)
	if err != nil {
		return nil, utils.MapError(err)

	}
	return &productpb.StandardResponse{
		Success:    true,
		Message:    "product fetched successfully",
		StatusCode: 200,
		Result: &productpb.StandardResponse_Product{
			Product: product,
		},
	}, nil

}

func (h *handler) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.StandardResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, utils.MapError(err)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapError(errors.New("missing authentication metadata"))
	}
	emails := md.Get("x-user-email")
	if len(emails) == 0 {
		return nil, utils.MapError(errors.New("user email not found in metadata"))
	}
	log.Println(emails[0])
	product, err := h.service.Update(ctx, req, emails[0])
	if err != nil {
		return nil, utils.MapError(err)

	}
	return &productpb.StandardResponse{
		Success:    true,
		Message:    "product update successful",
		StatusCode: 200,
		Result: &productpb.StandardResponse_UpdatedProduct{
			UpdatedProduct: product,
		},
	}, nil

}
