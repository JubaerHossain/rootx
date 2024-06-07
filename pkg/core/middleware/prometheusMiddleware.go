package middleware

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

// PrometheusMiddleware is a middleware for recording Prometheus metrics.
func PrometheusMiddleware(next http.Handler, requestsTotal *prometheus.CounterVec, requestDuration *prometheus.HistogramVec) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Create a response writer to capture the status code
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        statusCode := strconv.Itoa(rw.statusCode)

        // Record metrics
        requestsTotal.With(prometheus.Labels{
            "method": r.Method,
            "status": statusCode,
        }).Inc()

        requestDuration.With(prometheus.Labels{
            "method": r.Method,
            "status": statusCode,
        }).Observe(duration)

        // Optional: Add logging for request details
        log.Printf("Request: method=%s, status=%s, duration=%.6f", r.Method, statusCode, duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    return rw.ResponseWriter.Write(b)
}
