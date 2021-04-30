package routes

import (
	"go-react/controllers"

	"github.com/gofiber/fiber/v2"
)

//SetEventRoutes func is to set up routes for Event.
func SetEventRoutes(app *fiber.App) {
	router := app.Group("/api/event")
	router.Get("/", controllers.GetEvents)
	router.Get("/:id", controllers.GetEvent)
	router.Post("/create", controllers.CreateEvent)
	router.Patch("/update/:id", controllers.UpdateEvent)
	router.Delete("/delete/:id", controllers.DeleteEvent)
}
