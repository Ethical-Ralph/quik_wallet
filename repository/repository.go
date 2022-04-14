package repository

import "github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"

type Repository interface {
	CreateWallet(walletIdentifier string) error
	FindWallet(walletId string) (*models.Wallet, error)
	UpdateWalletBalance(walletId string, amount float64) error
}

func NewRepository(storageRepository Repository) Repository {
	return storageRepository
}
