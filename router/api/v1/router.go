package v1

import (
	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/controllers"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func ApiRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	services := services.New()
	authController := &controllers.AuthController{
		Service: services,
	}
	bankController := &controllers.BankController{
		Service: services,
	}
	favoriteController := &controllers.FavoriteController{
		Service: services,
	}

	route.Post("/register", authController.Register)
	route.Post("/login", authController.Login)

	routeAuth := route.Group("/", jwtware.New(jwtware.Config{
		SigningKey:   []byte(configs.GetInstance().Jwtconfig.Secret),
		ErrorHandler: authController.ErrorHandler,
	}))

	routeAuth.Post("/validate-account-pin", authController.ValidatePin)
	routeAuth.Get("/banks", bankController.GetAllBanks)
	routeAuth.Post("/check-account-number", bankController.CheckAccountNumber)
	routeAuth.Get("/bank-account-favorites", favoriteController.GetAllAccountFavoriteUser)
	routeAuth.Post("/bank-account-favorites", favoriteController.Store)
	routeAuth.Delete("/bank-account-favorites/:id", favoriteController.Destroy)
}
