package product

import (
	"context"
	"log"
	productservice "product_service/internal/services/productService"
	"product_service/internal/utils"
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

	if err := req.ValidateAll(); err != nil {
		log.Println(err)
	}
	log.Println(req)
	product, err := h.service.Create(ctx, req)
	log.Println(err)
	if err != nil {
		return nil, utils.MapError(err)

	}
	return product, nil
}
