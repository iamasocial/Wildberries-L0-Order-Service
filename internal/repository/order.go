package repository

import (
	"L0/internal/entities"
	"context"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*entities.Order, error)
	SaveOrder(ctx context.Context, order *entities.Order) error
	GetAllOrders(ctx context.Context) ([]entities.Order, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetOrderByUID(ctx context.Context, orderUID string) (*entities.Order, error) {
	var order entities.Order

	query := "SELECT " + orderColumns + " FROM orders WHERE order_uid=$1"
	err := r.db.Get(&order, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + deliveryColumns + " FROM delivery WHERE order_uid=$1"
	err = r.db.Get(&order.Delivery, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + paymentColumns + " FROM payments WHERE order_uid=$1"
	err = r.db.Get(&order.Payment, query, orderUID)
	if err != nil {
		return nil, err
	}

	query = "SELECT " + itemColumns + " FROM items WHERE order_uid=$1"
	err = r.db.Select(&order.Items, query, orderUID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) SaveOrder(ctx context.Context, order *entities.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
		}

		err = tx.Commit()
	}()

	orderQuery := `INSERT INTO orders (
	` + orderColumns + `
	) VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature,
	:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)`

	_, err = tx.NamedExecContext(ctx, orderQuery, order)
	if err != nil {
		return err
	}

	d := deliveryRow{
		OrderUID: order.OrderUID,
		Name:     order.Delivery.Name,
		Phone:    order.Delivery.Phone,
		Zip:      order.Delivery.Zip,
		City:     order.Delivery.City,
		Address:  order.Delivery.Address,
		Region:   order.Delivery.Region,
		Email:    order.Delivery.Email,
	}

	deliveryQuery := `
		INSERT INTO delivery (` + deliveryColumnsWithOrderUID + `
		) VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)
		`

	_, err = tx.NamedExecContext(ctx, deliveryQuery, d)
	if err != nil {
		return err
	}

	p := paymentRow{
		OrderUID:     order.OrderUID,
		Transaction:  order.Payment.Transaction,
		RequestID:    order.Payment.RequestID,
		Currency:     order.Payment.Currency,
		Provider:     order.Payment.Provider,
		Amount:       order.Payment.Amount,
		PaymentDt:    order.Payment.PaymentDt,
		Bank:         order.Payment.Bank,
		DeliveryCost: order.Payment.DeliveryCost,
		GoodsTotal:   order.Payment.GoodsTotal,
		CustomFee:    order.Payment.CustomFee,
	}

	paymentQuery := `
		INSERT INTO payments (` + paymentColumnsWithOrderUID + `
		) VALUES (:order_uid, :transaction, :request_id, :currency, :provider, :amount,
		:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
	`

	_, err = tx.NamedExecContext(ctx, paymentQuery, p)
	if err != nil {
		return err
	}

	itemQuery := `
		INSERT INTO items (` + itemColumnsWithOrderUID + `
		) VALUES (:order_uid, :chrt_id, :track_number, :price, :rid, :name, :sale,
		:size, :total_price, :nm_id, :brand, :status)
		`

	for _, it := range order.Items {
		i := itemRow{
			OrderUID:    order.OrderUID,
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			Rid:         it.Rid,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		}

		_, err = tx.NamedExecContext(ctx, itemQuery, i)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order

	query := "SELECT " + orderColumns + " FROM orders"
	err := r.db.Select(&orders, query)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		query = "SELECT " + deliveryColumns + " FROM delivery WHERE order_uid=$1"
		err = r.db.Get(&order.Delivery, query, order.OrderUID)
		if err != nil {
			return nil, err
		}

		query = "SELECT " + paymentColumns + " FROM payments WHERE order_uid=$1"
		err = r.db.Get(&order.Payment, query, order.OrderUID)
		if err != nil {
			return nil, err
		}

		query = "SELECT " + itemColumns + " FROM items WHERE order_uid=$1"
		err = r.db.Select(&order.Items, query, order.OrderUID)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}
