package main

import (
	"github.com/kysre/TurtleMQ/leader/internal/pkg/metrics"
	"github.com/kysre/TurtleMQ/leader/internal/pkg/metrics/prometheus"
)

var (
	providerMetrics metrics.Observer
)

func init() {
	providerMetrics = prometheus.NewHistogram("provider",
		"view metrics about provider", "provider_type", "method", "ok", "success")
}

func providePrometheus(config *Config) *prometheus.Server {
	return prometheus.NewServer(config.Metric.ListenPort)
}
