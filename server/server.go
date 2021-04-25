package main

import (
	"context"
	"go-react/database"
	"go-react/routes"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var config = fiber.Config{
	//override default error handle
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		// Statuscode defaults to 500
		code := fiber.StatusInternalServerError

		// Retreive the custom statuscode if it's an fiber.*Error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		// Return from handler
		return ctx.Status(code).JSON(fiber.Map{"message": err.Error()})
	},
}

func main() {
	app := fiber.New(config)
	app.Static("/", "../client/build")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Connect to MongoDB
	mongoClient := database.ConnectMongoDB()
	defer mongoClient.Disconnect(ctx)

	//Middleware
	app.Use(logger.New())

	//Set up routers
	routes.SetUserRoutes(app)

	// Last middleware to match anything
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":5000"))
}
