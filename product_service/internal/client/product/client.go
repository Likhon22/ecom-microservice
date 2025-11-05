package usersvc

import (
	"context"
	"fmt"
	productpb "product_service/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	stub productpb.UserServiceClient
	conn *grpc.ClientConn
}
type Client interface {
	CreateCustomer(ctx context.Context, req *productpb.CreateCustomerRequest) (*productpb.CreateCustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, req *productpb.GetCustomerByEmailRequest) (*productpb.CreateCustomerResponse, error)
	GetCustomers(ctx context.Context, req *productpb.GetCustomersRequest) (*productpb.GetCustomersResponse, error)
	DeleteCustomer(ctx context.Context, req *productpb.DeleteCustomerRequest) (*productpb.DeleteCustomerResponse, error)
	GetCustomerCredentials(ctx context.Context, req *productpb.GetCustomerByEmailRequest) (*productpb.CustomerCredentialsResponse, error)
}

func NewClient(ctx context.Context, addr string) (Client, func() error, error) {
	conn, err := grpc.DialContext(ctx,
		addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc dial %s: %w", addr, err)
	}
	return &client{
		stub: productpb.NewUserServiceClient(conn),
		conn: conn,
	}, conn.Close, nil
}

func (c *client) CreateCustomer(ctx context.Context, req *productpb.CreateCustomerRequest) (*productpb.CreateCustomerResponse, error) {
	return c.stub.CreateCustomer(ctx, req)

}

func (c *client) GetCustomerByEmail(ctx context.Context, req *productpb.GetCustomerByEmailRequest) (*productpb.CreateCustomerResponse, error) {
	return c.stub.GetCustomerByEmail(ctx, req)
}

func (c *client) GetCustomers(ctx context.Context, req *productpb.GetCustomersRequest) (*productpb.GetCustomersResponse, error) {
	return c.stub.GetCustomers(ctx, req)
}

func (c *client) DeleteCustomer(ctx context.Context, req *productpb.DeleteCustomerRequest) (*productpb.DeleteCustomerResponse, error) {
	return c.stub.DeleteCustomer(ctx, req)
}

func (c *client) GetCustomerCredentials(ctx context.Context, req *productpb.GetCustomerByEmailRequest) (*productpb.CustomerCredentialsResponse, error) {
	return c.stub.GetCustomerCredentials(ctx, req)
}
