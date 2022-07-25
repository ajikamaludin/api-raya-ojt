package router

import (
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  constants.STATUS_FAIL,
				"message": "Not Found",
			})
		},
	)
}
