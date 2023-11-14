package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/todesdev/promnatsfiber"
	"github.com/todesdev/promnatsfiber/middleware"
	"log"
	"time"
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

	// Subscribe to NATS JetStream messages
	go subscribeToJetStreamMessages(natsConnection)

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

func subscribeToJetStreamMessages(nc *nats.Conn) {
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	sub, err := js.QueueSubscribeSync("jetstream.messages", "myqueue", nats.Durable("my-durable-consumer"))
	if err != nil {
		log.Fatal(err)
	}

	// Process messages
	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			log.Println("Failed to receive a message:", err)
			continue
		}

		// Example of instrumented processing a JetStream message
		middleware.WrapProcessJetStreamMessage(jetStreamMessageHandler)(msg)
	}
}

func jetStreamMessageHandler(msg *nats.Msg) {
	log.Printf("Received a message: %s\n", string(msg.Data))
}
