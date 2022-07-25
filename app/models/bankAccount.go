package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankAccount struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	BankId        uuid.UUID `gorm:"not null"`
	Name          string    `gorm:"not null"`
	AccountNumber string    `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid"`
	DeletedBy     uuid.UUID      `gorm:"type:uuid"`
	ModCount      int
	// Relation BelongsTo
	Bank Bank `gorm:"references:BankId"`
}

func (bankAcc *BankAccount) BeforeUpdate(tx *gorm.DB) (err error) {
	bankAcc.ModCount = bankAcc.ModCount + 1

	return
}
