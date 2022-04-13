package service

import (
	"errors"

	"github.com/Ethical-Ralph/quik_wallet/infrastructure/cache"
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"github.com/Ethical-Ralph/quik_wallet/repository"
	"github.com/shopspring/decimal"
)

type Service struct {
	repository repository.Repository
	cache      cache.Cache
}

func NewService(repository repository.Repository, cache cache.Cache) *Service {
	return &Service{repository, cache}
}

func (s *Service) CreateWallet(wallet *models.Wallet) error {
	if wallet.WalletIdentifier == "" {
		return errors.New("Wallet Identifier can not be empty")
	}

	err := s.repository.CreateWallet(wallet.WalletIdentifier)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetWalletBalance(walletId string) (interface{}, error) {
	wallet, err := s.repository.FindWallet(walletId)
	if err != nil {
		return nil, err
	}

	return struct {
		Balance float64 `json:"balance"`
	}{
		Balance: wallet.Balance,
	}, nil
}

func (s *Service) CreditWallet(walletId string, amount float64) error {
	zero := decimal.NewFromInt(0)
	amountToAdd := decimal.NewFromFloat(amount)

	if amountToAdd.LessThanOrEqual(zero) {
		return errors.New("Invalid Amount")
	}

	wallet, err := s.repository.FindWallet(walletId)
	if err != nil {
		return err
	}

	walletBalance := decimal.NewFromFloat(wallet.Balance)

	totalAmount := walletBalance.Add(amountToAdd)
	totalAmountTofloat, _ := totalAmount.Float64()

	err = s.repository.UpdateWalletBalance(walletId, totalAmountTofloat)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DebitWallet(walletId string, amount float64) error {
	zero := decimal.NewFromInt(0)
	amountToSub := decimal.NewFromFloat(amount)

	if amountToSub.LessThanOrEqual(zero) {
		return errors.New("Invalid Amount")
	}

	wallet, err := s.repository.FindWallet(walletId)
	if err != nil {
		return err
	}

	walletBalance := decimal.NewFromFloat(wallet.Balance)
	totalAmount := walletBalance.Sub(amountToSub)

	if totalAmount.LessThan(zero) {
		return errors.New("Insufficent Balance")
	}

	totalAmountTofloat, _ := totalAmount.Float64()

	err = s.repository.UpdateWalletBalance(walletId, totalAmountTofloat)
	if err != nil {
		return err
	}

	return nil
}
