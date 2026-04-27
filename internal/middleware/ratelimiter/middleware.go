package ratelimiter

import (
	"course-go-avito-SitnikovArtem06/internal/logger"
	"course-go-avito-SitnikovArtem06/internal/observability"
	"fmt"
	"net/http"
)

func RateLimiterMiddleware(bucket *TokenBucket, logger logger.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bucket.Allow() {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				observability.RateLimitExceededTotal.WithLabelValues(
					r.URL.Path,
					r.Method,
				).Inc()
				logger.Log(fmt.Sprintf("rate limit exceeded method=%s path=%s", r.Method, r.URL.Path))

				_, err := w.Write([]byte("rate limit exceeded"))
				if err != nil {
					return
				}
			}
		})
	}
}
