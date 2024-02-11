package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetricServerOrPanic(listenPort int) *http.Server {
	prometheusServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", listenPort),
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	go listenAndServeMetrics(prometheusServer)
	return prometheusServer
}

func listenAndServeMetrics(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		logrus.Panic(err, "failed to start Prometheus http listener")
	}
}

func ShutDownMetricsServer(server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Panic(err, "Failed to shutdown prometheus metric server")
	}
}
