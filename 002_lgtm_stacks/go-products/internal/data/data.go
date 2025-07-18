package data

import (
	"context"
	"go-products/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
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
func NewData(c *conf.Bootstrap, conn *pgx.Conn, tracer trace.Tracer) (*Data, func(), error) {
	cleanup := func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Fatalw("data cleanup failed", map[string]string{"error": err.Error()})
		}
		log.Info("closing the data resources")
	}

	return &Data{
		config: c,
		tracer: tracer,
	}, cleanup, nil
}
