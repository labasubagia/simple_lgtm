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
}

func AbortJSON(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	traceID, _ := getTraceInfo(ctx)
	w.Header().Set("X-Trace-ID", traceID)
	w.Header().Set("Content-Type", "application/json")
	status, message := errs.MapHttp(err)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{
		Message: message,
	})
}

func JSON(ctx context.Context, w http.ResponseWriter, status int, message string, data any) {
	traceID, _ := getTraceInfo(ctx)
	w.Header().Set("X-Trace-ID", traceID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{
		Message: message,
		Data:    data,
	})
}

func getTraceInfo(ctx context.Context) (traceID string, spanID string) {
	if span := trace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID().String()
		spanID = span.SpanContext().SpanID().String()
	}
	return traceID, spanID
}
