package repository

import (
	"time"

	"github.com/ajikamaludin/api-raya-ojt/app/models"
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

func (r *Repository) ValidateAccountNumber() error {
	// check bank id is raya
	// search in account table , ok return , not found return
	// not raya account
	// mock to function
	return nil
}
