package helper

import (
	"fmt"

	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"gorm.io/gorm"
)

type ArtaJasaAccountNumberResponse struct {
	Status  string
	Message string
	Data    struct {
		BankName      string
		AccountNumber string
		HolderName    string
	}
}

func SeedBank() *[]models.Bank {
	banks := []models.Bank{
		{
			Name:           "Mandiri",
			ShortName:      "MDR",
			LogoUrl:        "https://bankmandiri.co.id/image/layout_set_logo?img_id=31567&t=1658510036346",
			Code:           "008",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Rakyat Indonesia",
			ShortName:      "BRI",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "002",
			TransactionFee: 0,
		},
		{
			Name:           "Bank Raya Indonesia",
			ShortName:      "RAYA",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "002",
			TransactionFee: 0,
		},
		{
			Name:           "Bank Tabungan Pensiunan Nasional - Jenius",
			ShortName:      "BTPN",
			LogoUrl:        "https://www.btpn.com/website/static/img/logo.png",
			Code:           "213",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Central Asia",
			ShortName:      "BCA",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "014",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Negara Indonesia",
			ShortName:      "BNI",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "009",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank CIMB Niaga",
			ShortName:      "CIMB",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "022",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Muamalat",
			ShortName:      "Muamalat",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "147",
			TransactionFee: 2500,
		},
		{
			Name:           "Permata Bank",
			ShortName:      "Permata",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "013",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Danamon",
			ShortName:      "Danamon",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "011",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Mega",
			ShortName:      "Mega",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "426",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Sinarmas",
			ShortName:      "Sinarmas",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "153",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank OCBC NISP",
			ShortName:      "OCBC",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "028",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Bukopin",
			ShortName:      "BBUK",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "441",
			TransactionFee: 2500,
		},
		{
			Name:           "Citibank",
			ShortName:      "CITI",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "031",
			TransactionFee: 2500,
		},
		{
			Name:           "Bank Jateng",
			ShortName:      "BJT",
			LogoUrl:        "https://bankraya.co.id/img/logo.png",
			Code:           "113",
			TransactionFee: 2500,
		},
	}

	return &banks
}

func Seed(db *gorm.DB) {
	tx := db.Begin()
	var banks *[]models.Bank

	db.Find(&banks)

	if len(*banks) <= 0 {
		banks = SeedBank()
		tx.Create(banks)
	}

	tx.Commit()
}

func CallArtaJasaApi(accNumber string, bank *models.Bank, isOk bool) *ArtaJasaAccountNumberResponse {
	// TODO: call api here and get response

	// NOTE : im mocking the api result
	if isOk {
		return &ArtaJasaAccountNumberResponse{
			Status:  "success",
			Message: "account number found",
			Data: struct {
				BankName      string
				AccountNumber string
				HolderName    string
			}{
				BankName:      bank.Name,
				AccountNumber: accNumber,
				HolderName:    "Budi " + bank.Name,
			},
		}
	}

	return &ArtaJasaAccountNumberResponse{
		Status:  "fail",
		Message: "account number not found",
	}
}
