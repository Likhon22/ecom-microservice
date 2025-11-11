package client

import (
	cartpb "cart_service/proto/gen"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	stub cartpb.ProductServiceClient
	conn *grpc.ClientConn
}

type Client interface {
	GetProductById(ctx context.Context, in *cartpb.GetProductByIdRequest) (*cartpb.Product, error)
}

func NewClient(ctx context.Context, clientAddr string) (Client, func() error, error) {

	productClient, err := grpc.DialContext(ctx, clientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)))
	if err != nil {
		return nil, nil, fmt.Errorf("grpc dial %s: %w", clientAddr, err)
	}
	return &client{
		stub: cartpb.NewProductServiceClient(productClient),
		conn: productClient,
	}, productClient.Close, nil

}

func (c *client) GetProductById(ctx context.Context, in *cartpb.GetProductByIdRequest) (*cartpb.Product, error) {

	resp, err := c.stub.GetProductById(ctx, &cartpb.GetProductByIdRequest{
		Category:  in.Category,
		ProductId: in.ProductId,
	})
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("failed to get product: %s", resp.Message)
	}
	productData := resp.GetProduct()
	if productData == nil || productData.Product == nil {
		return nil, fmt.Errorf("product not found in response")
	}
	return productData.Product, nil
}
