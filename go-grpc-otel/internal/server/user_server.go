package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/wahyurudiyan/go-grpc-otel/proto"
)

var tracer = otel.Tracer("user-service")

// UserServer implements the UserService
type UserServer struct {
	pb.UnimplementedUserServiceServer
	profileClient pb.ProfileServiceClient
}

// NewUserServer creates a new UserServer with profile service client
func NewUserServer() *UserServer {
	// Connect to profile service (simulating microservice communication)
	conn, err := grpc.Dial("localhost:9091",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Printf("Failed to connect to profile service: %v", err)
		return &UserServer{}
	}

	return &UserServer{
		profileClient: pb.NewProfileServiceClient(conn),
	}
} // GetUser retrieves a user by ID
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Start a new span
	ctx, span := tracer.Start(ctx, "GetUser",
		trace.WithAttributes(
			attribute.Int64("user.id", req.Id),
			attribute.String("operation", "get_user"),
		),
	)
	defer span.End()

	log.Printf("GetUser called with ID: %d", req.Id)

	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Add custom attributes to the span
	span.SetAttributes(
		attribute.String("user.lookup.method", "database"),
		attribute.Bool("user.found", true),
	)

	// Simulate database lookup
	user := &pb.User{
		Id:    req.Id,
		Name:  fmt.Sprintf("User %d", req.Id),
		Email: fmt.Sprintf("user%d@example.com", req.Id),
		Age:   25 + int32(req.Id%10),
	}

	// Add events to trace
	span.AddEvent("user.retrieved", trace.WithAttributes(
		attribute.String("user.name", user.Name),
		attribute.String("user.email", user.Email),
	))

	span.SetStatus(codes.Ok, "User retrieved successfully")

	return &pb.GetUserResponse{User: user}, nil
}

// CreateUser creates a new user
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := tracer.Start(ctx, "CreateUser",
		trace.WithAttributes(
			attribute.String("user.name", req.Name),
			attribute.String("user.email", req.Email),
			attribute.Int("user.age", int(req.Age)),
		),
	)
	defer span.End()

	log.Printf("CreateUser called for: %s", req.Name)

	// Simulate validation
	if req.Name == "" {
		span.RecordError(fmt.Errorf("user name is required"))
		span.SetStatus(codes.Error, "Invalid input")
		return nil, fmt.Errorf("user name is required")
	}

	// Simulate some processing
	time.Sleep(100 * time.Millisecond)

	// Create user (simulate database operation)
	user := &pb.User{
		Id:    time.Now().Unix(), // Simple ID generation
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	span.AddEvent("user.created", trace.WithAttributes(
		attribute.Int64("user.id", user.Id),
	))

	span.SetStatus(codes.Ok, "User created successfully")

	return &pb.CreateUserResponse{User: user}, nil
}

// GetUserProfile demonstrates cross-service communication with context propagation
func (s *UserServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfile",
		trace.WithAttributes(
			attribute.Int64("user.id", req.UserId),
		),
	)
	defer span.End()

	log.Printf("GetUserProfile called for user ID: %d", req.UserId)

	// First, get the user information
	userResp, err := s.GetUser(ctx, &pb.GetUserRequest{Id: req.UserId})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Then, get the profile from another service (context will be propagated)
	if s.profileClient != nil {
		profileResp, err := s.profileClient.GetProfile(ctx, &pb.GetProfileRequest{
			UserId: req.UserId,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to get profile")
			log.Printf("Failed to get profile: %v", err)
			// Continue without profile if service is unavailable
		} else {
			span.AddEvent("profile.retrieved")
			return &pb.GetUserProfileResponse{
				User:    userResp.User,
				Profile: profileResp.Profile,
			}, nil
		}
	}

	// Return without profile if service unavailable
	span.AddEvent("profile.service.unavailable")
	return &pb.GetUserProfileResponse{
		User: userResp.User,
	}, nil
}
