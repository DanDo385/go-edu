package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/metrics"
)

// Metrics records HTTP metrics to Prometheus
func Metrics(m *metrics.Metrics) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Increment active requests
			m.HTTPActiveRequests.Inc()
			defer m.HTTPActiveRequests.Dec()

			// Wrap response writer
			rw := NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rw.StatusCode())

			m.HTTPRequestsTotal.WithLabelValues(
				r.Method, r.URL.Path, status,
			).Inc()

			m.HTTPRequestDuration.WithLabelValues(
				r.Method, r.URL.Path, status,
			).Observe(duration)
		})
	}
}
