package repository

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (r Repository) GetUserByEmail(email string, user *models.User) error {
	db, _ := r.Gormdb.GetInstance()
	err := db.Where("email = ?", email).Preload("UserAccount").First(user).Error
	return err
}

func (r Repository) ValidatePin(userId uuid.UUID, pin string) error {
	user := &models.User{}
	db, _ := r.Gormdb.GetInstance()

	err := db.Find(user, userId).Error
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(pin))

	return err
}

func (r Repository) CreateUser(user *models.User) error {
	db, err := r.Gormdb.GetInstance()
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UTC().UnixNano())
	user.UserAccount = &models.Account{
		AccountNumber: fmt.Sprintf("000%08d", rand.Intn(9999999999)),
		Balance:       0, // default balance must be 0
	}

	err = db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Seed() {}
