package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/todesdev/promnatsfiber/internal/collectors"
	"strconv"
	"time"
)

func FiberPrometheusMiddleware(mc collectors.HttpMetricsCollector) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		method := c.Route().Method

		path := c.Route().Path

		if path == mc.GetMetricsUrl() {
			return c.Next()
		}

		mc.IncRequestsInProgress(method, path)
		defer mc.DecRequestsInProgress(method, path)

		err := c.Next()
		if err != nil {
			return err
		}

		statusCode := c.Response().StatusCode()

		mc.IncRequestCount(strconv.Itoa(statusCode), method, path)
		elapsed := float64(time.Since(startTime).Nanoseconds()) / 1e9
		mc.ObserveResponseTime(strconv.Itoa(statusCode), method, path, elapsed)

		return nil
	}

}
