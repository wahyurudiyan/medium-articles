package main

import (
	"context"
	"flag"
	"os"

	"go-products/internal/conf"
	"go-products/internal/data/product"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newZapLogger(c *conf.Logger) *kratoszap.Logger {
	var (
		err         error
		output      = os.Stdout
		encoderConf = zap.NewProductionEncoderConfig()
	)

	if c.Output == conf.LogOutput_file {
		output, err = os.OpenFile("/var/log/service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	var encoder = zapcore.NewConsoleEncoder(encoderConf)
	if c.Encoder == conf.LogEncoder_json {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	writeSyncer := zapcore.AddSync(output)
	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	z := zap.New(core)

	return kratoszap.NewLogger(z)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	ctx := context.Background()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	Name = bc.Application.GetName()
	Version = bc.Application.GetVersion()

	tracer := otel.Tracer(bc.Application.Name)
	logger := log.With(
		newZapLogger(bc.GetLogger()),
		"caller", log.DefaultCaller,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(logger)

	db, err := pgx.Connect(ctx, bc.Data.Database.Source)
	if err != nil {
		panic(err)
	}

	// Init service repositories
	productRepository := product.New(db)

	app, cleanup, err := wireApp(
		&bc,
		db,
		tracer,
		logger,
		productRepository,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
