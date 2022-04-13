package repository

type Repository interface {
	CreateWallet(walletIdentifier string) error
	// DebitWallet()
	// CreditWallet()
}

func NewRepository(storageRepository Repository) Repository {
	return storageRepository
}
