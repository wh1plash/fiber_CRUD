package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CheckHandler struct{}

func NewCheckHandler() *CheckHandler {
	return &CheckHandler{}
}

func (h CheckHandler) HandleHealthy(c *fiber.Ctx) error {
	fmt.Println("check")
	return c.JSON(fiber.Map{"result": "ok"})
}
