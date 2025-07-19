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
func NewGRPCServer(conf *conf.Bootstrap, product *service.ProductsService) *grpc.Server {
	ctx := context.Background()
	telemetry, err := InitTelemetry(ctx,
		WithExporterType(ConsoleExporter),
		WithComponents(TelemetryTrace),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := telemetry.AcquireTraceProvider()
	if err != nil {
		log.Fatal(err)
	}

	cs := conf.Server
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			tracing.Server(
				tracing.WithTracerProvider(tp),
			),
			recovery.Recovery(),
		),
	}
	if cs.Grpc.Network != "" {
		opts = append(opts, grpc.Network(cs.Grpc.Network))
	}
	if cs.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(cs.Grpc.Addr))
	}
	if cs.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(cs.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	productApi.RegisterProductsServer(srv, product)

	return srv
}
