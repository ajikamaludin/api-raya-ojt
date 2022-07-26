package repository

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/google/uuid"
)

func (r *Repository) GetAllAccountFavoriteUser(userId uuid.UUID, query string, favorites *[]models.BankAccountFavorite) error {
	db, _ := r.Gormdb.GetInstance()

	if query != "" {
		query = "%" + query + "%"
		err := db.Where("name ilike ? OR account_number ilike ?", query, query).
			Preload("Bankaccount").
			Preload("RayaAccount.UserAccount").
			Preload("Bank").
			Find(favorites, "user_id = ?", userId).Error
		return err
	}
	err := db.Preload("Bankaccount").
		Preload("RayaAccount.UserAccount").
		Preload("Bank").
		Find(favorites, "user_id = ?", userId).Error

	return err
}

func (r *Repository) GetAccountFavoriteUserByAccountNumber(
	userId uuid.UUID,
	accountNumber string,
	bank *models.Bank,
	favorite *models.BankAccountFavorite,
) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Where("user_id = ? AND account_number = ? AND bank_id = ?", userId, accountNumber, bank.ID).First(&favorite).Error

	return err
}

func (r *Repository) CreateAccountFavoriteUser(bankAccountFavorite *models.BankAccountFavorite) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Preload("Bankaccount").Preload("RayaAccount").Preload("Bank").Create(&bankAccountFavorite).Error

	return err
}

func (r *Repository) GetAccountFavoriteUserById(id string, bankAccountFavorite *models.BankAccountFavorite) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.First(&bankAccountFavorite, "id = ?", id).Error
	return err
}

func (r *Repository) DeleteAccountFavoriteUser(bankAccountFavorite *models.BankAccountFavorite) error {
	db, _ := r.Gormdb.GetInstance()

	err := db.Delete(&bankAccountFavorite).Error
	return err
}
