package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HttpMetricsCollector interface {
	IncRequestCount(statusCode, method, path string)
	ObserveResponseTime(statusCode, method, path string, duration float64)
	IncRequestsInProgress(method, path string)
	DecRequestsInProgress(method, path string)
	GetMetricsUrl() string
}

const (
	HttpSubsystem                   = "http"
	HttpRequestsTotal               = "requests_total"
	HttpRequestsHelp                = "Total number of HTTP requests."
	HttpRequestDurationSeconds      = "request_duration_seconds"
	HttpRequestsDurationSecondsHelp = "Duration of HTTP requests."
	HttpRequestsInProgressTotal     = "requests_in_progress_total"
	HttpRequestsInProgressHelp      = "Number of HTTP requests in progress."
	HttpStatusCodeLabel             = "status_code"
	HttpMethodLabel                 = "method"
	HttpPathLabel                   = "path"
)

type FiberMetricsCollector struct {
	metricsUrl              string
	requestCountMetric      *prometheus.CounterVec
	responseTimeMetric      *prometheus.HistogramVec
	requestsInProgressGauge *prometheus.GaugeVec
}

func NewFiberMetricsCollector(reg *prometheus.Registry, serviceName, metricsUrl string) HttpMetricsCollector {
	requestCountMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(serviceName, HttpSubsystem, HttpRequestsTotal),
			Help: HttpRequestsHelp,
		},
		[]string{HttpStatusCodeLabel, HttpMethodLabel, HttpPathLabel},
	)

	responseTimeMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prometheus.BuildFQName(serviceName, HttpSubsystem, HttpRequestDurationSeconds),
			Help:    HttpRequestsDurationSecondsHelp,
			Buckets: prometheus.DefBuckets,
		},
		[]string{HttpStatusCodeLabel, HttpMethodLabel, HttpPathLabel},
	)

	requestsInProgressGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(serviceName, HttpSubsystem, HttpRequestsInProgressTotal),
			Help: HttpRequestsInProgressHelp,
		},
		[]string{HttpMethodLabel, HttpPathLabel},
	)

	reg.MustRegister(requestCountMetric, responseTimeMetric)

	return &FiberMetricsCollector{
		metricsUrl:              metricsUrl,
		requestCountMetric:      requestCountMetric,
		responseTimeMetric:      responseTimeMetric,
		requestsInProgressGauge: requestsInProgressGauge,
	}
}

func (m *FiberMetricsCollector) IncRequestCount(statusCode, method, path string) {
	m.requestCountMetric.WithLabelValues(statusCode, method, path).Inc()
}

func (m *FiberMetricsCollector) ObserveResponseTime(statusCode, method, path string, duration float64) {
	m.responseTimeMetric.WithLabelValues(statusCode, method, path).Observe(duration)
}

func (m *FiberMetricsCollector) IncRequestsInProgress(method, path string) {
	m.requestsInProgressGauge.WithLabelValues(method, path).Inc()
}

func (m *FiberMetricsCollector) DecRequestsInProgress(method, path string) {
	m.requestsInProgressGauge.WithLabelValues(method, path).Dec()
}

func (m *FiberMetricsCollector) GetMetricsUrl() string {
	return m.metricsUrl
}
