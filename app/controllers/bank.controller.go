package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/converter"
	"github.com/gofiber/fiber/v2"
)

type BankController struct {
	Serv *services.Services
}

func (bank *BankController) GetAllBanks(c *fiber.Ctx) error {
	query := c.Query("query")

	var banks []models.Bank
	bank.Serv.Repository.GetAllBanks(&banks, query)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get list banks",
		"data":    converter.MapBanksToBankRes(banks),
	})
}

func (bank *BankController) CheckAccountNumber(c *fiber.Ctx) error {
	AccountNumberReq := new(models.AccountNumberReq)

	c.BodyParser(&AccountNumberReq)

	errors := bank.Serv.Validator.ValidateRequest(AccountNumberReq)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	accountBank, err := bank.Serv.Repository.CheckAccountNumber(AccountNumberReq.BankId, AccountNumberReq.AccountNumber)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "account bank number not found",
			"error":   err.Error(),
		})
	}

	if accountBank == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "account bank number not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "account bank number found",
		"data":    accountBank,
	})
}
