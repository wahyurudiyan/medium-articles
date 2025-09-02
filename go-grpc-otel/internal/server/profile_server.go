package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	pb "github.com/wahyurudiyan/go-grpc-otel/proto"
)

var profileTracer = otel.Tracer("profile-service")

// ProfileServer implements the ProfileService
type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
}

// NewProfileServer creates a new ProfileServer
func NewProfileServer() *ProfileServer {
	return &ProfileServer{}
}

// GetProfile retrieves a user profile
func (s *ProfileServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// Start a new span - this will be a child of the incoming span due to context propagation
	ctx, span := profileTracer.Start(ctx, "GetProfile",
		trace.WithAttributes(
			attribute.Int64("user.id", req.UserId),
			attribute.String("service", "profile-service"),
		),
	)
	defer span.End()

	log.Printf("ProfileService.GetProfile called for user ID: %d", req.UserId)

	// Check if we have a parent span (context propagation working)
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		span.AddEvent("context.propagated", trace.WithAttributes(
			attribute.String("trace.id", spanContext.TraceID().String()),
			attribute.String("span.id", spanContext.SpanID().String()),
		))
		log.Printf("Trace context propagated - TraceID: %s", spanContext.TraceID().String())
	}

	// Simulate some database/cache lookup time
	time.Sleep(30 * time.Millisecond)

	// Simulate profile lookup
	profile := &pb.Profile{
		UserId:    req.UserId,
		Bio:       fmt.Sprintf("This is the bio for user %d", req.UserId),
		AvatarUrl: fmt.Sprintf("https://avatar.example.com/user%d.jpg", req.UserId),
		Interests: []string{"technology", "music", "travel"},
	}

	// Add custom attributes and events
	span.SetAttributes(
		attribute.String("profile.bio.length", fmt.Sprintf("%d", len(profile.Bio))),
		attribute.Int("profile.interests.count", len(profile.Interests)),
	)

	span.AddEvent("profile.enriched", trace.WithAttributes(
		attribute.String("profile.avatar_url", profile.AvatarUrl),
	))

	span.SetStatus(codes.Ok, "Profile retrieved successfully")

	return &pb.GetProfileResponse{Profile: profile}, nil
}
