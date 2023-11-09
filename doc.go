// Package promnatsfiber provides a suite of observability tools for applications
// using the Fiber framework, with support for Prometheus metrics and NATS messaging system.
// It enables easy collection and exposition of critical operational metrics such as CPU usage,
// memory usage, request counts, request durations, and more.
//
// The package offers middleware for the Fiber framework that automatically collects
// and reports these metrics to a Prometheus server. Additionally, it provides tools
// to monitor NATS messaging system interactions, capturing metrics related to message
// processing.
//
// To use this package, create a new instance of MetricsCollector using NewMetricsCollector,
// which collects a predefined set of system and application metrics. You can then register
// this collector with your Fiber application using the provided middleware functions.
// The metrics are exposed via a configurable endpoint (by default '/metrics') for Prometheus
// to scrape.
//
// The MetricsCollector struct provides methods to manually collect metrics and the middleware
// automatically collects basic request metrics. The package also exposes Go runtime metrics
// such as garbage collection statistics and goroutine counts.
//
// For detailed usage and configuration, refer to the README.md and the examples provided in the
// 'examples' directory of this package.
//
// Example:
//
//	package main
//
//	import (
//	    "github.com/gofiber/fiber/v2"
//	    "github.com/todesdev/promnatsfiber"
//	)
//
//	func main() {
//	    app := fiber.New()
//
//	    // Create a new MetricsCollector instance
//	    mc := promnatsfiber.NewMetricsCollector()
//
//	    // Setup your Fiber app routes as needed
//	    app.Get("/", func(c *fiber.Ctx) error {
//	        return c.SendString("Hello, World!")
//	    })
//
//	    // Register the metrics middleware
//	    app.Use(mc.Middleware())
//
//	    // Start your Fiber application with metrics collection enabled
//	    app.Listen(":8080")
//	}
//
// This package assumes that a Prometheus server is configured to scrape the metrics
// from the '/metrics' endpoint exposed by your Fiber application.
//
// For further details on configuring Prometheus and visualizing metrics with tools
// like Grafana, please refer to the Prometheus and Grafana documentation.
package promnatsfiber
