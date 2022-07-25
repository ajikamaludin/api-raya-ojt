package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	CreatedBy      uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy      uuid.UUID      `gorm:"type:uuid"`
	DeletedBy      uuid.UUID      `gorm:"type:uuid"`
	ModCount       int
}

func (bank *Bank) BeforeUpdate(tx *gorm.DB) (err error) {
	bank.ModCount = bank.ModCount + 1

	return
}
