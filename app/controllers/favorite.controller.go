package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/converter"
	"github.com/gofiber/fiber/v2"
)

type FavoriteController struct {
	Serv *services.Services
}

func (favorite *FavoriteController) GetAllAccountFavoriteUser(c *fiber.Ctx) error {
	query := c.Query("query")

	var favorites []models.BankAccountFavorite
	userId := favorite.Serv.JwtManager.GetUserId(c)
	favorite.Serv.Repository.GetAllAccountFavoriteUser(userId, query, &favorites)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get list user favorite accounts",
		"data":    converter.MapBankAccountFavoriteToRes(favorites),
	})
}

func (favorite *FavoriteController) Store(c *fiber.Ctx) error {
	bankAccountFavoriteReq := new(models.BankAccountFavoriteReq)

	c.BodyParser(&bankAccountFavoriteReq)

	errors := favorite.Serv.Validator.ValidateRequest(bankAccountFavoriteReq)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	// validate bank is exists
	var bank models.Bank
	err := favorite.Serv.Repository.GetBankById(bankAccountFavoriteReq.BankID, &bank)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "bank id not registered",
			"error":   err.Error(),
		})
	}

	userId := favorite.Serv.JwtManager.GetUserId(c)
	var bankAccountFavorite models.BankAccountFavorite

	// validate bank account is not in favorite
	err = favorite.Serv.Repository.GetAccountFavoriteUserByAccountNumber(
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
		CreatedBy:     userId, // note: this is must be not here, it must be set automate from gorm hooks, find out leter.
	}

	// validate bank is raya or outside bank
	if bank.Code == constants.RAYA_BANK_CODE {
		var account models.Account
		err = favorite.Serv.Repository.GetBankRayaAccount(bankAccountFavoriteReq.AccountNumber, &account)
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
		err = favorite.Serv.Repository.GetBankAccountByAccountNumber(bankAccountFavoriteReq.AccountNumber, bank, &bankAccount)
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

	// create
	favorite.Serv.Repository.CreateAccountFavoriteUser(&bankAccountFavorite)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "create user favorite accounts",
		"data":    bankAccountFavorite.ToBankAccountFavoriteRes(),
	})
}

func (favorite *FavoriteController) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")

	var bankAccountFavorite models.BankAccountFavorite
	err := favorite.Serv.Repository.GetAccountFavoriteUserBy(id, &bankAccountFavorite)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "favorite account not found",
			"error":   err.Error(),
		})
	}

	favorite.Serv.Repository.DeleteAccountFavoriteUser(&bankAccountFavorite)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "delete user favorite accounts",
	})
}
