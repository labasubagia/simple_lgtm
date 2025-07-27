package metrics

import "github.com/prometheus/client_golang/prometheus"

func Init() (*prometheus.CounterVec, *prometheus.HistogramVec) {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path"},
	)
	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_latency_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(requestCounter, latencyHistogram)
	return requestCounter, latencyHistogram
}
