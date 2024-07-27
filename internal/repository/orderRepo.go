package repository

import (
	"context"
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make Order db with map
type OrderRepo struct {
	Orders map[int]domain.Order
	mutek  *sync.Mutex
}

func NewOrderRepo() OrderRepoInterface {
	return OrderRepo{
		Orders: map[int]domain.Order{},
		mutek:  &sync.Mutex{},
	}
}

type OrderRepoInterface interface {
	CreateOrder
	GetOrderByID
	GetAllOrders
}
type CreateOrder interface {
	CreateOrder(order *domain.Order, kontek context.Context) (*domain.Order, error)
}
type GetOrderByID interface {
	GetOrderByID(userID int, kontek context.Context) ([]domain.Order, error)
}
type GetAllOrders interface {
	GetAllOrders(kontek context.Context) ([]domain.Order, error)
}

func (repo OrderRepo) CreateOrder(order *domain.Order, kontek context.Context) (*domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		if repo.Orders == nil || len(repo.Orders) == 0 {
			order.ID = 1
		} else {
			order.ID = repo.Orders[len(repo.Orders)].ID + 1
		}

		repo.Orders[order.ID] = *order
		return order, nil
	}
}

// func to get All order by User ID
func (repo OrderRepo) GetOrderByID(userID int, kontek context.Context) ([]domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		var ordersUser []domain.Order
		for _, Order := range repo.Orders {
			if Order.User.ID == userID {
				ordersUser = append(ordersUser, Order)
			}
		}
		if len(ordersUser) == 0 {
			return nil, errors.New("THIS USER HAVENT BUY A TICKET")
		}
		return ordersUser, nil
	}
}

func (repo OrderRepo) GetAllOrders(kontek context.Context) ([]domain.Order, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		Orders := make([]domain.Order, 0, len(repo.Orders))
		for _, Order := range repo.Orders {
			Orders = append(Orders, Order)
		}
		return Orders, nil
	}
}
