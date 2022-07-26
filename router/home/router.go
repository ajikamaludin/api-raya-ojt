package home

import (
	controller "github.com/ajikamaludin/api-raya-ojt/app/controllers"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/gofiber/fiber/v2"
)

func HomeRoutes(app *fiber.App) {
	route := app.Group("/")

	homeController := controller.HomeController{
		Serv: services.New(),
	}

	route.Get("/", homeController.Home)
}
