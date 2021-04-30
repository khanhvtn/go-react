package routes

import (
	"go-react/controllers"

	"github.com/gofiber/fiber/v2"
)

//SetEventTypeRoutes func is to set up routes for event type.
func SetEventTypeRoutes(app *fiber.App) {
	router := app.Group("/api/eventType")
	router.Get("/", controllers.GetEventTypes)
	router.Get("/:id", controllers.GetEventType)
	router.Post("/create", controllers.CreateEventType)
	router.Patch("/update/:id", controllers.UpdateEventType)
	router.Delete("/delete/:id", controllers.DeleteEventType)
}
