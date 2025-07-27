package handler

import (
	"encoding/json"
	"net/http"
	"simple_lgtm/internal/model"
	"simple_lgtm/internal/service"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string `json:"message"`
}

type Handler struct {
	service          service.Service
	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

func NewHandler(svc service.Service, requestCounter *prometheus.CounterVec, latencyHistogram *prometheus.HistogramVec) *Handler {
	return &Handler{
		service:          svc,
		requestCounter:   requestCounter,
		latencyHistogram: latencyHistogram,
	}
}

func (h *Handler) CreateDataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	ctx, span := otel.Tracer("app-tracer").Start(r.Context(), "CreateDataHandler")
	defer span.End()

	h.requestCounter.WithLabelValues(method, path).Inc()
	defer func() {
		h.latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}()

	var payload model.DataItem
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := Response{Message: "Invalid request payload"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}
	if err := payload.Validate(); err != nil {
		res := Response{Message: "Validation error: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.CreateData(ctx, payload.ID, payload.Value)
	if err != nil {
		res := Response{Message: "Error creating data: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	res := Response{Message: "Data created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetDataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	ctx, span := otel.Tracer("app-tracer").Start(r.Context(), "GetDataHandler")
	defer span.End()

	h.requestCounter.WithLabelValues(method, path).Inc()
	defer func() {
		h.latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}()

	id := r.PathValue("id")
	if id == "" {
		res := Response{Message: "ID parameter is required"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	data, err := h.service.GetData(ctx, id)
	if err != nil {
		res := Response{Message: "Error retrieving data: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	res := Response{Message: data}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	ctx, span := otel.Tracer("app-tracer").Start(r.Context(), "UpdateDataHandler")
	defer span.End()

	h.requestCounter.WithLabelValues(method, path).Inc()
	defer func() {
		h.latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}()

	id := r.PathValue("id")
	var payload model.DataItem
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := Response{Message: "Invalid request payload"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	payload.ID = id
	if err := payload.Validate(); err != nil {
		res := Response{Message: "Validation error: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.UpdateData(ctx, payload.ID, payload.Value)
	if err != nil {
		res := Response{Message: "Error updating data: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	res := Response{Message: "Data updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	ctx, span := otel.Tracer("app-tracer").Start(r.Context(), "DeleteDataHandler")
	defer span.End()

	h.requestCounter.WithLabelValues(method, path).Inc()
	defer func() {
		h.latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}()

	id := r.PathValue("id")
	if id == "" {
		res := Response{Message: "ID parameter is required"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	err := h.service.DeleteData(ctx, id)
	if err != nil {
		res := Response{Message: "Error deleting data: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	res := Response{Message: "Data deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) ListAllDataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	ctx, span := otel.Tracer("app-tracer").Start(r.Context(), "ListAllDataHandler")
	defer span.End()

	h.requestCounter.WithLabelValues(method, path).Inc()
	defer func() {
		h.latencyHistogram.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	}()

	data, err := h.service.ListAllData(ctx)
	if err != nil {
		res := Response{Message: "Error listing all data: " + err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
