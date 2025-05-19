package api

import "github.com/gofiber/fiber/v2"

type TestHandler struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h TestHandler) HandleHealthy(c *fiber.Ctx) error {
	return c.JSON("result", "ok")
}
