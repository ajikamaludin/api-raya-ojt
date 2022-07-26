package repository

import (
	"errors"

	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/google/uuid"
)

func (r *Repository) GetUserBalance(userId uuid.UUID, account *models.Account) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.First(&account, "user_id = ?", userId).Error

	return err
}

func (r *Repository) CreateBankTransaction(trx *models.BankTransaction, bank *models.Bank) error {
	db, _ := r.Gormdb.GetInstance()
	credit := 0.0

	// get use account balance
	var account models.Account
	err := db.First(&account, "user_id = ?", trx.UserId).Error
	if err != nil {
		return err
	}

	trxDb := db.Begin()
	// record transaction free
	if bank.TransactionFee > 0 {
		credit += bank.TransactionFee
		transaction := &models.BankTransaction{
			UserId:    trx.UserId,
			Debit:     0,
			Credit:    bank.TransactionFee,
			Status:    constants.TRX_SUCCESS,
			CreatedBy: trx.UserId, // note: use gorm hook
		}
		err := trxDb.Create(&transaction).Error
		if err != nil {
			trxDb.Rollback()
			return err
		}

		trx.TransactionFeeId = transaction.ID
	}

	// validate balance
	credit += trx.Credit
	if account.Balance < credit {
		trxDb.Rollback()
		return errors.New("balance is not enough")
	}

	// update use balance
	account.Balance = account.Balance - credit
	account.UpdatedBy = trx.UserId
	trxDb.Save(&account)

	// record transaction
	err = trxDb.Create(&trx).Error
	if err != nil {
		trxDb.Rollback()
		return err
	}

	trxDb.Commit()

	trx.Bank = &models.Bank{}
	db.First(&trx.Bank, "id = ?", trx.BankId)

	trx.Bankaccount = &models.BankAccount{}
	db.First(&trx.Bankaccount, "id = ?", trx.BankAccountId)

	if trx.TransactionFeeId != (uuid.UUID{}) { // uuid not nil
		trx.TransactionFee = &models.BankTransaction{}
		db.First(&trx.TransactionFee, "id = ?", trx.TransactionFeeId)
	}

	// TODO: publish transaction to PUB/SUB
	return nil
}

func (r *Repository) GetTransactionById(id string, trx *models.BankTransaction) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Preload("TransactionFee").Preload("Bankaccount").
		Preload("Bank").First(&trx, "id = ?", id).Error
	return err
}

func (r *Repository) GetLatestTransactions(userId uuid.UUID, transactions *[]models.BankTransaction, query string) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Order("created_at DESC").
		Where("bank_account_id is not null").
		Preload("Bank").
		Preload("Bankaccount").
		Preload("TransactionFee").
		Find(&transactions, "user_id = ?", userId).Error

	return err
}