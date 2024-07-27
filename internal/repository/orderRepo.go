package repository

import (
	"context"
	"errors"
	"fmt"
	"pemesananTiketOnlineGo/internal/domain"
	"time"
)

// make Order db with map
type OrderRepo struct {
	Orders    map[int]domain.Order
	EventRepo EventRepoInterface
	UserRepo  UserRepoInterface
}

func NewOrderRepo(eventRepo EventRepoInterface, userRepo UserRepoInterface) OrderRepoInterface {
	return OrderRepo{
		Orders:    map[int]domain.Order{},
		EventRepo: eventRepo,
		UserRepo:  userRepo,
	}
}

type OrderRepoInterface interface {
	CreateOrder
	GetOrderByID
	GetAllOrders
}
type CreateOrder interface {
	CreateOrder(orderReq domain.OrderRequest, kontek context.Context) (*domain.Order, error)
}
type GetOrderByID interface {
	GetOrderByID(id int, kontek context.Context) (*domain.Order, error)
}
type GetAllOrders interface {
	GetAllOrders(kontek context.Context) ([]domain.Order, error)
}

func (repo OrderRepo) CreateOrder(orderReq domain.OrderRequest, kontek context.Context) (*domain.Order, error) {

	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		var order domain.Order

		defer func() {
			repo.Orders[order.ID] = order
		}()

		// auto increment ID
		if repo.Orders == nil || len(repo.Orders) == 0 {
			order.ID = 1
		} else {
			order.ID = repo.Orders[len(repo.Orders)].ID + 1
		}

		// get event first from event repo get by ID
		fmt.Println(orderReq.EventID)
		event, err := repo.EventRepo.GetEventByID(orderReq.EventID, kontek)
		if err != nil {
			return nil, err
		}

		// get user details from user repo by ID
		user, err := repo.UserRepo.GetUserByID(orderReq.UserID, kontek)
		if err != nil {
			order.Status = "FAILED COULDN'T FIND USER ID"
			return &order, err
		}

		order.Date = time.Now().Format("02-Jan-2006 15:04:05")
		order.User.ID = user.ID
		order.User.Name = user.Name
		order.Event.Name = event.Name
		order.Event.Date = event.Date
		order.Event.Location = event.Location

		// then we extract the ticket the user wants and append it to order
		var ticketTemp []domain.Ticket      // for storing updated stock
		var orderTicketTemp []domain.Ticket // for storing updated order stock
		for _, value := range event.Ticket {
			for _, value2 := range orderReq.Ticket {
				if value.Type == value2.Type || value.ID == value2.ID {
					// check if the quantity the user wants and the stock is enough
					if value.Quantity < value2.Quantity {
						order.Status = "FAILED STOCK TICKET NOT ENOUGH"
						return &order, errors.New("TICKET NOT ENOUGH")
					}
					order.EventTicket = append(order.EventTicket, value2)
					// we need to assign the type if the user input is only ticket ID
					// or if the user input the wrong type
					value2.Type = value.Type
					orderTicketTemp = append(orderTicketTemp, value2)
					// reduce the quantity and save it to ticket temp
					value.Quantity -= value2.Quantity
					ticketTemp = append(ticketTemp, value)
					order.TotalPrice += value.Price * float64(value2.Quantity)
				}
			}
		}
		// update the order ticket's type
		order.EventTicket = orderTicketTemp

		order.PaymentMethod = "QRIS"
		order.Status = "PENDING"

		//check if the user has enough money to buy the tickets
		if user.Balance < order.TotalPrice {
			order.Status = "PAYMENT FAILED"
			return &order, errors.New("INSUFFICIENT BALANCE")
		}

		// then we reduce the stock on event db if the payment is success
		event.Ticket = ticketTemp
		repo.EventRepo.UpdateEvent(event, kontek)

		// then we update the balance
		user.Balance -= order.TotalPrice
		repo.UserRepo.UpdateUser(user, kontek)
		order.Status = "SUCCESS"

		return &order, nil
	}
}

func (repo OrderRepo) GetOrderByID(id int, kontek context.Context) (*domain.Order, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, Order := range repo.Orders {
			if Order.ID == id {
				return &Order, nil
			}
		}
		return nil, errors.New("THERE'S NO ORDER WITH THAT ID")
	}
}

func (repo OrderRepo) GetAllOrders(kontek context.Context) ([]domain.Order, error) {
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
