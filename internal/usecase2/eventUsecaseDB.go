package usecase2

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository2"
)

// make a connection to repo
type EventUsecase struct {
	EventRepo  repository2.EventRepoInterface
	TicketRepo repository2.TicketRepoInterface
}

func NewEventUsecase(eventRepo repository2.EventRepoInterface, ticketRepo repository2.TicketRepoInterface) EventUsecaseInterface {
	return EventUsecase{
		EventRepo:  eventRepo,
		TicketRepo: ticketRepo,
	}
}

type EventUsecaseInterface interface {
	CreateEventDB
	ViewAllEvents
	ViewEventByIdDB
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
type DeleteEventDB interface {
	DeleteEventDB(eventID int, kontek context.Context) error
}
type UpdateEventTicketDB interface {
	UpdateEventTicketDB(eventID int, tickets *[]domain.Ticket, kontek context.Context) error
}

func (uc EventUsecase) CreateEventDB(event *domain.Event, kontek context.Context) error {
	err := uc.TicketRepo.CreateTicketDB(&event.Ticket, kontek)
	if err != nil {
		return err
	}

	return uc.EventRepo.CreateEventDB(event, kontek)
}

func (uc EventUsecase) ViewAllEvents(kontek context.Context) ([]domain.Event, error) {
	return uc.EventRepo.ViewAllEvents(kontek)
}

func (uc EventUsecase) ViewEventByIdDB(eventID int, kontek context.Context) (*domain.Event, error) {
	return uc.EventRepo.ViewEventByIdDB(eventID, kontek)
}

func (uc EventUsecase) DeleteEventDB(eventID int, kontek context.Context) error {
	return uc.EventRepo.DeleteEventDB(eventID, kontek)
}

func (uc EventUsecase) UpdateEventTicketDB(eventID int, tickets *[]domain.Ticket, kontek context.Context) error {
	err := uc.TicketRepo.CreateTicketDB(tickets, kontek)
	if err != nil {
		return err
	}

	return uc.EventRepo.UpdateEventTicketDB(eventID, tickets, kontek)
}
