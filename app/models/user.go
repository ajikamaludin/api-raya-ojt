package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserLoginReq struct {
	Email    string `validate:"required,min=3,email"`
	Password string `validate:"required,min=6"`
}

type PinReq struct {
	Pin string `validate:"required,min=6"`
}

type UserRegisterReq struct {
	Name     string `validate:"required,min=3"`
	Email    string `validate:"required,min=3,email"`
	Password string `validate:"required,min=6"`
	Pin      string `validate:"required,min=6"`
}

type UserRes struct {
	ID            uuid.UUID
	AccountNumber string
	Balance       float64
	Email         string
	CreatedAt     time.Time
}

type User struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Pin       string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy uuid.UUID      `gorm:"type:uuid"`
	DeletedBy uuid.UUID      `gorm:"type:uuid"`
	ModCount  int
	// Relation
	// BankTransactions []BankTransaction
	UserAccount *Account
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()

	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.ModCount = user.ModCount + 1

	return
}

func (user User) ToUserRes() *UserRes {
	return &UserRes{
		ID:            user.ID,
		AccountNumber: user.UserAccount.AccountNumber,
		Balance:       user.UserAccount.Balance,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
	}
}
