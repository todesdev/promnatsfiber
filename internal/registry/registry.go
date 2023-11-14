package registry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/todesdev/promnatsfiber/internal/collectors"
)

type MetricsRegistry struct {
	Registry               *prometheus.Registry
	HttpMetricsCollector   collectors.HttpMetricsCollector
	NatsMetricsCollector   collectors.AsyncMessageBrokerMetricsCollector
	SystemMetricsCollector collectors.SystemMetricsCollector
}

func NewPrometheusRegistry(serviceName, metricsUrl string) *MetricsRegistry {
	registry := prometheus.NewRegistry()

	return &MetricsRegistry{
		Registry:               registry,
		HttpMetricsCollector:   collectors.NewFiberMetricsCollector(registry, serviceName, metricsUrl),
		NatsMetricsCollector:   collectors.NewNatsMetricsCollector(registry, serviceName),
		SystemMetricsCollector: collectors.NewODSystemMetricsCollector(registry, serviceName),
	}
}
