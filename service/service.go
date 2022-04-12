package service

import (
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/cache"
	"github.com/Ethical-Ralph/quik_wallet/repository"
)

type Service struct {
	repository *repository.Repository
	cache      *cache.Cache
}

func NewService(repository *repository.Repository, cache *cache.Cache) *Service {
	return &Service{repository, cache}
}
