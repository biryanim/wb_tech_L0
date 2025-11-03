package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	appName   = "order_service"
	namespace = "order"
)

// Metric holds Prometheus metrics for the order service.
type Metric struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metric

// Init initializes Prometheus metrics for the order service.
func Init(_ context.Context) error {
	metrics = &Metric{
		requestCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      appName + "_requests_total",
			Help:      "The total number of HTTP requests.",
		}),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      appName + "_responses_total",
				Help:      "The total number of HTTP responses.",
			}, []string{"status", "status_code", "method", "path"}),
		histogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      appName + "_histogram_response_time",
				Help:      "Response time from server",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			}, []string{"status"}),
	}
	return nil
}

// IncRequestCounter increments the total HTTP request counter.
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

// IncResponseCounter increments the HTTP response counter for the given status, status code, method, and path.
func IncResponseCounter(status, statusCode, method, path string) {
	metrics.responseCounter.WithLabelValues(status, statusCode, method, path).Inc()
}

// HistogramResponseTimeObserve records an HTTP response time observation for the given status.
func HistogramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
