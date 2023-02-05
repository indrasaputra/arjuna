package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
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
	JaegerEndpoint string `env:"JAEGER_ENDPOINT,default=http://localhost:14268/api/traces"`
	AppEnv         string `env:"APP_ENV,default=development"`
	ServiceName    string `env:"SERVICE_NAME,required"`
}

// Provider provides tracing functionality.
type Provider struct {
	otrace.TracerProvider
}

// NewProvider creates an instance of Provider.
func NewProvider(cfg Config, exporter sdktrace.SpanExporter) *Provider {
	sampler := sdktrace.AlwaysSample()
	if cfg.AppEnv == EnvProduction {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplerRatio))
	}

	prov := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.AppEnv),
		)),
	)

	otel.SetTracerProvider(prov)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{prov}
}

// GetTraceIDFromContext gets trace id from context.
func GetTraceIDFromContext(ctx context.Context) string {
	return otrace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

// NewJaegerExporter creates an instance of Jaeger exporter.
func NewJaegerExporter(cfg Config) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}
