package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	Order Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Order: NewOrderRepository(db)}
}
