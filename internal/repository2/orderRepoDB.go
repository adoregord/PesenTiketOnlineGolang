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
	ViewAllOrdersDB
	ViewUsersOrder
}
type CreateOrderDB interface {
	CreateOrderDB(order *domain.Order, kontek context.Context) (*domain.Order, error)
}
type ViewAllOrdersDB interface {
	ViewAllOrdersDB(kontek context.Context) (*[]domain.Order, error)
}
type ViewUsersOrder interface {
	ViewUsersOrder(UserID int, kontek context.Context) (*[]domain.Order, error)
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

// func to view all orders
func (repo OrderRepo) ViewAllOrdersDB(kontek context.Context) (*[]domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	query := `
	select 
		t.id, 
		t.order_date, 
		t.status, 
		t.payment_method, 
		t.total_price, 
		a."name",
		t.event_id,
		e."name",
		e."date",
		e."location",
		e.description,
		t2.id,
		t2.type,
		tt.quantity
   	from 
		"transaction" t
   	join 
		akun a ON t.user_id = a.id 
   	join 
		"event" e ON t.event_id = e.id 
   	join 
		transaction_ticket tt on t.id = tt.transaction_id
   	join 
		ticket t2 on tt.ticket_id = t2.id 
	`

	rows, err := repo.db.QueryContext(kontek, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[int]domain.Order)

	for rows.Next() {
		var (
			order  domain.Order
			ticket domain.TicketReq
		)

		err := rows.Scan(
			&order.ID,
			&order.OrderDate,
			&order.Status,
			&order.PaymentMethod,
			&order.TotalPrice,
			&order.User.Name,
			&order.Event.ID,
			&order.Event.Name,
			&order.Event.Date,
			&order.Event.Location,
			&order.Event.Description,
			&ticket.ID,
			&ticket.Type,
			&ticket.Quantity,
		)
		if err != nil {
			return nil, err
		}

		if data, exists := orders[order.ID]; exists {
			data.EventTicket = append(data.EventTicket, ticket)
			orders[order.ID] = data
		} else {
			order.EventTicket = []domain.TicketReq{ticket}
			orders[order.ID] = order
		}
	}

	orderList := make([]domain.Order, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, order)
	}

	return &orderList, nil
}

// func to view orders by user id
func (repo OrderRepo) ViewUsersOrder(userID int, kontek context.Context) (*[]domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	query := `
	select 
		t.id, 
		t.order_date, 
		t.status, 
		t.payment_method, 
		t.total_price, 
		a."name",
		t.event_id,
		e."name",
		e."date",
		e."location",
		e.description,
		t2.id,
		t2.type,
		tt.quantity
   	from 
		"transaction" t
   	join 
		akun a ON t.user_id = a.id 
   	join 
		"event" e ON t.event_id = e.id 
   	join 
		transaction_ticket tt on t.id = tt.transaction_id
   	join 
		ticket t2 on tt.ticket_id = t2.id
	where
		a.id = $1
	`

	rows, err := repo.db.QueryContext(kontek, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[int]domain.Order)

	for rows.Next() {
		var (
			order  domain.Order
			ticket domain.TicketReq
		)

		err := rows.Scan(
			&order.ID,
			&order.OrderDate,
			&order.Status,
			&order.PaymentMethod,
			&order.TotalPrice,
			&order.User.Name,
			&order.Event.ID,
			&order.Event.Name,
			&order.Event.Date,
			&order.Event.Location,
			&order.Event.Description,
			&ticket.ID,
			&ticket.Type,
			&ticket.Quantity,
		)
		if err != nil {
			return nil, err
		}

		if data, exists := orders[order.ID]; exists {
			data.EventTicket = append(data.EventTicket, ticket)
			orders[order.ID] = data
		} else {
			order.EventTicket = []domain.TicketReq{ticket}
			orders[order.ID] = order
		}
	}

	orderList := make([]domain.Order, 0, len(orders))
	for _, order := range orders {
		orderList = append(orderList, order)
	}

	return &orderList, nil
}
