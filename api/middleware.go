package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingHandlerDecorator(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status := fiber.StatusOK
		message := "OK"
		start := time.Now()
		err := handler(c)
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				fmt.Printf("Request failed with code %d and message: %s\n", apiErr.Code, apiErr.Message)
				status = apiErr.Code
				message = err.Error()
			}
		}
		curTime := time.Now()
		duration := time.Since(start)
		method := c.Method()
		path := c.Path()

		fmt.Printf("%s %s %s %d %s %s\n", curTime, method, path, status, message, duration)

		return err
	}
}
