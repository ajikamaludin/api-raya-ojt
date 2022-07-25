package repository

import (
	"github.com/ajikamaludin/api-raya-ojt/app/models"
)

func (r Repository) GetAllBanks(banks *[]models.Bank, query string) (err error) {
	db, _ := r.Gormdb.GetInstance()

	if query != "" {
		query = "%" + query + "%"
		err = db.Where("name ILIKE ? OR code ILIKE ? OR short_name ILIKE ?", query, query, query).Find(&banks).Error
		return
	}

	err = db.Find(&banks).Error
	return
}

func (r *Repository) ValidateAccountNumber() error {
	// check bank id is raya
	// search in account table , ok return , not found return
	// not raya account
	// mock to function
	return nil
}
