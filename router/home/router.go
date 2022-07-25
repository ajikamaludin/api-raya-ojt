package home

import (
	home "github.com/ajikamaludin/api-raya-ojt/app/controllers/home"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/gofiber/fiber/v2"
)

func HomeRoutes(app *fiber.App) {
	route := app.Group("/")

	homeController := home.HomeController{
		Service: services.New(),
	}

	route.Get("/", homeController.Home)
}
