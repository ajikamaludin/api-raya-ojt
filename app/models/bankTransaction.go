package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankTransaction struct {
	ID               uuid.UUID `gorm:"primarykey;type:uuid"`
	BankAccountId    uuid.UUID `gorm:"not null;type:uuid"`
	BankId           uuid.UUID `gorm:"not null;type:uuid"`
	UserId           uuid.UUID `gorm:"not null;type:uuid"`
	Debit            float64   `gorm:"not null"`
	Credit           float64   `gorm:"not null"`
	Status           int       `gorm:"not null"`
	TransactionFeeId uuid.UUID `gorm:"type:uuid"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	CreatedBy        uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy        uuid.UUID      `gorm:"type:uuid"`
	DeletedBy        uuid.UUID      `gorm:"type:uuid"`
	ModCount         int
	// Relation BelongsTo
	Bankaccount    BankAccount      `gorm:"references:BankAccountId"`
	Bank           Bank             `gorm:"references:BankId"`
	User           User             `gorm:"references:UserId"`
	TransactionFee *BankTransaction `gorm:"references:TransactionFeeId"`
}

func (banktrx *BankTransaction) BeforeUpdate(tx *gorm.DB) (err error) {
	banktrx.ModCount = banktrx.ModCount + 1

	return
}
