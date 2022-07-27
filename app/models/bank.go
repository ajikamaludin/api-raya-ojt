package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountNumberReq struct {
	BankId        string `validate:"required,uuid"`
	AccountNumber string `validate:"required,min=6,numeric"`
}

type AccountNumberRes struct {
	BankId        string
	AccountNumber string
	Name          string
	Bank          BankRes
}

type BankRes struct {
	ID             uuid.UUID
	Name           string
	ShortName      string
	LogoUrl        string
	Code           string
	TransactionFee float64
}

type Bank struct {
	ID             uuid.UUID `gorm:"primarykey;type:uuid"`
	Name           string    `gorm:"not null"`
	ShortName      string    `gorm:"not null"`
	LogoUrl        string    `gorm:"not null"`
	Code           string    `gorm:"not null"`
	TransactionFee float64   `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;default:null"`
	UpdatedBy      uuid.UUID      `gorm:"type:uuid;default:null"`
	DeletedBy      uuid.UUID      `gorm:"type:uuid;default:null"`
	ModCount       int            `gorm:"default:0"`
}

func (bank *Bank) BeforeCreate(tx *gorm.DB) (err error) {
	bank.ID = uuid.New()

	return
}

func (bank *Bank) BeforeUpdate(tx *gorm.DB) (err error) {
	bank.ModCount = bank.ModCount + 1

	return
}

func (bank Bank) ToBankRes() BankRes {
	return BankRes{
		ID:             bank.ID,
		Name:           bank.Name,
		ShortName:      bank.ShortName,
		LogoUrl:        bank.LogoUrl,
		Code:           bank.Code,
		TransactionFee: bank.TransactionFee,
	}
}
