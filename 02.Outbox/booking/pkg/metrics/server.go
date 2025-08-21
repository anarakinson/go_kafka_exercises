package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	server *http.Server
}

func NewMetricServer() *MetricServer {
	return &MetricServer{}
}

func (s *MetricServer) Run(port string) error {

	// сбор системных метрик
	startGoroutineMonitor(5 * time.Second)
	startMemoryMonitor(10 * time.Second)
	startGCFreedMemoryMonitor(15 * time.Second)
	startGCMonitor(15 * time.Second)

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: metricsMux,
	}
	return s.server.ListenAndServe()
}

func (s *MetricServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
