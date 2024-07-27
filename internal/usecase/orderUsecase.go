package usecase

import (
	"context"
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository"
	"time"
)

// make a connection to repo
type OrderUsecase struct {
	OrderRepo repository.OrderRepoInterface
	EventRepo repository.EventRepoInterface
	UserRepo  repository.UserRepoInterface
}

func NewOrderUsecase(orderRepo repository.OrderRepoInterface, eventRepo repository.EventRepoInterface, userRepo repository.UserRepoInterface) OrderUsecaseInterface {
	return OrderUsecase{
		OrderRepo: orderRepo,
		EventRepo: eventRepo,
		UserRepo:  userRepo,
	}
}

type OrderUsecaseInterface interface {
	CreateOrder
	GetOrderByID
	GetAllOrders
}
type CreateOrder interface {
	CreateOrder(orederReq domain.OrderRequest, kontek context.Context) (*domain.Order, error)
}
type GetOrderByID interface {
	GetOrderByID(id int, kontek context.Context) ([]domain.Order, error)
}
type GetAllOrders interface {
	GetAllOrders(kontek context.Context) ([]domain.Order, error)
}

func (uc OrderUsecase) CreateOrder(orderReq domain.OrderRequest, kontek context.Context) (*domain.Order, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		var order domain.Order

		// get event first from event repo get by ID
		event, err := uc.EventRepo.GetEventByID(orderReq.EventID, kontek)
		if err != nil {
			order.Status = "FAILED COULDN'T FIND EVENT ID"
			return nil, err
		}

		// get user details from user repo by ID
		user, err := uc.UserRepo.GetUserByID(orderReq.UserID, kontek)
		if err != nil {
			order.Status = "FAILED COULDN'T FIND USER ID"
			return &order, err
		}

		order.User.ID = user.ID
		order.User.Name = user.Name
		order.Event.Name = event.Name
		order.Event.Date = event.Date
		order.Event.Location = event.Location
		order.Event.Description = event.Description

		defer func() {
			order.OrderDate = time.Now().Format("02-Jan-2006 15:04:05")
			uc.OrderRepo.CreateOrder(&order, kontek)
		}()

		// then we extract the ticket the user wants and append it to order
		var ticketTemp []domain.Ticket      // for storing updated stock
		var orderTicketTemp []domain.Ticket // for storing updated order stock
		for _, value := range event.Ticket {
			for _, value2 := range orderReq.Ticket {
				if value.ID == value2.ID || value.Type == value2.Type {
					// check if the quantity the user wants and the stock is enough
					if value.Quantity < value2.Quantity {
						order.Status = "FAILED STOCK TICKET NOT ENOUGH"
						return &order, errors.New("TICKET STOCK NOT ENOUGH")
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

		//check if the temp ticket is empty
		if len(ticketTemp) == 0 {
			order.Status = "FAILED TICKET NOT FOUND"
			return &order, errors.New("TICKET NOT FOUND")
		}

		order.PaymentMethod = "QRIS"
		order.Status = "PENDING"

		//check if the user has enough money to buy the tickets
		if user.Balance < order.TotalPrice {
			order.Status = "PAYMENT FAILED"
			return &order, errors.New("INSUFFICIENT BALANCE")
		}

		// then we reduce the stock on event db if the payment is success
		event.Ticket = ticketTemp
		uc.EventRepo.UpdateEvent(event, kontek)

		// then we update the balance
		user.Balance -= order.TotalPrice
		uc.UserRepo.UpdateUser(user, kontek)
		order.Status = "SUCCESS"

		return &order, nil
	}
}
func (uc OrderUsecase) GetOrderByID(userID int, kontek context.Context) ([]domain.Order, error) {
	return uc.OrderRepo.GetOrderByID(userID, kontek)
}
func (uc OrderUsecase) GetAllOrders(kontek context.Context) ([]domain.Order, error) {
	return uc.OrderRepo.GetAllOrders(kontek)
}
