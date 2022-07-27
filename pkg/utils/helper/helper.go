package helper

import (
	"encoding/json"
	"errors"
	. "net/http"
	"time"

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

type ApiReposnse struct {
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"accountNumber"`
	ID            string    `json:"id"`
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

func CallArtaJasaApi(accNumber string, bank *models.Bank) (*ArtaJasaAccountNumberResponse, error) {
	// TODO: call api here and get response
	const ENDPOINT = "https://62dfbf45976ae7460bf2b9e6.mockapi.io/api/v1/AccountNumber?accountNumber="

	request, err := NewRequest("GET", ENDPOINT+accNumber, nil)
	if err != nil {
		return nil, err
	}

	response, err := DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var apiReponses []ApiReposnse
	err = json.NewDecoder(response.Body).Decode(&apiReponses)
	if err != nil {
		return nil, err
	}

	if len(apiReponses) > 0 {
		// NOTE : im mocking the api result
		return &ArtaJasaAccountNumberResponse{
			Status:  "success",
			Message: "account number found",
			Data: struct {
				BankName      string
				AccountNumber string
				HolderName    string
			}{
				BankName:      bank.Name,
				AccountNumber: apiReponses[0].AccountNumber,
				HolderName:    apiReponses[0].Name,
			},
		}, nil
	}

	return nil, errors.New("not found")
}
