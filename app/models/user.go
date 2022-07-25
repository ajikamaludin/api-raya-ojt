package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Pin       int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy uuid.UUID      `gorm:"type:uuid"`
	UpdatedBy uuid.UUID      `gorm:"type:uuid"`
	DeletedBy uuid.UUID      `gorm:"type:uuid"`
	ModCount  int
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()

	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.ModCount = user.ModCount + 1

	return
}
