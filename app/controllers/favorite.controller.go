package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/converter"
	"github.com/gofiber/fiber/v2"
)

type FavoriteController struct {
	Service *services.Services
}

func (favorite *FavoriteController) GetAllAccountFavoriteUser(c *fiber.Ctx) error {
	query := c.Query("query")

	var favorites []models.BankAccountFavorite
	userId := favorite.Service.JwtManager.GetUserId(c)
	favorite.Service.Repository.GetAllAccountFavoriteUser(userId, query, &favorites)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get list user favorite accounts",
		"data":    converter.MapBankAccountFavoriteToRes(favorites),
	})
}

func (favorite *FavoriteController) Store(c *fiber.Ctx) error {
	bankAccountFavoriteReq := new(models.BankAccountFavoriteReq)

	c.BodyParser(&bankAccountFavoriteReq)

	errors := favorite.Service.Validator.ValidateRequest(bankAccountFavoriteReq)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	var bank models.Bank

	err := favorite.Service.Repository.GetBankById(bankAccountFavoriteReq.BankID, &bank)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "bank id not registered",
			"error":   err.Error(),
		})
	}

	userId := favorite.Service.JwtManager.GetUserId(c)
	var bankAccountFavorite models.BankAccountFavorite

	err = favorite.Service.Repository.GetAccountFavoriteUserByAccountNumber(
		userId, bankAccountFavoriteReq.AccountNumber, &bank, &bankAccountFavorite,
	)

	if (models.BankAccountFavorite{}) != bankAccountFavorite {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "account has been set favorite",
		})
	}

	bankAccountFavorite = models.BankAccountFavorite{
		UserId:        userId,
		BankId:        bank.ID,
		Name:          bankAccountFavoriteReq.Name,
		AccountNumber: bankAccountFavoriteReq.AccountNumber,
	}

	if bank.Code == constants.RAYA_BANK_CODE {
		var account models.Account
		err = favorite.Service.Repository.GetBankRayaAccount(bankAccountFavoriteReq.AccountNumber, &account)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  constants.STATUS_FAIL,
				"message": "account number not found",
				"error":   err.Error(),
			})
		}
		bankAccountFavorite.AccountId = account.ID
		bankAccountFavorite.RayaAccount = &account
	} else {
		var bankAccount models.BankAccount
		err = favorite.Service.Repository.GetBankAccountByAccountNumber(bankAccountFavoriteReq.AccountNumber, bank, &bankAccount)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  constants.STATUS_FAIL,
				"message": "account number not found",
				"error":   err.Error(),
			})
		}
		bankAccountFavorite.BankAccountId = bankAccount.ID
		bankAccountFavorite.Bankaccount = &bankAccount
	}

	favorite.Service.Repository.CreateAccountFavoriteUser(&bankAccountFavorite)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "create user favorite accounts",
		"data":    bankAccountFavorite.ToBankAccountFavoriteRes(),
	})
}

func (favorite *FavoriteController) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	var bankAccountFavorite models.BankAccountFavorite
	err := favorite.Service.Repository.GetAccountFavoriteUserBy(id, &bankAccountFavorite)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "favorite account not found",
			"error":   err.Error(),
		})
	}

	favorite.Service.Repository.DeleteAccountFavoriteUser(&bankAccountFavorite)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "delete user favorite accounts",
	})
}
