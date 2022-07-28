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

	err = trxDb.Commit().Error
	if err != nil {
		trxDb.Rollback()
		return err
	}

	// Call Pubsub Here With Gorutine
	go r.GooglePubsub.Publish(constants.TRANSACTION_TOPIC_NAME, trx.ID.String())

	trx.Bank = &models.Bank{}
	db.First(&trx.Bank, "id = ?", trx.BankId)

	if trx.AccountId != (uuid.UUID{}) {
		trx.RayaAccount = &models.Account{}
		db.Preload("UserAccount").First(&trx.RayaAccount, "id = ?", trx.AccountId)
	}
	if trx.BankAccountId != (uuid.UUID{}) {
		trx.Bankaccount = &models.BankAccount{}
		db.First(&trx.Bankaccount, "id = ?", trx.BankAccountId)
	}

	if trx.TransactionFeeId != (uuid.UUID{}) { // uuid not nil
		trx.TransactionFee = &models.BankTransaction{}
		db.First(&trx.TransactionFee, "id = ?", trx.TransactionFeeId)
	}

	// TODO: publish transaction to PUB/SUB use goroutine
	return nil
}

func (r *Repository) GetTransactionById(id string, trx *models.BankTransaction) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Preload("TransactionFee").Preload("Bankaccount").
		Preload("RayaAccount.UserAccount").
		Preload("Bank").First(&trx, "id = ?", id).Error
	return err
}

func (r *Repository) GetLatestTransactions(userId uuid.UUID, transactions *[]models.BankTransaction, query string) error {
	db, _ := r.Gormdb.GetInstance()

	// TODO: leter on
	// if query != "" {
	// 	query = `%` + query + `%`
	// 	var bankaccount []models.BankAccount
	// 	db.Find(&bankaccount, "name ILIKE ? OR account_number ILIKE ?", query, query)
	// 	var account []models.Account
	// 	db.Joins("UserAccount").Find(&account, "Accounts.account_number ILIKE ? OR UserAccount.name ILIKE ? ", query, query)

	// 	fmt.Println(bankaccount, account)
	// }

	err := db.Order("created_at DESC").
		Where("debit = 0 AND (bank_account_id IS NOT NULL OR account_id IS NOT NULL)").
		Preload("Bank").
		Preload("Bankaccount").
		Preload("RayaAccount.UserAccount").
		Preload("TransactionFee").
		Find(&transactions, "user_id = ?", userId).Error

	return err
}
