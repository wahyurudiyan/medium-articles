package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// This is a simplified demonstration of OpenTelemetry context propagation concepts
// Run this to understand the basics before diving into the gRPC example

func main() {
	// Initialize a simple stdout exporter for demonstration
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("context-propagation-demo"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Demonstrate context propagation
	demonstrateContextPropagation()

	// Flush traces
	tp.Shutdown(context.Background())
}

func demonstrateContextPropagation() {
	tracer := otel.Tracer("demo")
	ctx := context.Background()

	fmt.Println("=== OpenTelemetry Context Propagation Demo ===\n")

	// 1. Start a root span (like a client request)
	ctx, rootSpan := tracer.Start(ctx, "client_request",
		oteltrace.WithAttributes(
			attribute.String("http.method", "POST"),
			attribute.String("http.url", "/api/users/123/profile"),
		),
	)
	defer rootSpan.End()

	fmt.Println("1. 🚀 Started root span: client_request")
	printTraceInfo(ctx, "Client")

	// 2. Simulate calling UserService (context is propagated)
	simulateUserService(ctx)

	fmt.Println("\n=== Key Concepts Explained ===")
	fmt.Println("✅ Trace ID stays the same across all services")
	fmt.Println("✅ Each service creates child spans")
	fmt.Println("✅ Context flows automatically through function calls")
	fmt.Println("✅ In gRPC, context flows through network calls via metadata")
}

func simulateUserService(ctx context.Context) {
	tracer := otel.Tracer("user-service")

	// This span will automatically be a child of the client_request span
	ctx, userSpan := tracer.Start(ctx, "UserService.GetUserProfile",
		oteltrace.WithAttributes(
			attribute.Int64("user.id", 123),
		),
	)
	defer userSpan.End()

	fmt.Println("2. 👤 UserService received request")
	printTraceInfo(ctx, "UserService")

	// UserService calls GetUser internally
	simulateGetUser(ctx)

	// UserService calls ProfileService
	simulateProfileService(ctx)
}

func simulateGetUser(ctx context.Context) {
	tracer := otel.Tracer("user-service")

	ctx, span := tracer.Start(ctx, "UserService.GetUser")
	defer span.End()

	fmt.Println("3. 🔍 UserService.GetUser called internally")
	printTraceInfo(ctx, "UserService.GetUser")

	span.AddEvent("user.found", oteltrace.WithAttributes(
		attribute.String("user.name", "Alice Johnson"),
	))
}

func simulateProfileService(ctx context.Context) {
	tracer := otel.Tracer("profile-service")

	// This simulates a cross-service call - in real gRPC, the context
	// would be extracted from the incoming request metadata
	ctx, profileSpan := tracer.Start(ctx, "ProfileService.GetProfile",
		oteltrace.WithAttributes(
			attribute.String("service", "profile-service"),
		),
	)
	defer profileSpan.End()

	fmt.Println("4. 📝 ProfileService received cross-service call")
	printTraceInfo(ctx, "ProfileService")

	profileSpan.AddEvent("profile.loaded", oteltrace.WithAttributes(
		attribute.String("profile.bio", "Software Engineer passionate about distributed systems"),
	))
}

func printTraceInfo(ctx context.Context, serviceName string) {
	spanContext := oteltrace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		fmt.Printf("   %s - Trace ID: %s, Span ID: %s\n",
			serviceName,
			spanContext.TraceID().String()[:8]+"...",
			spanContext.SpanID().String()[:8]+"...")
	}
}
