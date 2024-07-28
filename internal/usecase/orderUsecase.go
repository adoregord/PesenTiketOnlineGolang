package usecase

import (
	"context"
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

		// check if the stock ticket is available and get the total value
		total, err := uc.EventRepo.CheckTotalValue(orderReq.EventID, orderReq.Ticket, kontek)
		if err != nil {
			order.Status = err.Error()
			return &order, err
		}

		// get user details and decrease the balance from user
		user, err := uc.UserRepo.DecreaseBalance(orderReq.UserID, total, kontek)
		if err != nil {
			order.Status = err.Error()
			return &order, err
		}

		defer func() {
			order.OrderDate = time.Now().Format("02-Jan-2006 15:04:05")
			order.User.ID = user.ID
			order.User.Name = user.Name
			order.Event.Name = event.Name
			order.Event.Date = event.Date
			order.Event.Location = event.Location
			order.Event.Description = event.Description
			uc.OrderRepo.CreateOrder(&order, kontek)
		}()

		// decrease the total amount of ticket
		uc.EventRepo.DecrementTicketStock(orderReq.EventID, orderReq.Ticket, kontek)

		order.PaymentMethod = "QRIS"
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
