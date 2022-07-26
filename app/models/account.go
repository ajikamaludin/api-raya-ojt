package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	UserId        uuid.UUID `gorm:"not null;type:uuid"`
	AccountNumber string    `gorm:"not null"`
	Balance       float64   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid"`
	DeletedBy     uuid.UUID      `gorm:"type:uuid"`
	ModCount      int
	UserAccount   *User `gorm:"foreignKey:user_id"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.New()
	account.ModCount = 1

	return
}

func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	account.ModCount = account.ModCount + 1

	return
}
