package promnatsfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/todesdev/promnatsfiber/internal/registry"
	"github.com/todesdev/promnatsfiber/middleware"
)

type Config struct {
	FiberApp        *fiber.App
	ServiceName     string
	MetricsEndpoint string
}

func New(config *Config) {

	reg := registry.NewPrometheusRegistry(config.ServiceName, config.MetricsEndpoint)

	// Set up the /metrics endpoint for Prometheus scraping using the custom registry
	h := adaptor.HTTPHandler(promhttp.HandlerFor(reg.Registry, promhttp.HandlerOpts{}))
	config.FiberApp.Get(config.MetricsEndpoint, h)

	// Register Fiber middleware
	config.FiberApp.Use(middleware.FiberPrometheusMiddleware(reg.HttpMetricsCollector))
}
