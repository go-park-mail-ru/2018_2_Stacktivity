package public_api_server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ApiMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_total_requests",
			Help: "Total number of requests on public api with status codes",
		},
		[]string{"code", "path", "method"},
	)
)
