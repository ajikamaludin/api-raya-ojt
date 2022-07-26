package controllers

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/app/services"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/converter"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	Serv *services.Services
}

func (trx *TransactionController) GetBalance(c *fiber.Ctx) error {
	userId := trx.Serv.JwtManager.GetUserId(c)

	var account models.Account
	trx.Serv.Repository.GetUserBalance(userId, &account)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get user balance",
		"data": fiber.Map{
			"accountNumber": account.AccountNumber,
			"balance":       account.Balance,
		},
	})
}

func (trx *TransactionController) GetLatestTransactions(c *fiber.Ctx) error {
	query := c.Query("query")
	userId := trx.Serv.JwtManager.GetUserId(c)
	var transactions []models.BankTransaction
	err := trx.Serv.Repository.GetLatestTransactions(userId, &transactions, query)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get user latest transaction",
		"data":    converter.MapBankTransactionToRes(transactions),
	})
}

func (trx *TransactionController) CreateTransactions(c *fiber.Ctx) error {
	createBankTransactionReq := new(models.CreateBankTransactionReq)

	c.BodyParser(&createBankTransactionReq)

	errors := trx.Serv.Validator.ValidateRequest(createBankTransactionReq)
	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "request body invalid",
			"error":   errors,
		})
	}

	// validate bank is exists
	var bank models.Bank
	err := trx.Serv.Repository.GetBankById(createBankTransactionReq.BankId, &bank)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "bank id not registered",
			"error":   err.Error(),
		})
	}

	var account models.BankAccount
	err = trx.Serv.Repository.GetBankAccountByAccountNumber(createBankTransactionReq.AccountNumber, bank, &account)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "bank account number not valid",
			"error":   err.Error(),
		})
	}

	userId := trx.Serv.JwtManager.GetUserId(c)

	transaction := models.BankTransaction{
		BankAccountId: account.ID,
		BankId:        bank.ID,
		UserId:        userId,
		Debit:         0,
		Credit:        createBankTransactionReq.Amount,
		Status:        constants.TRX_PENDING,
		CreatedBy:     userId, // note: use gorm hook
	}

	err = trx.Serv.Repository.CreateBankTransaction(&transaction, &bank)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "transaction is invalid",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "transaction created",
		"data":    transaction.ToBankTransactionRes(),
	})
}

func (trx *TransactionController) ShowTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.BankTransaction
	err := trx.Serv.Repository.GetTransactionById(id, &transaction)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  constants.STATUS_FAIL,
			"message": "transaction not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  constants.STATUS_SUCCESS,
		"message": "get transaction detail",
		"data":    transaction.ToBankTransactionRes(),
	})
}
