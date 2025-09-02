# 🎯 OpenTelemetry Context Propagation over gRPC - Complete Example

I've created a comprehensive example to help you explain **OpenTelemetry context propagation over gRPC** to your peers. This is a production-ready demonstration showing how distributed tracing works across microservices.

## 📁 Project Structure

```
go-grpc-otel/
├── proto/
│   ├── user.proto              # gRPC service definitions
│   ├── user.pb.go              # Generated protobuf code
│   └── user_grpc.pb.go         # Generated gRPC code
├── internal/
│   ├── otel/
│   │   └── tracer.go           # OpenTelemetry setup & configuration
│   └── server/
│       ├── user_server.go      # Main service with cross-service calls
│       └── profile_server.go   # Secondary service (demonstrates propagation)
├── cmd/
│   ├── user-service/main.go    # User service server (port 9090)
│   ├── profile-service/main.go # Profile service server (port 9091)
│   └── client/main.go          # Demo client with traced calls
├── examples/
│   └── context-propagation-basics/main.go  # Simple concept demo
├── setup.sh                   # Setup script
├── run-demo.sh                 # Complete demo runner
├── docker-compose.yml          # Docker setup with Jaeger
├── README.md                   # Detailed documentation
└── PRESENTATION.md             # Presentation guide for peers
```

## 🚀 Quick Start for Demo

### 1. **Setup Dependencies**
```bash
chmod +x setup.sh && ./setup.sh
```

### 2. **Start Jaeger (for visualization)**
```bash
docker run -d --name jaeger \
  -p 16686:16686 -p 14268:14268 \
  jaegertracing/all-in-one:latest
```

### 3. **Run the Demo**
```bash
chmod +x run-demo.sh && ./run-demo.sh
```

### 4. **View Traces**
Open http://localhost:16686 to see the distributed traces in Jaeger UI.

## 🎓 Key Teaching Points

### 1. **The Magic of Context Propagation**

Show your peers this flow:
```
Client Request
    ├─ UserService.GetUserProfile    [Same Trace ID]
    │   ├─ UserService.GetUser       [Child Span]
    │   └─ ProfileService.GetProfile [Cross-service span]
    └─ All correlated under one trace! 🎉
```

### 2. **Code That Makes It Work**

**Server Setup (Automatic Context Extraction):**
```go
s := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), // 🔑 Magic happens here
)
```

**Client Setup (Automatic Context Injection):**
```go
conn, err := grpc.Dial("localhost:9090",
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()), // 🔑 And here
)
```

**Service Implementation (Context Flows Automatically):**
```go
func (s *UserServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) {
    // ctx already contains the trace context from the client!
    ctx, span := tracer.Start(ctx, "GetUserProfile")
    defer span.End()
    
    // This call to another service propagates context automatically
    profileResp, err := s.profileClient.GetProfile(ctx, profileReq)
}
```

### 3. **What Happens Under the Hood**

1. **Client** starts a trace and makes gRPC call
2. **OTel Client Interceptor** injects trace context into gRPC metadata
3. **Network** carries the metadata to the server
4. **OTel Server Interceptor** extracts trace context from metadata
5. **Server** continues the trace with child spans
6. **Cross-service calls** repeat this process automatically

## 🎯 Demo Script for Your Presentation

### **Setup (5 minutes)**
1. Show the project structure
2. Explain the three services: Client → UserService → ProfileService
3. Start Jaeger and show empty UI

### **Live Demo (10 minutes)**
1. **Run the services** in separate terminals
2. **Execute the client** and watch traces appear in real-time
3. **Open Jaeger UI** and click on a trace
4. **Point out:**
   - Same Trace ID across all services
   - Parent-child span relationships
   - Timing information showing the call flow
   - Custom attributes and events

### **Code Walkthrough (10 minutes)**
1. **Show the interceptors** that make it work automatically
2. **Walk through a service method** showing span creation
3. **Demonstrate custom attributes** and events
4. **Explain error handling** in traces

### **Q&A Points**
- **Performance Impact**: Minimal with proper sampling
- **Production Setup**: Use OTLP exporters, not Jaeger directly
- **Custom Data**: Use Baggage for cross-service data
- **Debugging**: Jaeger shows where trace chains break

## 🛠️ Extending the Example

After your presentation, suggest these exercises:

1. **Add Custom Attributes**: Modify services to add business-specific data
2. **Add Error Scenarios**: Show how errors appear in traces
3. **Add Database Calls**: Demonstrate database span instrumentation
4. **Add Metrics**: Combine tracing with OTel metrics
5. **Add Sampling**: Show different sampling strategies

## 📚 Files for Detailed Study

| File | What Your Peers Should Focus On |
|------|--------------------------------|
| `PRESENTATION.md` | Complete presentation guide with diagrams |
| `internal/server/user_server.go` | Main example of cross-service calls |
| `internal/otel/tracer.go` | OpenTelemetry setup and configuration |
| `cmd/client/main.go` | How to initiate traced calls |
| `examples/context-propagation-basics/main.go` | Simplified concept demonstration |

## 🎉 Why This Example is Great for Teaching

1. **Complete & Realistic**: Shows real microservice architecture
2. **Visual**: Jaeger UI makes concepts immediately visible
3. **Hands-on**: Peers can run it themselves
4. **Production-Ready**: Uses actual OTel best practices
5. **Educational**: Clear comments and documentation

## 🔧 Troubleshooting

If you encounter dependency issues:
```bash
go mod tidy
go get ./...
go mod download
```

If protobuf generation fails:
```bash
# macOS
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 🎯 Key Takeaway for Your Peers

**"OpenTelemetry context propagation is like a baton in a relay race - it automatically passes the trace context from service to service, allowing you to see the complete journey of a request across your entire distributed system."**

---

**Happy Teaching! 🎓✨**

Your peers will understand distributed tracing much better after seeing this live demo!
