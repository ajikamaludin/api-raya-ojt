package converter

import "github.com/ajikamaludin/api-raya-ojt/app/models"

func MapBanksToBankRes(banks []models.Bank) []models.BankRes {
	var bankres []models.BankRes
	for _, v := range banks {
		bankres = append(bankres, v.ToBankRes())
	}
	return bankres
}
