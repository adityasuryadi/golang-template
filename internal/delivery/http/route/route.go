package route

import (
	"order-service/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App             *fiber.App
	OrderController *http.OrderController
}

func (c *RouteConfig) Setup() {
	c.App.Post("/", c.OrderController.Create)
}
