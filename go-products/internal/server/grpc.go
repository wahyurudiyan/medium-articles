package server

import (
	"context"
	productApi "go-products/api/product"
	"go-products/internal/conf"
	"go-products/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cf *conf.Bootstrap, product *service.ProductsService) *grpc.Server {
	ctx := context.Background()
	telemetry, err := InitTelemetry(ctx,
		WithComponents(TelemetryTrace),
		WithExporterType(GRPCExporter),
		WithServiceName(cf.Application.Name),
		WithEndpoint(cf.Telemetry.GetEndpoint()),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := telemetry.AcquireTraceProvider()
	if err != nil {
		log.Fatal(err)
	}

	serverConfig := cf.Server
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			tracing.Server(
				tracing.WithTracerProvider(tp),
			),
			recovery.Recovery(),
		),
	}
	if serverConfig.Grpc.Network != "" {
		opts = append(opts, grpc.Network(serverConfig.Grpc.Network))
	}
	if serverConfig.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(serverConfig.Grpc.Addr))
	}
	if serverConfig.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(serverConfig.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	productApi.RegisterProductsServer(srv, product)

	return srv
}
