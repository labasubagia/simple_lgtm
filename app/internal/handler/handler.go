package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple_lgtm/internal/model"
	"simple_lgtm/pkg/errs"
	"simple_lgtm/pkg/http_handler"

	"simple_lgtm/internal/service"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	service          *service.Service
	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

func NewHandler(svc *service.Service, requestCounter *prometheus.CounterVec, latencyHistogram *prometheus.HistogramVec) *Handler {
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
		http_handler.AbortJSON(ctx, w, errs.NewInvalidInput(fmt.Errorf("invalid request payload: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "invalid request payload")
		return
	}
	if err := payload.Validate(); err != nil {
		http_handler.AbortJSON(ctx, w, errs.NewInvalidInput(fmt.Errorf("validation error: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "validation error")
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.CreateData(ctx, payload.ID, payload.Value)
	if err != nil {
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "failed to create data")
		return
	}

	http_handler.JSON(ctx, w, http.StatusCreated, "Data created successfully", nil)
	span.SetStatus(codes.Ok, "success")
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
		err := errs.NewInvalidInput(fmt.Errorf("ID parameter is required"))
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "ID parameter is required")
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	data, err := h.service.GetData(ctx, id)
	if err != nil {
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "failed to get data")
		return
	}

	http_handler.JSON(ctx, w, http.StatusOK, "ok", data)
	span.SetStatus(codes.Ok, "success")
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
		http_handler.AbortJSON(ctx, w, errs.NewInvalidInput(fmt.Errorf("invalid request payload: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "invalid request payload")
		return
	}
	payload.ID = id
	if err := payload.Validate(); err != nil {
		http_handler.AbortJSON(ctx, w, errs.NewInvalidInput(fmt.Errorf("validation error: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "validation error")
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.UpdateData(ctx, payload.ID, payload.Value)
	if err != nil {
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "failed to update data")
		return
	}

	http_handler.JSON(ctx, w, http.StatusOK, "Data updated successfully", nil)
	span.SetStatus(codes.Ok, "success")
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
		http_handler.AbortJSON(ctx, w, errs.NewInvalidInput(fmt.Errorf("ID parameter is required")))
		span.SetStatus(codes.Error, "ID parameter is required")
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	err := h.service.DeleteData(ctx, id)
	if err != nil {
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "failed to delete data")
		return
	}

	http_handler.JSON(ctx, w, http.StatusOK, "Data deleted successfully", nil)
	span.SetStatus(codes.Ok, "success")
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
		http_handler.AbortJSON(ctx, w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		span.SetStatus(codes.Error, "Failed to list data")
		return
	}

	http_handler.JSON(ctx, w, http.StatusOK, "ok", data)
	span.SetStatus(codes.Ok, "success")
}
