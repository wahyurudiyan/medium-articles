package server

import (
	"context"
	productApi "go-products/api/product"
	"go-products/internal/conf"
	"go-products/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cf *conf.Bootstrap, product *service.ProductsService) *http.Server {
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

	cs := cf.Server
	var opts = []http.ServerOption{
		http.Middleware(
			tracing.Server(
				tracing.WithTracerProvider(tp),
			),
			recovery.Recovery(),
		),
	}
	if cs.Http.Network != "" {
		opts = append(opts, http.Network(cs.Http.Network))
	}
	if cs.Http.Addr != "" {
		opts = append(opts, http.Address(cs.Http.Addr))
	}
	if cs.Http.Timeout != nil {
		opts = append(opts, http.Timeout(cs.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	productApi.RegisterProductsHTTPServer(srv, product)

	return srv
}
