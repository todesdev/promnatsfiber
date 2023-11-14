# promnatsfiber

`promnatsfiber` is a Go package designed to seamlessly integrate Prometheus metrics collection into applications using
the Fiber web framework and NATS messaging system.
It offers easy-to-use middleware for Fiber to collect HTTP metrics and a wrapper for NATS message handlers for message
processing metrics.
The package also includes system metrics collection using the `gopsutil` library.

## Features

- Fiber Middleware for HTTP metrics collection.
- NATS Handler Wrapper for message processing metrics.
- System Metrics Collector for process-specific metrics (CPU, Memory, etc.).
- Easy integration with Prometheus for monitoring.

## Exposed Metrics

`promnatsfiber` is designed to expose all the minimally required metrics for production use. The following metrics are
exposed by the package:

### HTTP Metrics

| Metric Name                       | Metric Type | Description                                                               |
|-----------------------------------|-------------|---------------------------------------------------------------------------|
| `http_requests_total`             | Counter     | Total number of HTTP requests processed by the Fiber app.                 |
| `http_request_duration_seconds`   | Histogram   | Total duration of HTTP requests processed by the Fiber app.               |
| `http_requests_in_progress_total` | Gauge       | Total number of HTTP requests currently being processed by the Fiber app. |

### NATS Metrics

| Metric Name                                | Metric Type | Description                                                 |
|--------------------------------------------|-------------|-------------------------------------------------------------|
| `nats_messages_processed_total`            | Counter     | Total number of NATS messages processed by the Fiber app.   |
| `nats_message_processing_duration_seconds` | Histogram   | Total duration of NATS messages processed by the Fiber app. |
| `nats_publishing_messages_total`           | Counter     | Total number of NATS messages published by the Fiber app.   |
| `nats_publishing_message_duration_seconds` | Histogram   | Total duration of NATS messages published by the Fiber app. |

### System Metrics

| Metric Name                 | Metric Type    | Description                    |
|-----------------------------|----------------|--------------------------------|
| `system_cput_usage_percent` | Constant Gauge | Current CPU usage percentage.  |
| `system_memory_usage_bytes` | Constant Gauge | Current memory usage in bytes. |
| `system_memory_total_bytes` | Constant Gauge | Total memory in bytes.         |
| `system_gc_stats`           | Constant Gauge | Garbage collection statistics. |
| `system_go_routine_count`   | Constant Gauge | Current number of Go routines. |

## Installation

To install `promnatsfiber`, you need to have Go installed on your machine. You can then use the following command:

```bash
go get github.com/todesdev/promnatsfiber
```

## Usage

Here's a quick example to get you started:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/todesdev/promnatsfiber"
	"github.com/todesdev/promnatsfiber/middleware"
	"log"
)

func main() {
	// Create new Fiber instance
	app := fiber.New()

	// Connect to NATS
	natsConnection := connectToNats()

	// Initialize the metrics collectors, register Fiber middleware, and register the metrics endpoint
	promnatsfiber.New(&promnatsfiber.Config{
		FiberApp:        app,
		ServiceName:     "my-service",
		MetricsEndpoint: "/metrics",
	})

	// Subscribe to NATS messages
	go subscribeToMessages(natsConnection)

	// Start the Fiber app
	log.Fatal(app.Listen(":3000"))
}

func connectToNats() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	return nc
}

// Example of wrapping a NATS message handler with the promnatsfiber middleware
// NATS metrics are collected automatically by using the middleware
func subscribeToMessages(nc *nats.Conn) {
	sub, err := nc.Subscribe("messages", middleware.WrapProcessMessage(messageHandler))
	if err != nil {
		log.Fatal(err)
	}

	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			log.Fatal(err)
		}
	}(sub)
}

func messageHandler(msg *nats.Msg) {
	log.Printf("Received a message: %s\n", string(msg.Data))
}
```

For detailed usage and more examples, refer to [examples](examples):

- [Fiber and Simple NATS subscription](examples/fiber_simple_nats_subscription/example.go)
- [Fiber and Simple NATS publishing](examples/fiber_simple_nats_publishing/example.go)
- [Fiber and JetStream queue subscription](examples/fiber_jetstream_subscription/example.go)
- [Fiber and JetStream publishing](examples/fiber_jetstream_puslishing/example.go)

## Special Acknowledgements

This project was inspired by the [fiberprometheus](https://github.com/ansrivas/fiberprometheus) package. We extend our
gratitude to the authors and contributors of `fiberprometheus` for their invaluable work, which guided the development
of promnatsfiber.

## Acknowledgements

This project makes use of several open-source libraries:

- [Fiber](https://github.com/gofiber/fiber)
- [Prometheus Go client library](https://github.com/prometheus/client_golang)
- [NATS - Go Client](https://github.com/nats-io/nats.go)
- [gopsutil](https://github.com/shirou/gopsutil)

We are grateful to these authors for their contributions to the open-source community.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.MD) file for details.
