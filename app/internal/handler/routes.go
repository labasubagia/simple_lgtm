package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Routes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /metrics", promhttp.Handler().ServeHTTP)

	mux.HandleFunc("GET /data", otelhttp.NewHandler(http.HandlerFunc(handler.ListAllDataHandler), "ListData").ServeHTTP)
	mux.HandleFunc("GET /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.GetDataHandler), "GetData").ServeHTTP)
	mux.HandleFunc("POST /data", otelhttp.NewHandler(http.HandlerFunc(handler.CreateDataHandler), "CreateData").ServeHTTP)
	mux.HandleFunc("PUT /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.UpdateDataHandler), "UpdateData").ServeHTTP)
	mux.HandleFunc("DELETE /data/{id}", otelhttp.NewHandler(http.HandlerFunc(handler.DeleteDataHandler), "DeleteData").ServeHTTP)
}
