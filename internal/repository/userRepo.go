package repository

import (
	"context"
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make User db with map
type UserRepo struct {
	Users map[int]domain.User
	mutek *sync.Mutex
}

func NewUserRepo() UserRepoInterface {
	return UserRepo{
		Users: map[int]domain.User{},
		mutek: &sync.Mutex{},
	}
}

type UserRepoInterface interface {
	CreateUser
	GetUserByID
	GetUserByName
	UpdateUser
	DeleteUser
	GetAllUsers
	DecreaseBalance
}
type CreateUser interface {
	CreateUser(User *domain.User, kontek context.Context) (*domain.User, error)
}
type GetUserByID interface {
	GetUserByID(id int, kontek context.Context) (*domain.User, error)
}
type GetUserByName interface {
	GetUserByName(name string, kontek context.Context) (*domain.User, error)
}
type UpdateUser interface {
	UpdateUser(User *domain.User, kontek context.Context) error
}
type DeleteUser interface {
	DeleteUser(id int, kontek context.Context) error
}
type GetAllUsers interface {
	GetAllUsers(kontek context.Context) ([]domain.User, error)
}
type DecreaseBalance interface {
	DecreaseBalance(userID int, totalAmount float64, kontek context.Context) (*domain.User, error)
}

func (repo UserRepo) CreateUser(User *domain.User, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		if repo.Users == nil || len(repo.Users) == 0 {
			User.ID = 1
		} else {
			User.ID = repo.Users[len(repo.Users)].ID + 1
		}
		repo.Users[User.ID] = *User
		return User, nil
	}

}

func (repo UserRepo) GetUserByID(id int, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, User := range repo.Users {
			if User.ID == id {
				return &User, nil
			}
		}
		return nil, errors.New("THERE'S NO USER WITH THAT ID🤬🚨🤬🚨")
	}
}

func (repo UserRepo) GetUserByName(name string, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, User := range repo.Users {
			if User.Name == name {
				return &User, nil
			}
		}
		return nil, errors.New("THERE'S NO USER WITH THAT NAME🤬🚨🤬🚨")
	}
}

func (repo UserRepo) UpdateUser(User *domain.User, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return kontek.Err()
	default:
		if _, exist := repo.Users[User.ID]; !exist {
			return errors.New("THERE'S NO USER WITH THAT ID🤬🤬🤬🚨🚨")
		}
		repo.Users[User.ID] = *User
		return nil
	}
}

func (repo UserRepo) DeleteUser(id int, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return kontek.Err()
	default:
		if _, exist := repo.Users[id]; !exist {
			return errors.New("THERE'S NO USER WITH THAT ID🤬🚨🤬🚨")
		}
		delete(repo.Users, id)
		return nil
	}
}

func (repo UserRepo) GetAllUsers(kontek context.Context) ([]domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		Users := make([]domain.User, 0, len(repo.Users))
		for _, User := range repo.Users {
			Users = append(Users, User)
		}
		return Users, nil
	}
}

func (repo UserRepo) DecreaseBalance(userID int, totalAmount float64, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		user, exist := repo.Users[userID]
		if !exist {
			return nil, errors.New("THERE'S NO USER WITH THAT ID🤬🤬🤬🚨🚨")
		}
		if user.Balance < totalAmount {
			return nil, errors.New("INSUFFICIENT BALANCE🤬🤬🤬🚨🚨")
		}
		user.Balance -= totalAmount
		repo.Users[userID] = user
		return &user, nil
	}
}
