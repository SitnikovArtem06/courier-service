package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	RateLimitExceededTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of requests rejected by rate limiter",
		},
		[]string{"path", "method"},
	)

	GatewayRetriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_retries_total",
			Help: "Total number of gateway retry attempts",
		},
		[]string{},
	)
)

func Register() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(RateLimitExceededTotal)
	prometheus.MustRegister(GatewayRetriesTotal)
}
