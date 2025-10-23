package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	addr       string
	grpcServer *grpc.Server
	lis        net.Listener
}

func NewServer(addr string) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	return &Server{
		addr:       addr,
		grpcServer: grpcServer,
		lis:        lis,
	}, nil

}
func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcServer

}
func (s *Server) StartServer(ctx context.Context) {

	go func() {
		fmt.Printf("gRPC server listening at %s\n", s.addr)
		if err := s.grpcServer.Serve(s.lis); err != nil {
			log.Fatal("failed to server grpc", err)
		}
	}()
	<-ctx.Done()
	fmt.Println("shutting down the server")
	s.grpcServer.GracefulStop()

}
