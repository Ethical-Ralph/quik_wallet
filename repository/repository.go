package repository

type Repository interface {
	CreateWallet()
	// DebitWallet()
	// CreditWallet()
}

func NewRepository(storageRepository Repository) Repository {
	return storageRepository
}
