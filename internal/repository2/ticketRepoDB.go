package repository2

import (
	"context"
	"database/sql"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make event db with map
type TicketRepo struct {
	db    *sql.DB
	mutek *sync.Mutex
}

func NewTicketRepo(db *sql.DB) TicketRepoInterface {
	return TicketRepo{
		db:    db,
		mutek: &sync.Mutex{},
	}
}

type TicketRepoInterface interface {
	CreateTicketDB
}
type CreateTicketDB interface {
	CreateTicketDB(ticket *[]domain.Ticket, kontek context.Context) error
}

func (repo TicketRepo) CreateTicketDB(ticket *[]domain.Ticket, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	tx, err := repo.db.BeginTx(kontek, nil) // Memulai transaksi
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(kontek, `
			INSERT INTO ticket (type, quantity, price)
			VALUES ($1, $2, $3)
			RETURNING id
		`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close() // Menutup prepared statement setelah selesai

	for i, v := range *ticket {
		var id int
		err := stmt.QueryRowContext(kontek, v.Type, v.Quantity, v.Price).Scan(&id)
		if err != nil {
			tx.Rollback() // Membatalkan transaksi jika ada kesalahan
			return err
		}
		(*ticket)[i].ID = id // Menyimpan ID yang dihasilkan ke dalam struct ticket
	}

	if err = tx.Commit(); err != nil { // Komit transaksi
		return err
	}

	return nil
}
