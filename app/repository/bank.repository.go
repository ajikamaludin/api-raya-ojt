package repository

import (
	"errors"
	"time"

	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/helper"
)

func (r Repository) GetAllBanks(banks *[]models.Bank, query string) (err error) {
	db, _ := r.Gormdb.GetInstance()

	var expired time.Duration = 100
	key := "allbanks+" + query
	err = r.RedisClient.Get(key, &banks)
	if err != nil {
		if query != "" {
			query = "%" + query + "%"
			err = db.Where("name ILIKE ? OR code ILIKE ? OR short_name ILIKE ?", query, query, query).Find(&banks).Error
			r.RedisClient.Set(key, banks, expired*time.Second) // NOTE : i make it expired in 30 second for demo mode
			return
		}
		err = db.Find(&banks).Error
		r.RedisClient.Set(key, banks, expired*time.Second) // NOTE : i make it expired in 30 second for demo mode
	}

	return
}

func (r *Repository) CheckAccountNumber(bankId string, accNumber string) (*models.AccountNumberRes, error) {
	var bank models.Bank
	err := r.GetBankById(bankId, &bank)

	if err != nil {
		return nil, err
	}

	// Check Bank.Code is RAYA BANK
	if bank.Code == constants.RAYA_BANK_CODE {
		// search in account
		var account models.Account
		err := r.GetBankRayaAccount(accNumber, &account)

		if err != nil {
			return nil, err
		}

		return &models.AccountNumberRes{
			BankId:        bankId,
			AccountNumber: accNumber,
			Name:          account.UserAccount.Name,
			Bank:          bank.ToBankRes(),
		}, nil
	} else {
		// Search in bank_accounts table
		var bankAccount models.BankAccount
		err := r.GetBankAccountByAccountNumber(accNumber, bank, &bankAccount)
		if err != nil { // in bank_account table not found
			err = r.GetAccountNumberFromArtaJasa(accNumber, &bank, &bankAccount, true)
			if err != nil { // in arta jasa not found
				// TODO: if error timeout from artajasa, publish to pub/sub for proccess leter
				return nil, err
			}
		}

		return &models.AccountNumberRes{
			BankId:        bankId,
			AccountNumber: accNumber,
			Name:          bankAccount.Name,
			Bank:          bank.ToBankRes(),
		}, nil
	}
}

func (r *Repository) GetBankById(id string, bank *models.Bank) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.First(&bank, "id = ?", id).Error
	return err
}

func (r *Repository) GetBankRayaAccount(accountNumber string, account *models.Account) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Where("account_number = ?", accountNumber).Preload("UserAccount").First(&account).Error
	return err
}

func (r *Repository) GetBankAccountByAccountNumber(accountNumber string, bank models.Bank, bankAccount *models.BankAccount) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.First(&bankAccount, "account_number = ? AND bank_id = ?", accountNumber, bank.ID).Error
	return err
}

func (r *Repository) GetAccountNumberFromArtaJasa(accNumber string, bank *models.Bank, bankAccount *models.BankAccount, isSaveResult bool) error {
	// Mocking Api Call
	result := helper.CallArtaJasaApi(accNumber, bank, true)

	if result.Status != "success" {
		return errors.New("not found")
	}
	// save to db
	bankAccount = &models.BankAccount{
		BankID:        bank.ID,
		Name:          result.Data.HolderName,
		AccountNumber: result.Data.AccountNumber,
	}

	if isSaveResult {
		db, _ := r.Gormdb.GetInstance()
		db.Create(&bankAccount)
	}
	return nil
}
