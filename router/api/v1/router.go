package v1

import (
	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/controllers"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Routes(app *fiber.App, services *services.Services) {
	route := app.Group("/api/v1")

	authController := &controllers.AuthController{
		Serv: services,
	}
	bankController := &controllers.BankController{
		Serv: services,
	}
	favoriteController := &controllers.FavoriteController{
		Serv: services,
	}
	transactionController := &controllers.TransactionController{
		Serv: services,
	}

	route.Post("/register", authController.Register)
	route.Post("/login", authController.Login)

	routeAuth := route.Group("/", jwtware.New(jwtware.Config{
		SigningKey:   []byte(configs.GetInstance().Jwtconfig.Secret),
		ErrorHandler: authController.ErrorHandler,
	}))

	routeAuth.Post("/validate-account-pin", authController.ValidatePin)
	routeAuth.Get("/banks", bankController.GetAllBanks)
	routeAuth.Post("/banks/check-account-number", bankController.CheckAccountNumber)

	routeAuth.Get("/bank-account-favorites", favoriteController.GetAllAccountFavoriteUser)
	routeAuth.Post("/bank-account-favorites", favoriteController.Store)
	routeAuth.Delete("/bank-account-favorites/:id", favoriteController.Destroy)

	routeAuth.Get("/transactions/latest-transactions", transactionController.GetLatestTransactions)
	routeAuth.Get("/transactions/account-balance", transactionController.GetBalance)
	routeAuth.Post("/transactions", transactionController.CreateTransactions)
	routeAuth.Get("/transactions/:id", transactionController.ShowTransaction)
	routeAuth.Post("/test/injetct-balance", authController.InjectBalance)
}
