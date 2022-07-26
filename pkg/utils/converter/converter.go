package converter

import "github.com/ajikamaludin/api-raya-ojt/app/models"

func MapBanksToBankRes(banks []models.Bank) []models.BankRes {
	var bankres []models.BankRes
	for _, v := range banks {
		bankres = append(bankres, v.ToBankRes())
	}
	return bankres
}

func MapBankAccountFavoriteToRes(bankfavs []models.BankAccountFavorite) []models.BankAccountFavoriteRes {
	var bankfavres []models.BankAccountFavoriteRes
	for _, v := range bankfavs {
		bankfavres = append(bankfavres, *v.ToBankAccountFavoriteRes())
	}
	return bankfavres
}
