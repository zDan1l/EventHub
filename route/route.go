package route

import (
	"eventhub/app/service"
	"eventhub/helper"
	"eventhub/middleware"

	"github.com/gofiber/fiber/v2"
)

type Dependencies struct {
	EventService *service.EventService
}

func Setup(app *fiber.App, deps Dependencies) {
	app.Get("/health", healthCheck)

	api := app.Group("/api")
	api.Use(middleware.RequireJSON)

	v1 := api.Group("/v1")

	events := v1.Group("/events")
	events.Get("/", deps.EventService.GetAllEvents)
	events.Post("/", deps.EventService.CreateEvent)
	events.Get("/:id", deps.EventService.GetEventByID)
	events.Put("/:id", deps.EventService.UpdateEvent)
	events.Delete("/:id", deps.EventService.DeleteEvent)
	
}

func healthCheck(c *fiber.Ctx) error {
	return helper.Success(c, fiber.StatusOK, "service berjalan", fiber.Map{
		"service": "EventHub",
		"version": "v1",
	})
}