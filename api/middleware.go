package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func LoggingHandlerDecorator(handler fiber.Handler) fiber.Handler {
	logger := slog.Default()
	return func(c *fiber.Ctx) error {
		status := fiber.StatusOK
		message := "OK"
		start := time.Now()
		err := handler(c)
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				//fmt.Printf("Request failed with code %d and message: %s\n", apiErr.Code, apiErr.Message)
				status = apiErr.Code
				message = err.Error()
			}
		}
		duration := time.Since(start)
		method := c.Method()
		path := c.Path()

		logger.Info("New request:", "method", method, "path", path, "status", status, "message", message, "duration", duration)
		fmt.Println(string(c.Response().Body()))
		fmt.Println("-----------------------------------------------------")
		return err
	}
}

func (p *PromMetrics) WithMetrics(h fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := h(c)
		p.TotalRequests.Inc()
		fmt.Println("Total requests:", p.TotalRequests)
		return err
	}
}

type PromMetrics struct {
	TotalRequests prometheus.Counter `json:"total_requests"`
	//TotalErrors   int64              `json:"total_errors"`
	//RequestDurations
}

func NewPromMetrics() *PromMetrics {
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_requests",
		Help: "Total number of requests",
	})
	return &PromMetrics{
		TotalRequests: counter,
	}
}
