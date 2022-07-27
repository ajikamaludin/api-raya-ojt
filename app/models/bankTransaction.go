package models

import (
	"time"

	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateBankTransactionReq struct {
	BankId        string  `validate:"required,uuid"`
	AccountNumber string  `validate:"required,min=6,numeric"`
	Amount        float64 `validate:"required,min=10000,numeric"`
}

type BankTransactionRes struct {
	ID             uuid.UUID
	BankID         uuid.UUID
	UserId         uuid.UUID
	Debit          float64
	Credit         float64
	Status         int
	StatusText     string
	TransactionFee float64
	CreateAt       time.Time
	CreatedBy      uuid.UUID
	UpdatedAt      time.Time
	UpdatedBy      uuid.UUID
	Bank           BankRes
	BankAccount    BankAccountRes
}

type BankTransaction struct {
	ID               uuid.UUID `gorm:"primarykey;type:uuid"`
	BankAccountId    uuid.UUID `gorm:"type:uuid;default:null"`
	AccountId        uuid.UUID `gorm:"type:uuid;default:null"`
	BankId           uuid.UUID `gorm:"type:uuid;default:null"`
	UserId           uuid.UUID `gorm:"not null;type:uuid"`
	Debit            float64   `gorm:"not null;default:0"`
	Credit           float64   `gorm:"not null;default:0"`
	Status           int       `gorm:"not null;default:0"` // 0 is pending, 1 success, 2 fail
	TransactionFeeId uuid.UUID `gorm:"type:uuid;default:null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	CreatedBy        uuid.UUID      `gorm:"type:uuid;default:null"`
	UpdatedBy        uuid.UUID      `gorm:"type:uuid;default:null"`
	DeletedBy        uuid.UUID      `gorm:"type:uuid;default:null"`
	ModCount         int            `gorm:"default:0"`
	// Relation BelongsTo
	RayaAccount    *Account         `gorm:"foreignKey:AccountId"`
	Bankaccount    *BankAccount     `gorm:"foreignKey:BankAccountId"`
	Bank           *Bank            `gorm:"foreignKey:BankId"`
	User           User             `gorm:"foreignKey:UserId"`
	TransactionFee *BankTransaction `gorm:"foreignKey:TransactionFeeId"`
}

func (banktrx *BankTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	banktrx.ID = uuid.New()

	return
}

func (banktrx *BankTransaction) BeforeUpdate(tx *gorm.DB) (err error) {
	banktrx.ModCount = banktrx.ModCount + 1

	return
}

func (banktrx *BankTransaction) ToBankTransactionRes() BankTransactionRes {
	bankTrxRes := BankTransactionRes{
		ID:             banktrx.ID,
		BankID:         banktrx.BankId,
		UserId:         banktrx.UserId,
		Debit:          banktrx.Debit,
		Credit:         banktrx.Credit,
		Status:         banktrx.Status,
		TransactionFee: 0,
		CreateAt:       banktrx.CreatedAt,
		CreatedBy:      banktrx.CreatedBy,
		UpdatedAt:      banktrx.UpdatedAt,
		UpdatedBy:      banktrx.UpdatedBy,
	}

	if banktrx.Status == constants.TRX_PENDING {
		bankTrxRes.StatusText = "pending"
	} else if banktrx.Status == constants.TRX_SUCCESS {
		bankTrxRes.StatusText = "success"
	} else if banktrx.Status == constants.TRX_FAIL {
		bankTrxRes.StatusText = "fail"
	}

	if banktrx.TransactionFee != nil {
		bankTrxRes.TransactionFee = banktrx.TransactionFee.Credit
	}

	if banktrx.Bank != nil {
		bankTrxRes.Bank = banktrx.Bank.ToBankRes()
	}

	if banktrx.Bankaccount != nil {
		bankTrxRes.BankAccount = *banktrx.Bankaccount.ToBankAccountRes()
	}

	if banktrx.RayaAccount != nil && banktrx.RayaAccount.UserAccount != nil {
		bankTrxRes.BankAccount = BankAccountRes{
			ID:            banktrx.RayaAccount.ID,
			BankID:        banktrx.BankId,
			Name:          banktrx.RayaAccount.UserAccount.Name,
			AccountNumber: banktrx.RayaAccount.AccountNumber,
		}
	}

	return bankTrxRes
}
