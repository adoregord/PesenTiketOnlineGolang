package repository2

import (
	"context"
	"database/sql"
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make User db with map
type UserRepo struct {
	db    *sql.DB
	mutek *sync.Mutex
}

func NewUserRepo(db *sql.DB) UserRepoInterface {
	return UserRepo{
		db:    db,
		mutek: &sync.Mutex{},
	}
}

type UserRepoInterface interface {
	CreateUserDB
	GetAllUserDB
	GetUserIDDB
	DecreaseBalanceDB
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
type DecreaseBalanceDB interface {
	DecreaseBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error)
}
type AddBalanceDB interface {
	AddBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error)
}
type DeleteUserDB interface {
	DeleteUserDB(userID int, kontek context.Context) error
}

func (repo UserRepo) CreateUserDB(user *domain.User, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()
	query := `	
	insert into akun (name, balance)
	values ($1, $2)
	RETURNING id
	`

	err := repo.db.QueryRowContext(kontek, query, user.Name, user.Balance).Scan(&user.ID)
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (repo UserRepo) GetAllUserDB(kontek context.Context) ([]domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	var users []domain.User

	query := `
	SELECT *
	FROM akun`

	rows, err := repo.db.QueryContext(kontek, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo UserRepo) GetUserIDDB(userID int, kontek context.Context) (*domain.User, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	var user domain.User

	query := `
	SELECT name, balance 
	FROM akun
	WHERE id = $1
	`

	row := repo.db.QueryRowContext(kontek, query, userID)

	err := row.Scan(
		&user.Name,
		&user.Balance,
	)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	return &user, err
}

func (repo UserRepo) DecreaseBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error) {
	var user domain.User

	query := `
	SELECT id, name, balance
	FROM akun 
	WHERE id = $1
	`

	err := repo.db.QueryRowContext(kontek, query, userID).Scan(&user.ID, &user.Name, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("THERE'S NO USER WITH THAT ID🤬🤬🤬🚨🚨")
		}
		return nil, err
	}

	if user.Balance < totalAmount {
		return nil, errors.New("INSUFFICIENT BALANCE🤬🤬🤬🚨🚨")
	}

	user.Balance -= totalAmount

	// update the user
	query2 := `
	UPDATE akun
	SET balance = $1
	WHERE id = $2
	`
	_, err = repo.db.ExecContext(kontek, query2, user.Balance, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (repo UserRepo) AddBalanceDB(userID int, totalAmount float64, kontek context.Context) (*domain.User, error) {
	var user domain.User

	query := `
	SELECT id, name, balance
	FROM akun 
	WHERE id = $1
	`

	err := repo.db.QueryRowContext(kontek, query, userID).Scan(&user.ID, &user.Name, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("THERE'S NO USER WITH THAT ID🤬🤬🤬🚨🚨")
		}
		return nil, err
	}

	user.Balance += totalAmount

	// update the user
	query2 := `
	UPDATE akun
	SET balance = $1
	WHERE id = $2
	`
	_, err = repo.db.ExecContext(kontek, query2, user.Balance, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo UserRepo) DeleteUserDB(userID int, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	var check int

	// query for deleting user
	query := `
	delete from akun where id = $1
	returning id
	`

	err := repo.db.QueryRowContext(kontek, query, userID).Scan(&check)
	if err != nil {
		return errors.New("NO USER WITH SUCH ID🚨🤬🚨🤬")
	}

	return nil
}
