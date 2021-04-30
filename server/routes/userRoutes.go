package routes

import (
	"go-react/controllers"

	"github.com/gofiber/fiber/v2"
)

//SetUserRoutes is a function to set up routes for users.
func SetUserRoutes(app *fiber.App) {
	router := app.Group("/api/user")
	router.Get("/", controllers.GetUsers)
	router.Post("/login", controllers.Login)
	router.Post("/create", controllers.CreateUser)
	router.Patch("/update/:id", controllers.UpdateUser)
	router.Delete("/delete/:id", controllers.DeleteUser)
}
