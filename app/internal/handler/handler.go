package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple_lgtm/internal/model"
	"simple_lgtm/internal/pkg/errs"
	"simple_lgtm/internal/pkg/http_handler"
	"simple_lgtm/internal/service"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
		http_handler.AbortJSON(w, errs.NewInvalidInput(fmt.Errorf("invalid request payload: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}
	if err := payload.Validate(); err != nil {
		http_handler.AbortJSON(w, errs.NewInvalidInput(fmt.Errorf("validation error: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.CreateData(ctx, payload.ID, payload.Value)
	if err != nil {
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	http_handler.JSON(w, http.StatusCreated, model.Response{
		Message: "Data created successfully",
	})
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
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	data, err := h.service.GetData(ctx, id)
	if err != nil {
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	http_handler.JSON(w, http.StatusOK, model.Response{
		Message: "Data retrieved successfully",
		Data:    data,
	})
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
		http_handler.AbortJSON(w, errs.NewInvalidInput(fmt.Errorf("invalid request payload: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}
	payload.ID = id
	if err := payload.Validate(); err != nil {
		http_handler.AbortJSON(w, errs.NewInvalidInput(fmt.Errorf("validation error: %s", err.Error())))
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	span.SetAttributes(
		attribute.String("request.id", payload.ID),
		attribute.String("request.value", payload.Value),
	)

	err := h.service.UpdateData(ctx, payload.ID, payload.Value)
	if err != nil {
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	http_handler.JSON(w, http.StatusOK, model.Response{
		Message: "Data updated successfully",
	})
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
		http_handler.AbortJSON(w, errs.NewInvalidInput(fmt.Errorf("ID parameter is required")))
		return
	}

	span.SetAttributes(attribute.String("request.id", id))

	err := h.service.DeleteData(ctx, id)
	if err != nil {
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	http_handler.JSON(w, http.StatusOK, model.Response{
		Message: "Data deleted successfully",
	})
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
		http_handler.AbortJSON(w, err)
		span.RecordError(err, trace.WithAttributes(attribute.String("error.message", err.Error())))
		return
	}

	http_handler.JSON(w, http.StatusOK, model.Response{
		Message: "Data retrieved successfully",
		Data:    data,
	})
}
