package usecase2

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository2"
)

// make a connection to repo
type UserUsecase struct {
	UserRepo repository2.UserRepoInterface
}

func NewUserUsecase(UserRepo repository2.UserRepoInterface) UserUsecaseInterface {
	return UserUsecase{
		UserRepo: UserRepo,
	}
}

type UserUsecaseInterface interface {
	CreateUserDB
	GetAllUserDB
	GetUserIDDB
	AddBalanceDB
	DeleteUserDB
}
type CreateUserDB interface {
	CreateUserDB(user *domain.User, kontek context.Context) (*domain.User, error)
}
type GetAllUserDB interface {
	GetAllUserDB(kontek context.Context) ([]domain.User, error)
}
type GetUserIDDB interface {
	GetUserIDDB(userID int, kontek context.Context) (*domain.User, error)
}
type AddBalanceDB interface {
	AddBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error)
}
type DeleteUserDB interface {
	DeleteUserDB(userID int, kontek context.Context) error
}

func (uc UserUsecase) CreateUserDB(user *domain.User, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.CreateUserDB(user, kontek)
}
func (uc UserUsecase) GetAllUserDB(kontek context.Context) ([]domain.User, error) {
	return uc.UserRepo.GetAllUserDB(kontek)
}
func (uc UserUsecase) GetUserIDDB(userID int, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.GetUserIDDB(userID, kontek)
}
func (uc UserUsecase) AddBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.AddBalanceDB(userID, totalAmount, kontek)
}
func (uc UserUsecase) DeleteUserDB(userID int, kontek context.Context) error {
	return uc.UserRepo.DeleteUserDB(userID, kontek)
}
