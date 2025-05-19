package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type TestHandler struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h TestHandler) HandleHealthy(c *fiber.Ctx) error {
	fmt.Println("check")
	return c.JSON(fiber.Map{"result": "ok"})
}
