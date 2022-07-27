package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankAccountRes struct {
	ID            uuid.UUID
	BankID        uuid.UUID
	Name          string
	AccountNumber string
}

type BankAccount struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	BankID        uuid.UUID `gorm:"not null;type:uuid"`
	Name          string    `gorm:"not null"`
	AccountNumber string    `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	DeletedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	ModCount      int            `gorm:"default:0"`
	// Relation BelongsTo
	BankAccountBank Bank `gorm:"foreignKey:bank_id"`
}

func (bankAcc *BankAccount) BeforeCreate(tx *gorm.DB) (err error) {
	bankAcc.ID = uuid.New()
	bankAcc.ModCount = 1

	return
}

func (bankAcc *BankAccount) BeforeUpdate(tx *gorm.DB) (err error) {
	bankAcc.ModCount = bankAcc.ModCount + 1

	return
}

func (bankAcc *BankAccount) ToBankAccountRes() *BankAccountRes {
	return &BankAccountRes{
		ID:            bankAcc.ID,
		BankID:        bankAcc.BankID,
		Name:          bankAcc.Name,
		AccountNumber: bankAcc.AccountNumber,
	}
}
