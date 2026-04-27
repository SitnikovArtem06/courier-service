package middleware

import (
	"course-go-avito-SitnikovArtem06/internal/logger"
	"course-go-avito-SitnikovArtem06/internal/observability"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func ObservabilityMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(sw, r)

		duration := time.Since(start).Seconds()

		logger.Log(fmt.Sprintf(" method=%s path=%s status=%d duration=%fs",
			r.Method,
			r.URL.Path,
			sw.status,
			duration))

		observability.HttpRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(sw.status),
		).Inc()

		observability.HttpRequestDuration.WithLabelValues(
			r.URL.Path,
		).Observe(duration)
	})
}
