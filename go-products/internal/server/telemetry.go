package server

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type ExporterType int

const (
	NoopExporter ExporterType = iota
	ConsoleExporter
	GRPCExporter
	HTTPExporter
)

type TelemetryComponent int

const (
	TelemetryNone  TelemetryComponent = 0
	TelemetryTrace TelemetryComponent = 1 << iota
	TelemetryMetrics
	TelemetryAll = TelemetryTrace | TelemetryMetrics
)

type telemetryOpt struct {
	endpoint     string
	exporterType ExporterType
	components   TelemetryComponent
}

type TelemetryOpt func(*telemetryOpt)

func WithEndpoint(endpoint string) TelemetryOpt {
	return func(opt *telemetryOpt) {
		opt.endpoint = endpoint
	}
}

func WithExporterType(t ExporterType) TelemetryOpt {
	return func(opt *telemetryOpt) {
		opt.exporterType = t
	}
}

func WithComponents(c TelemetryComponent) TelemetryOpt {
	return func(opt *telemetryOpt) {
		opt.components = c
	}
}

func applyTelemetryOpts(opts ...TelemetryOpt) *telemetryOpt {
	opt := &telemetryOpt{
		endpoint:     "localhost:4317",
		exporterType: NoopExporter,
		components:   TelemetryAll,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// === TRACE ===

func (o *telemetryOpt) selectTraceExporter(ctx context.Context) (tracesdk.SpanExporter, error) {
	switch o.exporterType {
	case GRPCExporter:
		return otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(o.endpoint),
			otlptracegrpc.WithInsecure(),
		)
	case HTTPExporter:
		return otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(o.endpoint),
			otlptracehttp.WithInsecure(),
		)
	case ConsoleExporter:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		return nil, nil
	}
}

func initTracer(ctx context.Context, opt *telemetryOpt) (*tracesdk.TracerProvider, error) {
	exp, err := opt.selectTraceExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("init trace exporter: %w", err)
	}
	if exp == nil {
		// fallback to Noop provider
		return nil, nil
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			attribute.String("exporter", "otel"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// === METRICS ===

func (o *telemetryOpt) selectMetricExporter(ctx context.Context) (metricsdk.Exporter, error) {
	switch o.exporterType {
	case GRPCExporter:
		return otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(o.endpoint),
			otlpmetricgrpc.WithInsecure(),
		)
	case HTTPExporter:
		return otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(o.endpoint),
			otlpmetrichttp.WithInsecure(),
		)
	case ConsoleExporter:
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	default:
		return nil, nil
	}
}

func initMeter(ctx context.Context, opt *telemetryOpt) (*metricsdk.MeterProvider, error) {
	exp, err := opt.selectMetricExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("init metric exporter: %w", err)
	}
	if exp == nil {
		// fallback to Noop provider
		return nil, nil
	}

	mp := metricsdk.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp)),
		metric.WithResource(resource.NewSchemaless(
			attribute.String("exporter", "otel"),
		)),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

// === MAIN ENTRY POINT ===

type Telemetry struct {
	opt           *telemetryOpt
	traceprovider *tracesdk.TracerProvider
	meterprovider *metricsdk.MeterProvider
}

func (t *Telemetry) AcquireTraceProvider() (*tracesdk.TracerProvider, error) {
	if t.opt.components&TelemetryTrace == 0 {
		return nil, errors.New("trace provider undefined")
	}

	return t.traceprovider, nil
}

func (t *Telemetry) AcquireMeterProvider() (*metricsdk.MeterProvider, error) {
	if t.opt.components&TelemetryMetrics == 0 {
		return nil, errors.New("trace provider undefined")
	}

	return t.meterprovider, nil
}

func InitTelemetry(ctx context.Context, opts ...TelemetryOpt) (*Telemetry, error) {
	var err error
	var tel Telemetry

	o := applyTelemetryOpts(opts...)
	tel.opt = o

	if o.components&TelemetryTrace != 0 {
		tel.traceprovider, err = initTracer(ctx, o)
		if err != nil {
			return nil, fmt.Errorf("telemetry: failed to init tracer: %w", err)
		}
	}

	if o.components&TelemetryMetrics != 0 {
		tel.meterprovider, err = initMeter(ctx, o)
		if err != nil {
			return nil, fmt.Errorf("telemetry: failed to init meter: %w", err)
		}
	}

	return &tel, nil
}
