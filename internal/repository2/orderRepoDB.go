package repository2

import (
	"context"
	"database/sql"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make Order db with map
type OrderRepo struct {
	db    *sql.DB
	mutek *sync.Mutex
}

func NewOrderRepo(db *sql.DB) OrderRepoInterface {
	return OrderRepo{
		db:    db,
		mutek: &sync.Mutex{},
	}
}

type OrderRepoInterface interface {
	CreateOrderDB
}
type CreateOrderDB interface {
	CreateOrderDB(order *domain.Order, kontek context.Context) (*domain.Order, error)
}

func (repo OrderRepo) CreateOrderDB(order *domain.Order, kontek context.Context) (*domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	query := `
	INSERT INTO transaction (status, payment_method, total_price, user_id, event_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, order_date
	`
	err := repo.db.QueryRowContext(kontek, query, order.Status, order.PaymentMethod, order.TotalPrice, order.User.ID, order.Event.ID).Scan(&order.ID, &order.OrderDate)
	if err != nil {
		return nil, err
	}

	query2 := `
	INSERT INTO transaction_ticket (transaction_id, ticket_id, quantity)
	VALUES ($1, $2, $3)
	`
	for _, ticket := range order.EventTicket {
		_, err := repo.db.ExecContext(kontek, query2, order.ID, ticket.ID, ticket.Quantity)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}
