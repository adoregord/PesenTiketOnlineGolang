package usecase

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository"
)

// make a connection to repo
type OrderUsecase struct {
	OrderRepo repository.OrderRepoInterface
}

func NewOrderUsecase(orderRepo repository.OrderRepoInterface) OrderUsecaseInterface {
	return OrderUsecase{
		OrderRepo: orderRepo,
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
	GetOrderByID(id int, kontek context.Context) (*[]domain.Order, error)
}
type GetAllOrders interface {
	GetAllOrders(kontek context.Context) ([]domain.Order, error)
}

func (uc OrderUsecase) CreateOrder(orederReq domain.OrderRequest, kontek context.Context) (*domain.Order, error) {
	return uc.OrderRepo.CreateOrder(orederReq, kontek)
}
func (uc OrderUsecase) GetOrderByID(userID int, kontek context.Context) (*[]domain.Order, error) {
	return uc.OrderRepo.GetOrderByID(userID, kontek)
}
func (uc OrderUsecase) GetAllOrders(kontek context.Context) ([]domain.Order, error) {
	return uc.OrderRepo.GetAllOrders(kontek)
}
