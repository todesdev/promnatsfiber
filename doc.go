// Package promnatsfiber provides a convenient way to instrument Go applications
// using the Fiber web framework and NATS messaging system with Prometheus metrics.
//
// This package offers middleware for the Fiber framework to collect HTTP metrics
// and a wrapper for NATS message handlers to collect metrics on message processing.
// Additionally, it includes a system metrics collector leveraging gopsutil to
// gather process-specific metrics such as CPU and memory usage, garbage collection
// stats, and Go routine counts.
//
// The primary components of this package are:
//   - PrometheusNatsFiber: A struct that initializes and contains the metrics registry.
//   - MetricsRegistry: A struct holding different types of metrics collectors,
//     including HTTP, NATS, and system metrics collectors.
//   - HttpMetricsCollector: An interface and its implementation for collecting
//     HTTP-related metrics from Fiber requests.
//   - AsyncMessageBrokerMetricsCollector: An interface for collecting metrics
//     on NATS message processing, publishing including JetStream compatability.
//   - SystemMetricsCollector: An interface for collecting system-level metrics.
//
// Usage:
// To use this package, create a new instance of PrometheusNatsFiber with the
// necessary configuration, including the Fiber app instance, service name, and
// metrics endpoint. This initialization will automatically set up the /metrics
// endpoint for Prometheus scraping and register the necessary middleware for
// collecting HTTP metrics in Fiber. For NATS message handling, use the
// WrapWithProm function to wrap your NATS message handlers.
//
// Example:
//
//	func main() {
//	    app := fiber.New()
//	    natsConnection := connectToNats()
//	    promnatsfiber.New(&promnatsfiber.Config{
//	        FiberApp:        app,
//	        ServiceName:     "my-service",
//	        MetricsEndpoint: "/metrics",
//	    })
//	    go subscribeToMessages(natsConnection)
//	    log.Fatal(app.Listen(":3000"))
//	}
//
// This package is designed to be easy to integrate into existing applications
// using Fiber and NATS, providing out-of-the-box metrics collection for
// monitoring with Prometheus.
package promnatsfiber
