package usecase

import (
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
	CreateUser(User domain.User) error
}
type GetUserByID interface {
	GetUserByID(id int) (*domain.User, error)
}
type GetUserByName interface {
	GetUserByName(name string) (*domain.User, error)
}
type UpdateUser interface {
	UpdateUser(User domain.User) error
}
type DeleteUser interface {
	DeleteUser(id int) error
}
type GetAllUsers interface {
	GetAllUsers() ([]domain.User, error)
}

func (uc UserUsecase) CreateUser(User domain.User) error {
	return uc.UserRepo.CreateUser(&User)
}
func (uc UserUsecase) GetUserByID(id int) (*domain.User, error) {
	return uc.UserRepo.GetUserByID(id)
}
func (uc UserUsecase) GetUserByName(name string) (*domain.User, error) {
	return uc.UserRepo.GetUserByName(name)
}
func (uc UserUsecase) UpdateUser(User domain.User) error {
	return uc.UserRepo.UpdateUser(&User)
}
func (uc UserUsecase) DeleteUser(id int) error {
	return uc.UserRepo.DeleteUser(id)
}
func (uc UserUsecase) GetAllUsers() ([]domain.User, error) {
	return uc.UserRepo.GetAllUsers()
}