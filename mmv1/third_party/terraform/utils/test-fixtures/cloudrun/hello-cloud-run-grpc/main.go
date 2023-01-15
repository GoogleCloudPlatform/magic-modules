package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	echopb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type echoServer struct {
	echopb.UnimplementedEchoServer
}

func (e *echoServer) UnaryEcho(ctx context.Context, req *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	return &echopb.EchoResponse{Message: req.Message}, nil
}

// A basic gRPC server with the standard-compliant health check for deployment to Cloud Run
func main() {
	server := grpc.NewServer()

	// Provides the gRPC service grpc.health.v1.Health/Check
	// The implementation is based on the standardized health checking protocol: https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthcheck)

	// Provides the gRPC service grpc.examples.echo.Echo/UnaryEcho
	echopb.RegisterEchoServer(server, &echoServer{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("gRPC service starts listening on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
