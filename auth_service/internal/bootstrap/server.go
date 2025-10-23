package bootstrap

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

type server struct {
	addr       string
	grpcServer *grpc.Server
	lis        net.Listener
}

func NewServer(addr string) (*server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	return &server{
		addr:       addr,
		grpcServer: grpcServer,
		lis:        lis,
	}, nil

}
func (s *server) StartServer() {

	go func() {
		fmt.Printf("gRPC server listening at %s\n", s.addr)
		if err := s.grpcServer.Serve(s.lis); err != nil {
			log.Fatal("failed to server grpc", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("shutting down the server")
	s.grpcServer.GracefulStop()

}
