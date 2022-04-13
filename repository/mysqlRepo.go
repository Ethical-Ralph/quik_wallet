package repository

import (
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"gorm.io/gorm"
)

type MySqlRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) *MySqlRepository {
	return &MySqlRepository{conn}
}

func (mrepo *MySqlRepository) CreateWallet(walletIdentifier string) error {
	return mrepo.conn.Create(&models.Wallet{
		WalletIdentifier: walletIdentifier,
		Balance:          0,
	}).Error
}
