package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"simple_lgtm/internal/config"
	"simple_lgtm/internal/handler"
	"simple_lgtm/internal/repository"
	"simple_lgtm/internal/service"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initMetrics() (*prometheus.CounterVec, *prometheus.HistogramVec) {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path"},
	)
	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_latency_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	prometheus.MustRegister(requestCounter, latencyHistogram)
	return requestCounter, latencyHistogram
}

func initTracer(ctx context.Context, cfg *config.Config) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.AppName),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func main() {

	var loggerHandler slog.Handler = slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	)
	loggerHandler = NewTraceHandler(loggerHandler)
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	requestCounter, latencyHistogram := initMetrics()
	ctx := context.Background()
	shutdownTracer := initTracer(ctx, cfg)
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	repo := repository.NewInMemoryRepository()
	svc := service.NewAppService(repo)
	handler := handler.NewHandler(svc, requestCounter, latencyHistogram)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /data", otelhttp.NewHandler(http.HandlerFunc(handler.ListAllDataHandler), "list").ServeHTTP)
	mux.HandleFunc("GET /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.GetDataHandler), "get").ServeHTTP)
	mux.HandleFunc("POST /data", otelhttp.NewHandler(http.HandlerFunc(handler.CreateDataHandler), "create").ServeHTTP)
	mux.HandleFunc("PUT /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.UpdateDataHandler), "update").ServeHTTP)
	mux.HandleFunc("DELETE /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.DeleteDataHandler), "delete").ServeHTTP)

	mux.HandleFunc("GET /metrics", promhttp.Handler().ServeHTTP)

	slog.Info("app started", slog.Any("port", cfg.Port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		return
	}
}
