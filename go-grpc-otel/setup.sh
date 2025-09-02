#!/bin/bash

# Setup script for OpenTelemetry gRPC Context Propagation Demo

echo "🚀 Setting up OpenTelemetry gRPC Context Propagation Demo"

# Create proto directory if it doesn't exist
mkdir -p proto

# Install Protocol Buffers compiler if not present (macOS)
if ! command -v protoc &> /dev/null; then
    echo "📦 Installing protoc..."
    if command -v brew &> /dev/null; then
        brew install protobuf
    else
        echo "❌ Please install Protocol Buffers compiler (protoc) manually"
        echo "   Visit: https://grpc.io/docs/protoc-installation/"
        exit 1
    fi
fi

# Install Go protobuf plugins
echo "📦 Installing Go protobuf plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod tidy

# Generate protobuf code
echo "🔧 Generating protobuf code..."
mkdir -p proto/user
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto

echo "✅ Setup complete!"
echo ""
echo "🐳 To run the demo, you'll need Jaeger running:"
echo "   docker run -d --name jaeger \\"
echo "     -p 16686:16686 \\"
echo "     -p 14268:14268 \\"
echo "     jaegertracing/all-in-one:latest"
echo ""
echo "🔥 Then run the services:"
echo "   Terminal 1: go run cmd/profile-service/main.go"
echo "   Terminal 2: go run cmd/user-service/main.go"
echo "   Terminal 3: go run cmd/client/main.go"
echo ""
echo "📊 View traces at: http://localhost:16686"
