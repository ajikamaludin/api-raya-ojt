package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/converter"
	"github.com/gofiber/fiber/v2"
)

type BankController struct {
	Service *services.Services
}

func (bank *BankController) GetAllBanks(c *fiber.Ctx) error {
	query := c.Query("query")

	var banks []models.Bank
	bank.Service.Repository.GetAllBanks(&banks, query)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get list banks",
		"data":    converter.MapBanksToBankRes(banks),
	})
}

func (bank *BankController) ValidateAccountNumber(c *fiber.Ctx) error {
	return nil
}
