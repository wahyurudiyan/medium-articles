#!/bin/bash

# Start all services for the OpenTelemetry gRPC demo

echo "🚀 Starting OpenTelemetry gRPC Context Propagation Demo"

# Check if Jaeger is running
if ! curl -s http://localhost:16686 > /dev/null; then
    echo "⚠️  Jaeger doesn't seem to be running!"
    echo "Start it with: docker run -d --name jaeger -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:latest"
    echo "Continuing anyway..."
fi

# Function to cleanup background processes
cleanup() {
    echo "🛑 Stopping services..."
    jobs -p | xargs -r kill
    exit 0
}

# Set trap to cleanup on exit
trap cleanup SIGINT SIGTERM

echo "📱 Starting Profile Service (port 9091)..."
go run cmd/profile-service/main.go &
PROFILE_PID=$!

echo "👤 Starting User Service (port 9090)..."
go run cmd/user-service/main.go &
USER_PID=$!

# Wait a moment for services to start
sleep 2

echo "🎯 Running Client Demo..."
echo "=================================="
go run cmd/client/main.go

echo ""
echo "✅ Demo completed!"
echo "📊 Check traces at: http://localhost:16686"
echo "🔍 Look for traces from services: grpc-client, user-service, profile-service"

# Keep services running for a bit so you can check Jaeger
echo "⏳ Keeping services running for 30 seconds so you can check Jaeger..."
sleep 30

cleanup
