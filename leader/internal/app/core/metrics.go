package core

import (
	"github.com/kysre/TurtleMQ/leader/internal/pkg/metrics"
	"github.com/kysre/TurtleMQ/leader/internal/pkg/metrics/prometheus"
)

var (
	Latency metrics.Observer
	Total   metrics.Observer
)

func init() {
	Latency = prometheus.NewHistogram("leader_latency",
		"Leader requests latency", "provider", "method", "status")
	Total = prometheus.NewHistogram("leader_total",
		"Leader requests total", "provider", "method", "status")
}
