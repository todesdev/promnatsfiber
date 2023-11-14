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
