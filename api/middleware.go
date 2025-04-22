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
		errorType := "none"
		start := time.Now()
		errors := make(map[string]string)
		err := handler(c)
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				//fmt.Printf("Request failed with code %d and message: %s\n", apiErr.Code, apiErr.Message)
				status = apiErr.Code
				errors["error"] = apiErr.Message
				errorType = "Api error"
			} else {
				if valErr, ok := err.(ValidationError); ok {
					status = valErr.Status
					errors = valErr.Errors
					errorType = "Validation error"
				} else {
					status = fiber.StatusInternalServerError
					errors["error"] = err.Error()
					errorType = "Internal server error"
				}
			}
		}
		duration := time.Since(start)
		method := c.Method()
		path := c.Path()

		logger.Info("New request:", "method", method, "path", path, "status", status, "errors", errors, "message", errorType, "duration", duration)
		fmt.Println(string(c.Response().Body()))
		fmt.Println("-----------------------------------------------------")
		return err
	}
}

func (p *PromMetrics) WithMetrics(h fiber.Handler, handlerName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := h(c)
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				p.TotalErrors.WithLabelValues(handlerName).Inc()
				_ = apiErr.Code
			}
		}
		p.TotalRequests.WithLabelValues(handlerName).Inc()
		time := time.Since(c.Context().Time())
		p.RequestLatency.WithLabelValues(handlerName).Observe(float64(time.Milliseconds()))
		return err
	}
}

type PromMetrics struct {
	TotalRequests  *prometheus.CounterVec   `json:"total_requests"`
	RequestLatency *prometheus.HistogramVec `json:"request_latency"`
	TotalErrors    *prometheus.CounterVec   `json:"total_errors"`
}

func NewPromMetrics() *PromMetrics {
	reqCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests",
		},
		[]string{"handler"})

	reqLatency := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency",
			Help:    "Request latency in seconds",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
		},
		[]string{"handler"})

	reqErrCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_errors",
			Help: "Total number of errors",
		},
		[]string{"handler"})

	return &PromMetrics{
		TotalRequests:  reqCounter,
		RequestLatency: reqLatency,
		TotalErrors:    reqErrCounter,
	}
}
