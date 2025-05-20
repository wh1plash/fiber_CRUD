package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthy(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON("result", "ok")
	// })
	checkHandler := NewCheckHandler()
	app.Get("/", checkHandler.HandleHealthy)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

}
