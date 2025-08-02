package data

import (
	"go-products/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	// TODO wrapped database client
	config *conf.Bootstrap
	tracer trace.Tracer
}

// NewData .
func NewData(c *conf.Bootstrap, conn *pgxpool.Conn, tracer trace.Tracer) (*Data, func(), error) {
	cleanup := func() {
		conn.Release()
		log.Info("closing the data resources")
	}

	return &Data{
		config: c,
		tracer: tracer,
	}, cleanup, nil
}
