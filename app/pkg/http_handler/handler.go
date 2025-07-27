package http_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"simple_lgtm/pkg/errs"

	"go.opentelemetry.io/otel/trace"
)

type response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
	SpanID  string `json:"span_id,omitempty"`
}

func AbortJSON(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	traceID, spanID := getTraceInfo(ctx)
	status, message := errs.MapHttp(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{
		Message: message,
		TraceID: traceID,
		SpanID:  spanID,
	})
}

func JSON(ctx context.Context, w http.ResponseWriter, status int, message string, data any) {
	traceID, spanID := getTraceInfo(ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{
		Message: message,
		Data:    data,
		TraceID: traceID,
		SpanID:  spanID,
	})
}

func getTraceInfo(ctx context.Context) (traceID string, spanID string) {
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID().String()
		spanID = span.SpanContext().SpanID().String()
	}
	return traceID, spanID
}
