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
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				p.TotalErrors.Inc()
				_ = apiErr.Code
			}
		}
		p.TotalRequests.Inc()
		time := time.Since(c.Context().Time())
		p.RequestLatency.Observe(float64(time.Milliseconds()))
		return err
	}
}

type PromMetrics struct {
	TotalRequests  prometheus.Counter   `json:"total_requests"`
	RequestLatency prometheus.Histogram `json:"request_latency"`
	TotalErrors    prometheus.Counter   `json:"total_errors"`
	//RequestDurations
}

func NewPromMetrics() *PromMetrics {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_requests",
		Help: "Total number of requests",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_latency",
		Help:    "Request latency in seconds",
		Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
	})
	reqErrCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_errors",
		Help: "Total number of errors",
	})
	return &PromMetrics{
		TotalRequests:  reqCounter,
		RequestLatency: reqLatency,
		TotalErrors:    reqErrCounter,
	}
}
