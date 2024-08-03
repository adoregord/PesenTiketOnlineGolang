package repository2

import (
	"context"
	"database/sql"
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
	"sync"
)

// make event db with map
type EventRepo struct {
	db    *sql.DB
	mutek *sync.Mutex
}

func NewEventRepo(db *sql.DB) EventRepoInterface {
	return EventRepo{
		db:    db,
		mutek: &sync.Mutex{},
	}
}

type EventRepoInterface interface {
	CreateEventDB
	ViewAllEvents
	ViewEventByIdDB
	CheckTotalValueDB
	DecrementTicketStockDB
	DeleteEventDB
	UpdateEventTicketDB
}
type CreateEventDB interface {
	CreateEventDB(event *domain.Event, kontek context.Context) error
}
type ViewAllEvents interface {
	ViewAllEvents(kontek context.Context) ([]domain.Event, error)
}
type ViewEventByIdDB interface {
	ViewEventByIdDB(eventID int, kontek context.Context) (*domain.Event, error)
}
type CheckTotalValueDB interface {
	CheckTotalValueDB(eventID int, tickets []domain.TicketReq, kontek context.Context) (float64, error)
}
type DecrementTicketStockDB interface {
	DecrementTicketStockDB(eventID int, tickets []domain.TicketReq, kontek context.Context) error
}
type DeleteEventDB interface {
	DeleteEventDB(eventID int, kontek context.Context) error
}
type UpdateEventTicketDB interface {
	UpdateEventTicketDB(eventID int, tickets *[]domain.Ticket, kontek context.Context) error
}

