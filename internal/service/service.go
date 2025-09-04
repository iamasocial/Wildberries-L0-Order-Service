package service

import (
	"L0/internal/cache"
	"L0/internal/repository"
	"L0/pkg/logger"
)

type Service struct {
	Order Order
}

func NewService(repo *repository.Repository, cache cache.Cache, logger logger.Logger) *Service {
	return &Service{Order: NewOrder(repo.Order, cache, logger)}
}
