package route

import (
	"eventhub/helper"
	"eventhub/middleware"

	"github.com/gofiber/fiber/v2"
)

type Dependencies struct {

}

func Setup(app *fiber.App, deps Dependencies) {
	app.Get("/health", healthCheck)

	api := app.Group("/api")
	api.Use(middleware.RequireJSON)

	v1 := api.Group("/v1")
}

func healthCheck(c *fiber.Ctx) error {
	return helper.Success(c, fiber.StatusOK, "service berjalan", fiber.Map{
		"service": "EventHub",
		"version": "v1",
	})
}