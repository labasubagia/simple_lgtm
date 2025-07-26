package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path"},
	)
	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_latency_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func initMetrics() {
	prometheus.MustRegister(requestCounter, latencyHistogram)
}

func initTracer(ctx context.Context) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("lgtm:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-app"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	path := r.URL.Path

	requestCounter.WithLabelValues(path).Inc()
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	latencyHistogram.WithLabelValues(path).Observe(time.Since(start).Seconds())

	log.Printf("Handled request: %s", path)
	w.Write([]byte("Hello from LGTM-integrated Go app"))
}

func main() {
	initMetrics()
	ctx := context.Background()
	defer initTracer(ctx)(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(mainHandler), "main"))
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Go app running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
