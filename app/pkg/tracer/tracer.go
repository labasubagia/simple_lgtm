package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Init(ctx context.Context, appName string) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
