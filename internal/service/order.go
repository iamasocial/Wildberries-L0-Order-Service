package service

import (
	"L0/internal/entities"
	"L0/internal/repository"
)

type Order interface {
	GetOrderByUID(orderUID string) (*entities.Order, error)
}

type orderService struct {
	repo repository.Order
}

func NewOrder(repo repository.Order) *orderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetOrderByUID(orderUID string) (*entities.Order, error) {
	return s.repo.GetOrderByUID(orderUID)
}
