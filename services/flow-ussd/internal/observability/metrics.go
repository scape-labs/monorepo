// Package observability owns the Prometheus + OTel instruments for flow-ussd.
//
// TODO: declare the request counter / latency histogram / error counter that
// every handler should populate via kit/observability middleware.
package observability

import "github.com/prometheus/client_golang/prometheus"

// TODO: define real metrics. Use the kit/observability registry so they
// surface on /metrics alongside platform-level metrics.
var (
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow-ussd_requests_total",
			Help: "Total HTTP requests served by flow-ussd.",
		},
		[]string{"route", "method", "status"},
	)
)
