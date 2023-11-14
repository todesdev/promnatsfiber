package collectors

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
)

type AsyncMessageBrokerMetricsCollector interface {
	IncProcessedMessageCount(subject, messageType string)
	ObserveMessageProcessingDuration(subject, messageType string, duration float64)
	IncPublishedMessageCount(subject, messageType string)
	ObserveMessagePublishingDuration(subject, messageType string, duration float64)
}

const (
	NatsSubsystem = "nats"

	NatsProcessedMessagesTotal        = "processed_messages_total"
	NatsMessagesTotalHelp             = "Total number of NATS messages processed."
	NatsMessageProcessingDuration     = "message_processing_duration_seconds"
	NatsMessageProcessingDurationHelp = "Duration of NATS message processing."

	NatsPublishedMessagesTotal        = "published_messages_total"
	NatsPublishedMessagesHelp         = "Total number of NATS messages published."
	NatsPublishingMessageDuration     = "publishing_message_duration_seconds"
	NatsPublishingMessageDurationHelp = "Duration of NATS message publishing."

	NatsSubjectLabel = "subject"
	NatsTypeLabel    = "type"

	NatsSimpleMessageType    = "simple"
	NatsJetStreamMessageType = "jetstream"
)

var natsMetricsCollector AsyncMessageBrokerMetricsCollector

type NatsMetricsCollector struct {
	processedMessageCountMetric     *prometheus.CounterVec
	messageProcessingDurationMetric *prometheus.HistogramVec

	publishedMessageCountMetric     *prometheus.CounterVec
	messagePublishingDurationMetric *prometheus.HistogramVec
}

func NewNatsMetricsCollector(reg *prometheus.Registry, serviceName string) AsyncMessageBrokerMetricsCollector {
	processedMessageCountMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(serviceName, NatsSubsystem, NatsProcessedMessagesTotal),
			Help: NatsMessagesTotalHelp,
		},
		[]string{NatsSubjectLabel, NatsTypeLabel},
	)

	messageProcessingDurationMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prometheus.BuildFQName(serviceName, NatsSubsystem, NatsMessageProcessingDuration),
			Help:    NatsMessageProcessingDurationHelp,
			Buckets: HistogramBuckets,
		},
		[]string{NatsSubjectLabel, NatsTypeLabel},
	)

	publishedMessageCountMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(serviceName, NatsSubsystem, NatsPublishedMessagesTotal),
			Help: NatsPublishedMessagesHelp,
		},
		[]string{NatsSubjectLabel},
	)

	messagePublishingDurationMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prometheus.BuildFQName(serviceName, NatsSubsystem, NatsPublishingMessageDuration),
			Help:    NatsPublishingMessageDurationHelp,
			Buckets: HistogramBuckets,
		},
		[]string{NatsSubjectLabel},
	)

	reg.MustRegister(
		processedMessageCountMetric,
		publishedMessageCountMetric,
		messageProcessingDurationMetric,
		messagePublishingDurationMetric,
	)

	return &NatsMetricsCollector{
		processedMessageCountMetric:     processedMessageCountMetric,
		messageProcessingDurationMetric: messageProcessingDurationMetric,
		publishedMessageCountMetric:     publishedMessageCountMetric,
		messagePublishingDurationMetric: messagePublishingDurationMetric,
	}
}

func GetNatsMetricsCollector() (AsyncMessageBrokerMetricsCollector, error) {
	if natsMetricsCollector == nil {
		return nil, errors.New("natsMetricsCollector is nil")
	}
	return natsMetricsCollector, nil
}

func (m *NatsMetricsCollector) IncProcessedMessageCount(subject, messageType string) {
	m.processedMessageCountMetric.WithLabelValues(subject, messageType).Inc()
}

func (m *NatsMetricsCollector) ObserveMessageProcessingDuration(subject, messageType string, duration float64) {
	m.messageProcessingDurationMetric.WithLabelValues(subject, messageType).Observe(duration)
}

func (m *NatsMetricsCollector) IncPublishedMessageCount(subject, messageType string) {
	m.publishedMessageCountMetric.WithLabelValues(subject, messageType).Inc()
}

func (m *NatsMetricsCollector) ObserveMessagePublishingDuration(subject, messageType string, duration float64) {
	m.messagePublishingDurationMetric.WithLabelValues(subject, messageType).Observe(duration)
}
