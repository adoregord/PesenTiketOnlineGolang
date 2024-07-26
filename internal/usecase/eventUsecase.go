package usecase

import (
	"context"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/repository"
)

// make a connection to repo
type EventUsecase struct {
	EventRepo repository.EventRepoInterface
}

func NewEventUsecase(eventRepo repository.EventRepoInterface) EventUsecaseInterface {
	return EventUsecase{
		EventRepo: eventRepo,
	}
}

type EventUsecaseInterface interface {
	CreateEvent
	GetEventByID
	GetEventByName
	UpdateEvent
	DeleteEvent
	GetAllEvents
}
type CreateEvent interface {
	CreateEvent(event domain.Event, kontek context.Context) (*domain.Event, error)
}
type GetEventByID interface {
	GetEventByID(id int, kontek context.Context) (*domain.Event, error)
}
type GetEventByName interface {
	GetEventByName(name string, kontek context.Context) (*domain.Event, error)
}
type UpdateEvent interface {
	UpdateEvent(event domain.Event, kontek context.Context) error
}
type DeleteEvent interface {
	DeleteEvent(id int, kontek context.Context) error
}
type GetAllEvents interface {
	GetAllEvents(kontek context.Context) ([]domain.Event, error)
}

func (uc EventUsecase) CreateEvent(event domain.Event, kontek context.Context) (*domain.Event, error) {
	return uc.EventRepo.CreateEvent(&event, kontek)
}
func (uc EventUsecase) GetEventByID(id int, kontek context.Context) (*domain.Event, error) {
	return uc.EventRepo.GetEventByID(id, kontek)
}
func (uc EventUsecase) GetEventByName(name string, kontek context.Context) (*domain.Event, error) {
	return uc.EventRepo.GetEventByName(name, kontek)
}
func (uc EventUsecase) UpdateEvent(event domain.Event, kontek context.Context) error {
	return uc.EventRepo.UpdateEvent(&event, kontek)
}
func (uc EventUsecase) DeleteEvent(id int, kontek context.Context) error {
	return uc.EventRepo.DeleteEvent(id, kontek)
}
func (uc EventUsecase) GetAllEvents(kontek context.Context) ([]domain.Event, error) {
	return uc.EventRepo.GetAllEvents(kontek)
}
