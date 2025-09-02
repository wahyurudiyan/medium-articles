package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	oteltracer "github.com/wahyurudiyan/go-grpc-otel/internal/otel"
	pb "github.com/wahyurudiyan/go-grpc-otel/proto"
)

var tracer = otel.Tracer("client")

func main() {
	// Initialize OpenTelemetry
	tp, err := oteltracer.InitTracer("grpc-client")
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := oteltracer.Shutdown(context.Background(), tp); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Connect to the user service with OTel stats handler
	conn, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Demonstrate different operations with context propagation
	demonstrateContextPropagation(client)
}

func demonstrateContextPropagation(client pb.UserServiceClient) {
	ctx := context.Background()

	log.Println("=== OpenTelemetry Context Propagation Demo ===")

	// 1. Simple GetUser call
	demonstrateGetUser(ctx, client)

	// 2. CreateUser call
	demonstrateCreateUser(ctx, client)

	// 3. Complex GetUserProfile call (demonstrates cross-service propagation)
	demonstrateGetUserProfile(ctx, client)
}

func demonstrateGetUser(ctx context.Context, client pb.UserServiceClient) {
	// Start a root span for the client operation
	ctx, span := tracer.Start(ctx, "client.get_user_demo",
		trace.WithAttributes(
			attribute.String("demo.operation", "get_user"),
			attribute.String("client.type", "demonstration"),
		),
	)
	defer span.End()

	log.Println("\n--- Demonstrating GetUser ---")
	span.AddEvent("demo.started", trace.WithAttributes(
		attribute.String("operation", "get_user"),
	))

	// Call GetUser - the context will be propagated to the server
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 123})
	if err != nil {
		span.RecordError(err)
		log.Printf("GetUser failed: %v", err)
		return
	}

	span.AddEvent("user.retrieved", trace.WithAttributes(
		attribute.String("user.name", resp.User.Name),
		attribute.String("user.email", resp.User.Email),
	))

	log.Printf("Retrieved user: %s (%s)", resp.User.Name, resp.User.Email)
}

func demonstrateCreateUser(ctx context.Context, client pb.UserServiceClient) {
	ctx, span := tracer.Start(ctx, "client.create_user_demo",
		trace.WithAttributes(
			attribute.String("demo.operation", "create_user"),
		),
	)
	defer span.End()

	log.Println("\n--- Demonstrating CreateUser ---")

	resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Age:   28,
	})
	if err != nil {
		span.RecordError(err)
		log.Printf("CreateUser failed: %v", err)
		return
	}

	span.AddEvent("user.created", trace.WithAttributes(
		attribute.Int64("user.id", resp.User.Id),
	))

	log.Printf("Created user: %s (ID: %d)", resp.User.Name, resp.User.Id)
}

func demonstrateGetUserProfile(ctx context.Context, client pb.UserServiceClient) {
	// This demonstrates the most interesting case: cross-service context propagation
	ctx, span := tracer.Start(ctx, "client.get_user_profile_demo",
		trace.WithAttributes(
			attribute.String("demo.operation", "get_user_profile"),
			attribute.String("demo.complexity", "cross_service"),
		),
	)
	defer span.End()

	log.Println("\n--- Demonstrating GetUserProfile (Cross-service call) ---")
	log.Println("This will call UserService.GetUserProfile, which internally calls:")
	log.Println("1. UserService.GetUser")
	log.Println("2. ProfileService.GetProfile (different service!)")
	log.Println("Watch how the trace context propagates across both services!")

	span.AddEvent("cross_service_demo.started")

	// Add some delay to make the trace more visible
	time.Sleep(100 * time.Millisecond)

	resp, err := client.GetUserProfile(ctx, &pb.GetUserProfileRequest{
		UserId: 456,
	})
	if err != nil {
		span.RecordError(err)
		log.Printf("GetUserProfile failed: %v", err)
		return
	}

	span.AddEvent("user_profile.retrieved", trace.WithAttributes(
		attribute.Bool("profile.available", resp.Profile != nil),
	))

	log.Printf("Retrieved user profile for: %s", resp.User.Name)
	if resp.Profile != nil {
		log.Printf("Profile bio: %s", resp.Profile.Bio)
		log.Printf("Interests: %v", resp.Profile.Interests)
	} else {
		log.Println("Profile service was unavailable")
	}

	log.Println("\n=== Trace Information ===")
	spanContext := trace.SpanContextFromContext(ctx)
	log.Printf("Trace ID: %s", spanContext.TraceID().String())
	log.Printf("Span ID: %s", spanContext.SpanID().String())
	log.Println("Check Jaeger UI at http://localhost:16686 to see the complete trace!")
}
