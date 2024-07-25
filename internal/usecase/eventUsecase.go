package usecase

import (
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
	CreateEvent(event domain.Event) error
}
type GetEventByID interface {
	GetEventByID(id int) (*domain.Event, error)
}
type GetEventByName interface {
	GetEventByName(name string) (domain.Event, error)
}
type UpdateEvent interface {
	UpdateEvent(event domain.Event) error
}
type DeleteEvent interface {
	DeleteEvent(id int) error
}
type GetAllEvents interface {
	GetAllEvents() ([]domain.Event, error)
}

func (uc EventUsecase) CreateEvent(event domain.Event) error {
	return uc.EventRepo.CreateEvent(&event)
}
func (uc EventUsecase) GetEventByID(id int) (*domain.Event, error) {
	return uc.EventRepo.GetEventByID(id)
}
func (uc EventUsecase) GetEventByName(name string) (domain.Event, error) {
	return uc.EventRepo.GetEventByName(name)
}
func (uc EventUsecase) UpdateEvent(event domain.Event) error {
	return uc.EventRepo.UpdateEvent(&event)
}
func (uc EventUsecase) DeleteEvent(id int) error {
	return uc.EventRepo.DeleteEvent(id)
}
func (uc EventUsecase) GetAllEvents() ([]domain.Event, error) {
	return uc.EventRepo.GetAllEvents()
}

