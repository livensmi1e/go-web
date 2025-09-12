package http

import "github.com/prometheus/client_golang/prometheus"

var (
	httpRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go-web",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests processed, labeled by status code and method.",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequest)
}
