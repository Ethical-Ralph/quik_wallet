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

	walletExist, err := s.repository.FindWallet(wallet.WalletIdentifier)
	if err != nil {
		return errors.New("An internal error ocurred")
	}

	if walletExist != nil && walletExist.WalletIdentifier == wallet.WalletIdentifier {
		return errors.New("Wallet with this identifier already exist")
	}

	err = s.repository.CreateWallet(wallet.WalletIdentifier)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetWalletBalance(walletId string) (interface{}, error) {
	fromCache := s.cache.Get(walletId)
	var amount float64

	getFromDb := func() (*float64, error) {
		wallet, err := s.repository.FindWallet(walletId)

		if wallet == nil {
			return nil, errors.New("Wallet doesn't exist")
		}

		if err != nil {
			return nil, err
		}
		s.cache.Set(walletId, wallet.Balance)
		return &wallet.Balance, nil
	}

	if fromCache != "" {
		decimalAmount, err := decimal.NewFromString(fromCache)
		if err != nil {
			balance, err := getFromDb()
			if err != nil {
				return nil, err
			}

			amount = *balance
		} else {
			amount, _ = decimalAmount.Float64()
		}

	} else {
		balance, err := getFromDb()
		if err != nil {
			return nil, err
		}

		amount = *balance
	}

	return struct {
		Balance float64 `json:"balance"`
	}{
		Balance: amount,
	}, nil
}

func (s *Service) CreditWallet(walletId string, amount float64) (*float64, error) {
	zero := decimal.NewFromInt(0)
	amountToAdd := decimal.NewFromFloat(amount)

	if amountToAdd.LessThanOrEqual(zero) {
		return nil, errors.New("Invalid Amount")
	}

	wallet, err := s.repository.FindWallet(walletId)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New("Wallet doesn't exist")
	}

	walletBalance := decimal.NewFromFloat(wallet.Balance)

	totalAmount := walletBalance.Add(amountToAdd)
	totalAmountTofloat, _ := totalAmount.Float64()

	err = s.repository.UpdateWalletBalance(walletId, totalAmountTofloat)
	if err != nil {
		return nil, err
	}
	s.cache.Del(walletId)

	return &totalAmountTofloat, nil
}

func (s *Service) DebitWallet(walletId string, amount float64) (*float64, error) {
	zero := decimal.NewFromInt(0)
	amountToSub := decimal.NewFromFloat(amount)

	if amountToSub.LessThanOrEqual(zero) {
		return nil, errors.New("Invalid Amount")
	}

	wallet, err := s.repository.FindWallet(walletId)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New("Wallet doesn't exist")
	}

	walletBalance := decimal.NewFromFloat(wallet.Balance)
	totalAmount := walletBalance.Sub(amountToSub)

	if totalAmount.LessThan(zero) {
		return nil, errors.New("Insufficent Balance")
	}

	totalAmountTofloat, _ := totalAmount.Float64()

	err = s.repository.UpdateWalletBalance(walletId, totalAmountTofloat)
	if err != nil {
		return nil, err
	}
	s.cache.Del(walletId)

	return &totalAmountTofloat, nil
}
