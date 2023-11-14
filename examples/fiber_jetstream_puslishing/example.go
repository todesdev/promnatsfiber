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

	app.Post("/publish", func(c *fiber.Ctx) error {
		message := c.BodyRaw()
		err := publishJetStreamMessage(natsConnection, "jetstream.messages", message)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).SendString("Message published")
	})

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

// Example of instrumented publishing a message to JetStream
func publishJetStreamMessage(nc *nats.Conn, subject string, message []byte) error {
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	stream, err := js.AddStream(&nats.StreamConfig{
		Name:     "my-stream",
		Subjects: []string{"jetstream.messages"},
	})
	if err != nil && stream == nil {
		log.Println("Failed to create a stream:", err)
	}

	instrumentedPublish := middleware.WrapPublishJetStreamMessage(js)
	err = instrumentedPublish(subject, message)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
