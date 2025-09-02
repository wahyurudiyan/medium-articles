package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/wahyurudiyan/go-grpc-otel/internal/otel"
	"github.com/wahyurudiyan/go-grpc-otel/internal/server"
	pb "github.com/wahyurudiyan/go-grpc-otel/proto"
)

func main() {
	// Initialize OpenTelemetry
	tp, err := otel.InitTracer("profile-service")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := otel.Shutdown(context.Background(), tp); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Create gRPC server with OTel stats handler
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// Register services
	profileServer := server.NewProfileServer()
	pb.RegisterProfileServiceServer(s, profileServer)

	// Listen on port 9091
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	log.Println("Profile service starting on port 9091...")
	log.Println("Tracing enabled - send traces to Jaeger at http://localhost:16686")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
