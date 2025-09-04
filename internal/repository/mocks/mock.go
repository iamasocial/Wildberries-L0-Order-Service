package mock_repository

import (
	"L0/internal/entities"
	"context"
	"fmt"
)

type MockOrderRepository struct {
	Orders map[string]*entities.Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{Orders: make(map[string]*entities.Order)}
}

func (m *MockOrderRepository) GetOrderByUID(ctx context.Context, orderUID string) (*entities.Order, error) {
	order, ok := m.Orders[orderUID]
	if !ok {
		return nil, fmt.Errorf("order with UID %s not found", orderUID)
	}

	return order, nil
}

func (m *MockOrderRepository) SaveOrder(ctx context.Context, order *entities.Order) error {
	m.Orders[order.OrderUID] = order
	return nil
}

func (m *MockOrderRepository) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	orders := make([]entities.Order, 0, len(m.Orders))
	for _, order := range m.Orders {
		orders = append(orders, *order)
	}

	return orders, nil
}
