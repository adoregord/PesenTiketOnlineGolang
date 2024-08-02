package usecase2

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository2"
)

// make a connection to repo
type OrderUsecase struct {
	OrderRepo repository2.OrderRepoInterface
	EventRepo repository2.EventRepoInterface
	UserRepo  repository2.UserRepoInterface
}

func NewOrderUsecase(orderRepo repository2.OrderRepoInterface, eventRepo repository2.EventRepoInterface, userRepo repository2.UserRepoInterface) OrderUsecaseInterface {
	return OrderUsecase{
		OrderRepo: orderRepo,
		EventRepo: eventRepo,
		UserRepo:  userRepo,
	}
}

type OrderUsecaseInterface interface {
	CreateOrderDB
}

type CreateOrderDB interface {
	CreateOrderDB(orderReq *domain.OrderRequest, kontek context.Context) (*domain.Order, error)
}

func (uc OrderUsecase) CreateOrderDB(orderReq *domain.OrderRequest, kontek context.Context) (*domain.Order, error) {
	var order domain.Order

	// this one is just for checking the user so i can use defer func
	user, err := uc.UserRepo.GetUserIDDB(orderReq.UserID, kontek)
	if err != nil {
		return nil, err
	}

	// get event info first from db to check if the event is there or not
	event, err := uc.EventRepo.ViewEventByIdDB(orderReq.EventID, kontek)
	if err != nil {
		return nil, err
	}
	for i, orderTicket := range orderReq.Ticket {
		for _, eventTicket := range event.Ticket {
			if orderTicket.ID == eventTicket.ID {
				orderReq.Ticket[i].Type = eventTicket.Type
			}
		}
	}
	
	total, err := uc.EventRepo.CheckTotalValueDB(orderReq.EventID, orderReq.Ticket, kontek)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			order.Status = "FAILED " + err.Error()
			order.PaymentMethod = "-"
		}
		order.User.ID = orderReq.UserID
		order.User.Name = user.Name
		order.Event.ID = orderReq.EventID
		order.Event.Name = event.Name
		order.Event.Date = event.Date
		order.Event.Location = event.Location
		order.Event.Description = event.Description
		order.TotalPrice = total
		order.EventTicket = orderReq.Ticket
		uc.OrderRepo.CreateOrderDB(&order, kontek)
	}()

	// check and decrease user's balance
	_, err = uc.UserRepo.DecreaseBalanceDB(orderReq.UserID, total, kontek)
	if err != nil {
		return &order, err
	}

	// decrease stock ticket in event repo
	err = uc.EventRepo.DecrementTicketStockDB(orderReq.EventID, orderReq.Ticket, kontek)
	if err != nil {
		return &order, err
	}

	order.PaymentMethod = "QRIS"
	order.Status = "SUCCESS"

	return &order, nil
}
