package observability

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ccrsxx/api/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

// ServiceNamespace groups every signal emitted by this application. It matches
// the service_namespace label the log/metric collector applies, so logs,
// metrics and traces line up in Grafana.
const ServiceNamespace = "personal-api"

// ShutdownFunc flushes and releases telemetry resources.
type ShutdownFunc func(ctx context.Context) error

func noopShutdown(context.Context) error { return nil }

// InitTracing configures the global OpenTelemetry tracer provider.
//
// It is a no-op (returning a no-op shutdown) when tracing is disabled or no
// OTLP endpoint is configured, so local development and tests need zero
// infrastructure.
//
// Sampling is ParentBased(TraceIDRatioBased): if an upstream caller already
// decided to sample a trace we respect that decision, otherwise we sample a
// ratio of new traces. Keeping errors and slow requests regardless is handled
// by tail sampling in the collector, which needs the whole trace to decide.
func InitTracing(ctx context.Context, cfg config.AppConfig) (ShutdownFunc, error) {
	if !cfg.TracingEnabled || cfg.OtlpEndpoint == "" {
		slog.Info("tracing disabled",
			"tracing_enabled", cfg.TracingEnabled,
			"otlp_endpoint", cfg.OtlpEndpoint,
		)

		return noopShutdown, nil
	}

	// The OTLP HTTP exporter reads OTEL_EXPORTER_OTLP_ENDPOINT itself, but we
	// pass it explicitly so the value always comes from our validated config.
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(cfg.OtlpEndpoint),
	)

	if err != nil {
		return noopShutdown, fmt.Errorf("otlp trace exporter create error: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceNamespace(ServiceNamespace),
			semconv.DeploymentEnvironmentNameKey.String(string(cfg.AppEnv)),
		),
	)

	if err != nil {
		return noopShutdown, fmt.Errorf("otel resource create error: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(cfg.TraceSampleRatio),
			),
		),
	)

	otel.SetTracerProvider(provider)

	// W3C trace context + baggage so trace ids propagate over HTTP in and out.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Never let a telemetry failure take down or spam the application.
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Warn("otel error", "error", err)
	}))

	slog.Info("tracing enabled",
		"otlp_endpoint", cfg.OtlpEndpoint,
		"sample_ratio", cfg.TraceSampleRatio,
	)

	return func(ctx context.Context) error {
		// Give the exporter a bounded window to flush buffered spans.
		flushCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		defer cancel()

		if err := provider.Shutdown(flushCtx); err != nil {
			return fmt.Errorf("tracer provider shutdown error: %w", err)
		}

		return nil
	}, nil
}
