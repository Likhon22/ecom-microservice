package product

import (
	"context"
	"log"
	productservice "product_service/internal/services/productService"
	productpb "product_service/proto/gen"
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

func (h *handler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {

	data := req
	if err := req.ValidateAll(); err != nil {
		log.Println(err)
	}
	return &productpb.CreateProductResponse{
		Name:  data.GetName(),
		Price: data.GetPrice(),
	}, nil

}
