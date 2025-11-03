package product

import (
	"context"
	productpb "product_service/proto/gen"
)

type handler struct {
	productpb.UnimplementedProductServiceServer
}

func NewProductHandler() *handler {
	return &handler{}

}

func (h *handler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {

	data := req
	return &productpb.CreateProductResponse{
		Name:  data.GetName(),
		Price: data.GetPrice(),
	}, nil

}
