package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	UserId        uuid.UUID `gorm:"unique;not null;type:uuid"`
	AccountNumber string    `gorm:"unique;not null"`
	Balance       float64   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	DeletedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	ModCount      int            `gorm:"default:0"`
	UserAccount   *User          `gorm:"foreignKey:user_id"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.New()

	return
}

func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	account.ModCount = account.ModCount + 1

	return
}
