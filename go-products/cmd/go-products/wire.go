//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"go-products/internal/biz"
	"go-products/internal/conf"
	"go-products/internal/data/product"
	"go-products/internal/server"
	"go-products/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

// wireApp init kratos application.
func wireApp(
	*conf.Bootstrap,
	*pgx.Conn,
	trace.Tracer,
	log.Logger,
	// inject repositories
	*product.Queries,
) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			newApp,
			biz.ProviderSet,
			server.ProviderSet,
			service.ProviderSet,
		),
	)
}