func (repo EventRepo) CreateEventDB(event *domain.Event, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	tx, err := repo.db.BeginTx(kontek, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO event (name, date, description, location)
		VAlUES ($1, $2, $3, $4)
		RETURNING id
		`
	stmt, err := tx.PrepareContext(kontek, query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	query2 := `
		INSERT into event_ticket (event_id, ticket_id)
		values ($1, $2)
		`
	stmt2, err := tx.PrepareContext(kontek, query2)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt2.Close()

	err = stmt.QueryRowContext(kontek, event.Name, event.Date, event.Description, event.Location).Scan(&event.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, value := range event.Ticket {
		_, err = stmt2.ExecContext(kontek, event.ID, value.ID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo EventRepo) ViewAllEvents(kontek context.Context) ([]domain.Event, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	query := `
	SELECT 
		e.id AS event_id,
		e.name AS event_name,
		e.date AS event_date,
		e.description AS event_description,
		e.location AS event_location,
		t.id AS ticket_id,
		t.type AS ticket_type,
		t.quantity AS ticket_quantity,
		t.price AS ticket_price
	FROM 
		event e
	JOIN 
		event_ticket et ON e.id = et.event_id
	JOIN 
		ticket t ON et.ticket_id = t.id;
	`
	rows, err := repo.db.QueryContext(kontek, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[int]*domain.Event)

	for rows.Next() {
		var (
			eventID     int
			eventName   string
			eventDate   string
			eventDesc   string
			eventLoc    string
			ticketID    int
			ticketType  string
			ticketQty   int
			ticketPrice float64
		)

		err := rows.Scan(
			&eventID,
			&eventName,
			&eventDate,
			&eventDesc,
			&eventLoc,
			&ticketID,
			&ticketType,
			&ticketQty,
			&ticketPrice,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := eventMap[eventID]; !exists {
			eventMap[eventID] = &domain.Event{
				ID:          eventID,
				Name:        eventName,
				Date:        eventDate,
				Description: eventDesc,
				Location:    eventLoc,
				Ticket:      []domain.Ticket{},
			}
		}

		eventMap[eventID].Ticket = append(eventMap[eventID].Ticket, domain.Ticket{
			ID:       ticketID,
			Type:     ticketType,
			Quantity: ticketQty,
			Price:    ticketPrice,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var events []domain.Event
	for _, event := range eventMap {
		events = append(events, *event)
	}

	return events, nil
}

func (repo EventRepo) ViewEventByIdDB(eventID int, kontek context.Context) (*domain.Event, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	var event domain.Event
	var ticket domain.Ticket

	query := `
	SELECT 
		e.id AS event_id,
		e.name AS event_name,
		e.date AS event_date,
		e.description AS event_description,
		e.location AS event_location,
		t.id AS ticket_id,
		t.type AS ticket_type,
		t.quantity AS ticket_quantity,
		t.price AS ticket_price
	FROM 
		event e
	JOIN 
		event_ticket et ON e.id = et.event_id
	JOIN 
		ticket t ON et.ticket_id = t.id
	WHERE event_id = $1;
	`

	rows, err := repo.db.QueryContext(kontek, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Date,
			&event.Description,
			&event.Location,
			&ticket.ID,
			&ticket.Type,
			&ticket.Quantity,
			&ticket.Price,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	event.Ticket = tickets

	return &event, err
}

func (repo EventRepo) CheckTotalValueDB(eventID int, tickets []domain.TicketReq, kontek context.Context) (float64, error) {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	// Memulai transaksi
	tx, err := repo.db.BeginTx(kontek, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Query untuk mendapatkan tiket event
	query := `
	SELECT 
		t.type,
		t.price
	FROM 
		event_ticket et
	JOIN 
		ticket t ON et.ticket_id = t.id
	WHERE 
		t.id = $1
	`

	stmt, err := tx.PrepareContext(kontek, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var total float64
	for _, ticket := range tickets {
		row := stmt.QueryRowContext(kontek, ticket.ID)

		err := row.Scan(&ticket.Type, &ticket.Price)
		if err != nil {
			return 0, nil
		}
		total += ticket.Price * (float64(ticket.Quantity))
	}

	return total, nil
}

func (repo EventRepo) DecrementTicketStockDB(eventID int, tickets []domain.TicketReq, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	tx, err := repo.db.BeginTx(kontek, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `
	SELECT
		t.id,
		t.type,
		t.quantity
	FROM
		event_ticket et
	JOIN 
		ticket as t ON et.ticket_id = t.id
	WHERE
		et.event_id = $1
	`

	stmt, err := tx.PrepareContext(kontek, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(kontek, eventID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// scan to eventTickets
	eventTickets := make(map[int]domain.Ticket)
	for rows.Next() {
		var ticket domain.Ticket
		err := rows.Scan(&ticket.ID, &ticket.Type, &ticket.Quantity)
		if err != nil {
			return err
		}
		eventTickets[ticket.ID] = ticket
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// check the ticket availability and update the quantities
	for _, ticket := range tickets {
		eventTickets, exist := eventTickets[ticket.ID]
		if !exist {
			return errors.New("ticket not found for event")
		}
		if eventTickets.Quantity < ticket.Quantity {
			return errors.New("NOT ENOUGH TICKET STOCK🤬🚨🤬🚨")
		}
		eventTickets.Quantity -= ticket.Quantity

		// update the ticket stock to database
		query2 := `
		UPDATE
			ticket
		SET
			quantity = $1
		WHERE 
			id = $2
		`
		_, err := tx.ExecContext(kontek, query2, eventTickets.Quantity, eventTickets.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// repo func for deleting event
func (repo EventRepo) DeleteEventDB(eventID int, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	var check int

	// query for deleting event DB
	query := `
	DELETE FROM event WHERE id = $1
	returning id
	`

	err := repo.db.QueryRowContext(kontek, query, eventID).Scan(&check)
	if err != nil {
		return errors.New("NO EVENT WITH SUCH ID🚨🤬🚨🤬")
	}

	return nil
}

// repo func for updating event ticket
func (repo EventRepo) UpdateEventTicketDB(eventID int, tickets *[]domain.Ticket, kontek context.Context) error {
	repo.mutek.Lock()
	defer repo.mutek.Unlock()

	tx, err := repo.db.BeginTx(kontek, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := `
		INSERT into event_ticket (event_id, ticket_id)
		values ($1, $2)
		`
	stmt, err := tx.PrepareContext(kontek, query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, value := range *tickets {
		_, err = stmt.ExecContext(kontek, eventID, value.ID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
