package main

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type traceSlogHandler struct {
	baseHandler slog.Handler
}

func NewTraceSlogHandler(baseHandler slog.Handler) slog.Handler {
	return &traceSlogHandler{baseHandler: baseHandler}
}

func (h *traceSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.baseHandler.Handle(ctx, r)
}

func (h *traceSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceSlogHandler{baseHandler: h.baseHandler.WithAttrs(attrs)}
}

func (h *traceSlogHandler) WithGroup(name string) slog.Handler {
	return &traceSlogHandler{baseHandler: h.baseHandler.WithGroup(name)}
}

func (h *traceSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}
