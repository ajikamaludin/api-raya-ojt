package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankAccountFavorite struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	BankAccountId uuid.UUID `gorm:"not null"`
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
	Bankaccount BankAccount `gorm:"references:BankAccountId"`
}

func (bankAccFav *BankAccountFavorite) BeforeUpdate(tx *gorm.DB) (err error) {
	bankAccFav.ModCount = bankAccFav.ModCount + 1

	return
}
