package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	UserId        uuid.UUID `gorm:"not null"`
	AccountNumber string    `gorm:"not null"`
	Balance       float64   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid"`
	DeletedBy     uuid.UUID      `gorm:"type:uuid"`
	ModCount      int
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.New()

	return
}

func (acc *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	acc.ModCount = acc.ModCount + 1

	return
}
