package home

import (
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/gofiber/fiber/v2"
)

type HomeController struct {
	Service *services.Services
}

func (hc *HomeController) Home(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   constants.STATUS_SUCCESS,
		"message":  "Ok",
		"app_name": hc.Service.Configs.Appconfig.Name,
		"app_env":  hc.Service.Configs.Appconfig.Env,
	})
}
