package repository

import (
	"L0/internal/entities"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	GetOrderByUID(orderUID string) (*entities.Order, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetOrderByUID(orderUID string) (*entities.Order, error) {
	var order entities.Order

	query := "SELECT " + strings.Join(orderColumns, ", ") + " FROM orders WHERE order_uid=$1"
	err := r.db.Get(&order, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + strings.Join(deliveryColumns, ", ") + " FROM delivery WHERE order_uid=$!"
	err = r.db.Get(&order.Delivery, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + strings.Join(paymentColumns, ", ") + " FROM payment WHERE order_uid=$1"
	err = r.db.Get(&order.Payment, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + strings.Join(itemColumns, ", ") + " FROM items WHERE order_uid=$1"
	err = r.db.Select(&order.Items, query, orderUID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
