package repository

import (
	"errors"

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

func (mrepo *MySqlRepository) FindWallet(walletId string) (*models.Wallet, error) {
	wallet := models.Wallet{}

	err := mrepo.conn.First(&wallet, models.Wallet{
		WalletIdentifier: walletId,
	}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &wallet, nil
}

func (mrepo *MySqlRepository) UpdateWalletBalance(walletId string, amount float64) error {
	if err := mrepo.conn.Model(&models.Wallet{}).Where("wallet_identifier = ?", walletId).Update("balance", amount).Error; err != nil {
		return err
	}

	return nil
}
