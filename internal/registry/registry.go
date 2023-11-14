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

	formattedServiceName := toSnakeCase(serviceName)

	return &MetricsRegistry{
		Registry:               registry,
		HttpMetricsCollector:   collectors.NewFiberMetricsCollector(registry, formattedServiceName, metricsUrl),
		NatsMetricsCollector:   collectors.NewNatsMetricsCollector(registry, formattedServiceName),
		SystemMetricsCollector: collectors.NewODSystemMetricsCollector(registry, formattedServiceName),
	}
}
