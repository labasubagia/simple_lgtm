package tracer

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type slogHandler struct {
	baseHandler slog.Handler
}

func NewSlogHandler(baseHandler slog.Handler) slog.Handler {
	return &slogHandler{baseHandler: baseHandler}
}

func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.baseHandler.Handle(ctx, r)
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{baseHandler: h.baseHandler.WithAttrs(attrs)}
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{baseHandler: h.baseHandler.WithGroup(name)}
}

func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}
