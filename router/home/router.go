package home

import (
	homeController "github.com/ajikamaludin/api-raya-ojt/app/controllers/home"
	"github.com/gofiber/fiber/v2"
)

func HomeRoutes(app *fiber.App) {
	route := app.Group("/")

	route.Get("/", homeController.Home)
}
