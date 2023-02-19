package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	otrace "go.opentelemetry.io/otel/trace"
)

const (
	// EnvDevelopment is set to development.
	EnvDevelopment = "development"
	// EnvProduction is set to production.
	EnvProduction = "production"

	samplerRatio = 0.2
)

// Config holds configuration for tracing.
type Config struct {
	OtelCollectorAddress string `env:"OPENTELEMETRY_COLLECTOR_ADDRESS,default=localhost:4317"`
	AppEnv               string `env:"APP_ENV,default=development"`
	ServiceName          string `env:"SERVICE_NAME,required"`
}

// Provider provides tracing functionality.
type Provider struct {
	otrace.TracerProvider
}

// NewProvider creates an instance of Provider.
// func NewProvider(cfg Config, exporter sdktrace.SpanExporter) *Provider {
func NewProvider(ctx context.Context, cfg Config) (*Provider, error) {
	sampler := sdktrace.AlwaysSample()
	if cfg.AppEnv == EnvProduction {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplerRatio))
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OtelCollectorAddress),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.AppEnv),
		),
	)
	if err != nil {
		return nil, err
	}

	sp := sdktrace.NewBatchSpanProcessor(exporter)
	prov := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithSpanProcessor(sp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(prov)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	prov.Tracer(cfg.ServiceName)
	return &Provider{prov}, nil
}

// GetTraceIDFromContext gets trace id from context.
func GetTraceIDFromContext(ctx context.Context) string {
	return otrace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
