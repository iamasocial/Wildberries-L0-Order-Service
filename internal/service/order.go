package service

import (
	"L0/internal/cache"
	"L0/internal/entities"
	"L0/internal/repository"
	"L0/pkg/logger"
	"context"
)

type Order interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*entities.Order, error)
	SaveOrder(ctx context.Context, order *entities.Order) error
	LoadCache(ctx context.Context) error
}

type orderService struct {
	repo   repository.Order
	cache  cache.Cache
	logger logger.Logger
}

func NewOrder(repo repository.Order, cache cache.Cache, logger logger.Logger) *orderService {
	return &orderService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *orderService) GetOrderByUID(ctx context.Context, orderUID string) (*entities.Order, error) {
	if order, ok := s.cache.Get(orderUID); ok {
		s.logger.Debug("order retrieved from cache", "order_uid", orderUID)
		return order, nil
	}
	s.logger.Debug("order not found in cache", "order_uid", orderUID)

	order, err := s.repo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		s.logger.Error("failed to get order from DB", "order_uid", orderUID, "error", err.Error())
		return nil, err
	}

	s.cache.Set(order)
	s.logger.Debug("order cached", "order_uid", orderUID)
	s.logger.Info("order retrieved from DB", "order_uid", orderUID)
	return order, nil
}

func (s *orderService) SaveOrder(ctx context.Context, order *entities.Order) error {
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		s.logger.Error("failed to save order to DB", "order_uid", order.OrderUID, "error", err.Error())
		return err
	}

	s.cache.Set(order)
	s.logger.Info("order saved and cached", "order_uid", order.OrderUID)

	return nil
}

func (s *orderService) LoadCache(ctx context.Context) error {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		s.logger.Error("failed to get all orders from DB", "error", err.Error())
		return err
	}

	for _, order := range orders {
		s.cache.Set(&order)
		s.logger.Debug("order cached", "order_uid", order.OrderUID)
	}

	s.logger.Info("cache loaded", "orders_count", len(orders))

	return nil
}
