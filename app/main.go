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
	"simple_lgtm/pkg/metrics"
	"simple_lgtm/pkg/tracer"
)

func main() {

	var loggerHandler slog.Handler = slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	)
	loggerHandler = tracer.NewSlogHandler(loggerHandler)
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	cfg := config.Load()

	requestCounter, latencyHistogram := metrics.Init()
	ctx := context.Background()
	shutdownTracer := tracer.Init(ctx, cfg)
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	repo := repository.NewInMemoryRepository()
	svc := service.NewService(repo)
	hldr := handler.NewHandler(svc, requestCounter, latencyHistogram)

	mux := http.NewServeMux()
	handler.Routes(mux, hldr)

	slog.Info("app started", slog.Any("port", cfg.Port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		return
	}
}
