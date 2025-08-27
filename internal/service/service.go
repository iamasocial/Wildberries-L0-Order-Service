package service

import "L0/internal/repository"

type Service struct {
	Order Order
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Order: NewOrder(repo.Order)}
}
