package main

import (
	"fmt"
	"log"

	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/router"
	apiv1 "github.com/ajikamaludin/api-raya-ojt/router/api/v1"
	home "github.com/ajikamaludin/api-raya-ojt/router/home"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	services := services.New()
	defer services.Db.Conn.Close()

	// default here : /
	home.Routes(app, services)
	// api route : api/v1
	apiv1.Routes(app, services)
	// handle 404
	router.Routes(app)

	config := configs.GetInstance()
	listen := fmt.Sprintf(":%v", config.Appconfig.Port)
	log.Fatal(app.Listen(listen))
}
