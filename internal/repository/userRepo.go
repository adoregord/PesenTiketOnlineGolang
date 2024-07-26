package repository

import (
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
)

// make User db with map
type UserRepo struct {
	Users map[int]domain.User
}

func NewUserRepo() UserRepoInterface {
	return UserRepo{
		Users: map[int]domain.User{},
	}
}

type UserRepoInterface interface {
	CreateUser
	GetUserByID
	GetUserByName
	UpdateUser
	DeleteUser
	GetAllUsers
}
type CreateUser interface {
	CreateUser(User *domain.User) error
}
type GetUserByID interface {
	GetUserByID(id int) (*domain.User, error)
}
type GetUserByName interface {
	GetUserByName(name string) (*domain.User, error)
}
type UpdateUser interface {
	UpdateUser(User *domain.User) error
}
type DeleteUser interface {
	DeleteUser(id int) error
}
type GetAllUsers interface {
	GetAllUsers() ([]domain.User, error)
}

func (repo UserRepo) CreateUser(User *domain.User) error {
	if repo.Users == nil || len(repo.Users) == 0 {
		User.ID = 1
	} else {
		User.ID = repo.Users[len(repo.Users)].ID + 1
	}
	repo.Users[User.ID] = *User
	return nil
}

func (repo UserRepo) GetUserByID(id int) (*domain.User, error) {
	for _, User := range repo.Users {
		if User.ID == id {
			return &User, nil
		}
	}
	return nil, errors.New("THERE'S NO USER WITH THAT ID")
}

func (repo UserRepo) GetUserByName(name string) (*domain.User, error) {
	for _, User := range repo.Users {
		if User.Name == name {
			return &User, nil
		}
	}
	return nil, errors.New("THERE'S NO USER WITH THAT NAME🤬🚨🤬🚨")
}

func (repo UserRepo) UpdateUser(User *domain.User) error {
	if _, exist := repo.Users[User.ID]; !exist {
		return errors.New("THERE'S NO USER WITH THAT ID🤬🤬🤬🚨🚨")
	}
	repo.Users[User.ID] = *User
	return nil
}

func (repo UserRepo) DeleteUser(id int) error {
	if _, exist := repo.Users[id]; !exist {
		return errors.New("THERE'S NO USER WITH THAT ID🤬🚨🤬🚨")
	}
	delete(repo.Users, id)
	return nil
}

func (repo UserRepo) GetAllUsers() ([]domain.User, error) {
	Users := make([]domain.User, 0, len(repo.Users))
	for _, User := range repo.Users {
		Users = append(Users, User)
	}
	return Users, nil
}
