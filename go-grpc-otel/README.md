# OpenTelemetry Context Propagation over gRPC Demo

This project demonstrates **OpenTelemetry context propagation** across gRPC microservices in Go. It's designed to help you understand how distributed tracing works in a microservices architecture.

## 🎯 What This Demo Teaches

### Key Concepts Demonstrated

1. **Context Propagation**: How trace context flows between services
2. **Distributed Tracing**: Correlating operations across multiple services
3. **gRPC Instrumentation**: Automatic and manual span creation
4. **Span Relationships**: Parent-child relationships in distributed calls
5. **Custom Attributes & Events**: Adding business context to traces

### Architecture

```
Client → UserService → ProfileService
   |         |            |
   ├─ Span1  ├─ Span2     ├─ Span3
   └─ TraceID────────────────┘
```

## 📋 Prerequisites

- Go 1.21+
- Docker (for Jaeger)
- Protocol Buffers compiler (`protoc`)

## 🚀 Quick Start

### 1. Setup

```bash
chmod +x setup.sh
./setup.sh
```

### 2. Start Jaeger (for trace visualization)

```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  jaegertracing/all-in-one:latest
```

### 3. Run the Services

In separate terminals:

```bash
# Terminal 1 - Profile Service (port 9091)
go run cmd/profile-service/main.go

# Terminal 2 - User Service (port 9090)
go run cmd/user-service/main.go

# Terminal 3 - Client (demonstrates the calls)
go run cmd/client/main.go
```

### 4. View Traces

Open [http://localhost:16686](http://localhost:16686) in your browser to see the Jaeger UI.

## 🔍 Understanding the Flow

### 1. Simple Call: `GetUser`
```
Client → UserService.GetUser
```
- Single service call
- Basic span creation and attributes

### 2. Service Creation: `CreateUser`
```
Client → UserService.CreateUser
```
- Demonstrates error handling in traces
- Custom events and attributes

### 3. **Cross-Service Call: `GetUserProfile`** ⭐
```
Client → UserService.GetUserProfile
                ├→ UserService.GetUser
                └→ ProfileService.GetProfile
```
This is the **key demo** showing context propagation:
- Client starts a trace
- UserService receives the context and creates child spans
- UserService calls ProfileService, propagating the context
- ProfileService creates spans as children of the original trace
- **All operations are correlated under one trace ID**

## 🎓 Teaching Points for Your Peers

### 1. **Context Propagation Magic**

Show how the trace ID flows automatically:

```go
// In client
ctx, span := tracer.Start(ctx, "client_operation")
client.GetUserProfile(ctx, request) // Context goes to server

// In UserService  
func (s *UserServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) {
    // ctx already contains the trace context!
    ctx, span := tracer.Start(ctx, "GetUserProfile") // Child span
    
    // Call another service - context propagates automatically
    s.profileClient.GetProfile(ctx, profileReq)
}
```

### 2. **gRPC Interceptors**

Explain how OTel interceptors handle the heavy lifting:

```go
// Server side - automatically extracts context from gRPC metadata
grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
)

// Client side - automatically injects context into gRPC metadata  
grpc.Dial(address,
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
)
```

### 3. **Custom Instrumentation**

Show how to add business context:

```go
span.SetAttributes(
    attribute.String("user.email", user.Email),
    attribute.Int64("user.id", user.Id),
)

span.AddEvent("user.validated", trace.WithAttributes(
    attribute.Bool("validation.passed", true),
))
```

## 🔧 Key Files to Review

| File | Purpose |
|------|---------|
| `internal/otel/tracer.go` | OTel setup and configuration |
| `internal/server/user_server.go` | Main service with cross-service calls |
| `internal/server/profile_server.go` | Demonstrates receiving propagated context |
| `cmd/client/main.go` | Shows how to initiate traced calls |

## 🎯 Demonstration Script for Peers

1. **Start the demo** and show Jaeger UI (empty)
2. **Run the client** and watch traces appear in real-time
3. **Click on a trace** to show the span hierarchy
4. **Point out**:
   - Same Trace ID across all services
   - Parent-child span relationships
   - How timing shows the call flow
   - Custom attributes and events
5. **Show the code** that creates spans and propagates context
6. **Modify attributes** and re-run to show customization

## 🤔 Common Questions & Answers

**Q: How does the context actually travel over gRPC?**  
A: It's embedded in gRPC metadata headers automatically by the OTel interceptors.

**Q: What if a service doesn't have OTel?**  
A: The trace stops there, but downstream services can start new traces.

**Q: Can we propagate custom data?**  
A: Yes, using Baggage (for data) and custom attributes (for metadata).

**Q: What's the performance impact?**  
A: Minimal with sampling. Production typically uses 1-5% sampling.

## 🛠 Extending the Demo

Try these modifications to deepen understanding:

1. **Add Baggage**: Propagate user permissions across services
2. **Add Database Spans**: Show how to trace database operations  
3. **Add Errors**: Demonstrate error propagation in traces
4. **Add Sampling**: Show different sampling strategies
5. **Add Metrics**: Combine with OTel metrics collection

## 📚 Further Reading

- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [gRPC Instrumentation](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)

---

**Happy Tracing! 🔍✨**
