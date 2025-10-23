package usersvc

import (
	userpb "auth_service/proto/gen"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	stub userpb.UserServiceClient
	conn *grpc.ClientConn
}
type Client interface {
	CreateCustomer(ctx context.Context, req *userpb.CreateCustomerRequest) (*userpb.CreateCustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, req *userpb.GetCustomerByEmailRequest) (*userpb.CreateCustomerResponse, error)
	GetCustomers(ctx context.Context, req *userpb.GetCustomersRequest) (*userpb.GetCustomersResponse, error)
	DeleteCustomer(ctx context.Context, req *userpb.DeleteCustomerRequest) (*userpb.DeleteCustomerResponse, error)
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
		stub: userpb.NewUserServiceClient(conn),
		conn: conn,
	}, conn.Close, nil
}

func (c *client) CreateCustomer(ctx context.Context, req *userpb.CreateCustomerRequest) (*userpb.CreateCustomerResponse, error) {
	return c.stub.CreateCustomer(ctx, req)
}

func (c *client) GetCustomerByEmail(ctx context.Context, req *userpb.GetCustomerByEmailRequest) (*userpb.CreateCustomerResponse, error) {
	return c.stub.GetCustomerByEmail(ctx, req)
}

func (c *client) GetCustomers(ctx context.Context, req *userpb.GetCustomersRequest) (*userpb.GetCustomersResponse, error) {
	return c.stub.GetCustomers(ctx, req)
}

func (c *client) DeleteCustomer(ctx context.Context, req *userpb.DeleteCustomerRequest) (*userpb.DeleteCustomerResponse, error) {
	return c.stub.DeleteCustomer(ctx, req)
}
