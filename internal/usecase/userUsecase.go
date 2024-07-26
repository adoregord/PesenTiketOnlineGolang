package usecase

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository"
)

// make a connection to repo
type UserUsecase struct {
	UserRepo repository.UserRepoInterface
}

func NewUserUsecase(UserRepo repository.UserRepoInterface) UserUsecaseInterface {
	return UserUsecase{
		UserRepo: UserRepo,
	}
}

type UserUsecaseInterface interface {
	CreateUser
	GetUserByID
	GetUserByName
	UpdateUser
	DeleteUser
	GetAllUsers
}
type CreateUser interface {
	CreateUser(User domain.User, kontek context.Context) (*domain.User, error)
}
type GetUserByID interface {
	GetUserByID(id int, kontek context.Context) (*domain.User, error)
}
type GetUserByName interface {
	GetUserByName(name string, kontek context.Context) (*domain.User, error)
}
type UpdateUser interface {
	UpdateUser(User domain.User, kontek context.Context) error
}
type DeleteUser interface {
	DeleteUser(id int, kontek context.Context) error
}
type GetAllUsers interface {
	GetAllUsers(kontek context.Context) ([]domain.User, error)
}

func (uc UserUsecase) CreateUser(User domain.User, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.CreateUser(&User, kontek)
}
func (uc UserUsecase) GetUserByID(id int, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.GetUserByID(id, kontek)
}
func (uc UserUsecase) GetUserByName(name string, kontek context.Context) (*domain.User, error) {
	return uc.UserRepo.GetUserByName(name, kontek)
}
func (uc UserUsecase) UpdateUser(User domain.User, kontek context.Context) error {
	return uc.UserRepo.UpdateUser(&User, kontek)
}
func (uc UserUsecase) DeleteUser(id int, kontek context.Context) error {
	return uc.UserRepo.DeleteUser(id, kontek)
}
func (uc UserUsecase) GetAllUsers(kontek context.Context) ([]domain.User, error) {
	return uc.UserRepo.GetAllUsers(kontek)
}