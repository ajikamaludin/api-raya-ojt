package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/gofiber/fiber/v2"
)

type HomeController struct {
	Serv *services.Services
}

func (hc *HomeController) Home(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "Ok",
		"data": map[string]string{
			"app_name": hc.Serv.Configs.Appconfig.Name,
			"app_env":  hc.Serv.Configs.Appconfig.Env,
		},
		"error": "",
	})
}
