package monitor

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    // requestsTotal is the Prometheus counter for total requests
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "myapp_requests_total",
            Help: "Total number of HTTP requests processed",
        },
        []string{"method", "status"},
    )
    
    // requestDuration is the Prometheus histogram for request duration
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "myapp_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "status"},
    )
)

// RegisterMetrics registers Prometheus metrics.
func RegisterMetrics() {
    // Register metrics
    prometheus.MustRegister(requestsTotal)
    prometheus.MustRegister(requestDuration)
}

// MetricsHandler returns an HTTP handler function that serves Prometheus metrics.
func MetricsHandler() http.Handler {
    return promhttp.Handler()
}

// RequestsTotal returns the Prometheus counter for total requests
func RequestsTotal() *prometheus.CounterVec {
    return requestsTotal
}

// RequestDuration returns the Prometheus histogram for request duration
func RequestDuration() *prometheus.HistogramVec {
    return requestDuration
}
