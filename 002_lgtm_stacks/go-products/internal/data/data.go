package data

import (
	"database/sql"
	"go-products/ent"
	"go-products/internal/conf"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewProductRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	conf  *conf.Data
	dbCli *ent.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	conn, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(c.Database.Driver, conn)
	if err := drv.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	dbCli := ent.NewClient(ent.Driver(drv))
	d := &Data{
		conf:  c,
		dbCli: dbCli,
	}

	cleanup := func() {
		if err := drv.Close(); err != nil {
			log.Fatalw("data cleanup failed", map[string]string{"error": err.Error()})
		}
		log.NewHelper(logger).Info("closing the data resources")
	}

	return d, cleanup, nil
}
