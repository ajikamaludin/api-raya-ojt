package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankAccountFavoriteReq struct {
	BankID        string `validate:"required,uuid"`
	Name          string `validate:"required"`
	AccountNumber string `validate:"required,min=6,numeric"`
}

type BankAccountFavoriteRes struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	BankID        uuid.UUID
	Name          string
	AccountNumber string
	Bankaccount   BankAccountRes
}

type BankAccountFavorite struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid"`
	UserId        uuid.UUID `gorm:"not null;type:uuid"`
	BankId        uuid.UUID `gorm:"not null;type:uuid"`
	BankAccountId uuid.UUID `gorm:"type:uuid;default:null"`
	AccountId     uuid.UUID `gorm:"type:uuid;default:null"`
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
	Bankaccount *BankAccount `gorm:"foreignKey:BankAccountId"`
	RayaAccount *Account     `gorm:"foreignKey:AccountId"`
	Bank        *Bank        `gorm:"foreignKey:BankId"`
}

func (bankAccFav *BankAccountFavorite) BeforeCreate(tx *gorm.DB) (err error) {
	bankAccFav.ID = uuid.New()

	return
}

func (bankAccFav *BankAccountFavorite) BeforeUpdate(tx *gorm.DB) (err error) {
	bankAccFav.ModCount = bankAccFav.ModCount + 1

	return
}

func (bankAccFav *BankAccountFavorite) ToBankAccountFavoriteRes() *BankAccountFavoriteRes {
	bankAccountFavoriteRes := &BankAccountFavoriteRes{
		ID:            bankAccFav.ID,
		BankID:        bankAccFav.BankId,
		UserID:        bankAccFav.UserId,
		Name:          bankAccFav.Name,
		AccountNumber: bankAccFav.AccountNumber,
	}

	if bankAccFav.Bankaccount != nil {
		bankAccountFavoriteRes.Bankaccount = BankAccountRes{
			ID:            bankAccFav.Bankaccount.ID,
			BankID:        bankAccFav.Bankaccount.BankID,
			Name:          bankAccFav.Bankaccount.Name,
			AccountNumber: bankAccFav.Bankaccount.AccountNumber,
		}
	}

	if bankAccFav.RayaAccount != nil {
		bankAccountFavoriteRes.Bankaccount = BankAccountRes{
			ID:            bankAccFav.RayaAccount.ID,
			BankID:        bankAccFav.BankId,
			AccountNumber: bankAccFav.RayaAccount.AccountNumber,
		}

		if bankAccFav.RayaAccount.UserAccount != nil {
			bankAccountFavoriteRes.Bankaccount.Name = bankAccFav.RayaAccount.UserAccount.Name
		}
	}

	return bankAccountFavoriteRes
}
