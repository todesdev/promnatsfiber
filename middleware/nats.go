package middleware

import (
	"github.com/nats-io/nats.go"
	"github.com/todesdev/promnatsfiber/internal/collectors"
	"time"
)

func WrapProcessMessage(funcToWrap func(*nats.Msg)) func(*nats.Msg) {
	return func(msg *nats.Msg) {
		mc, err := collectors.GetNatsMetricsCollector()
		if err != nil {
			panic(err)
		}
		startTime := time.Now()
		funcToWrap(msg)

		mc.IncProcessedMessageCount(msg.Subject, collectors.NatsSimpleMessageType)

		elapsed := float64(time.Since(startTime).Nanoseconds()) / 1e9
		mc.ObserveMessageProcessingDuration(msg.Subject, collectors.NatsSimpleMessageType, elapsed)
	}
}

func WrapProcessJetStreamMessage(funcToWrap func(*nats.Msg)) func(*nats.Msg) {
	return func(msg *nats.Msg) {
		mc, err := collectors.GetNatsMetricsCollector()
		if err != nil {
			panic(err)
		}
		startTime := time.Now()
		funcToWrap(msg)

		mc.IncProcessedMessageCount(msg.Subject, collectors.NatsJetStreamMessageType)

		elapsed := float64(time.Since(startTime).Nanoseconds()) / 1e9
		mc.ObserveMessageProcessingDuration(msg.Subject, collectors.NatsJetStreamMessageType, elapsed)
	}
}

func WrapPublishMessage(nc *nats.Conn) func(string, []byte) error {
	return func(subject string, data []byte) error {
		mc, err := collectors.GetNatsMetricsCollector()
		if err != nil {
			panic(err)
		}
		startTime := time.Now()
		err = nc.Publish(subject, data)
		if err != nil {
			return err
		}

		mc.IncPublishedMessageCount(subject, collectors.NatsSimpleMessageType)

		elapsed := float64(time.Since(startTime).Nanoseconds()) / 1e9
		mc.ObserveMessagePublishingDuration(subject, collectors.NatsSimpleMessageType, elapsed)

		return nil
	}
}

func WrapPublishJetStreamMessage(js nats.JetStreamContext) func(string, []byte) error {
	return func(subject string, data []byte) error {
		mc, err := collectors.GetNatsMetricsCollector()
		if err != nil {
			panic(err)
		}
		startTime := time.Now()
		_, err = js.Publish(subject, data)
		if err != nil {
			return err
		}

		mc.IncPublishedMessageCount(subject, collectors.NatsJetStreamMessageType)

		elapsed := float64(time.Since(startTime).Nanoseconds()) / 1e9
		mc.ObserveMessagePublishingDuration(subject, collectors.NatsJetStreamMessageType, elapsed)

		return nil
	}
}
