package repository

import (
	"errors"
	"pemesananTiketOnlineGo/internal/domain"
)

// make event db with map
type EventRepo struct {
	Events map[int]domain.Event
}

func NewEventRepo() EventRepoInterface {
	return EventRepo{
		Events: map[int]domain.Event{},
	}
}

type EventRepoInterface interface {
	CreateEvent
	GetEventByID
	GetEventByName
	UpdateEvent
	DeleteEvent
	GetAllEvents
}
type CreateEvent interface {
	CreateEvent(event *domain.Event) error
}
type GetEventByID interface {
	GetEventByID(id int) (*domain.Event, error)
}
type GetEventByName interface {
	GetEventByName(name string) (domain.Event, error)
}
type UpdateEvent interface {
	UpdateEvent(event *domain.Event) error
}
type DeleteEvent interface {
	DeleteEvent(id int) error
}
type GetAllEvents interface {
	GetAllEvents() ([]domain.Event, error)
}

func (repo EventRepo) CreateEvent(event *domain.Event) error {
	for _, value := range repo.Events {
		if value.Name == event.Name {
			return errors.New("EVENT WITH THAT NAME ALREADY EXIST")
		}
	}

	if repo.Events == nil || len(repo.Events) == 0 {
		event.ID = 1
	} else {
		event.ID = repo.Events[len(repo.Events)].ID + 1
	}
	repo.Events[event.ID] = *event
	return nil
}

func (repo EventRepo) GetEventByID(id int) (*domain.Event, error) {
	for _, event := range repo.Events {
		if event.ID == id {
			return &event, nil
		}
	}
	return nil, errors.New("THERE'S NO EVENT WITH THAT ID")
}

func (repo EventRepo) GetEventByName(name string) (domain.Event, error) {
	for _, event := range repo.Events {
		if event.Name == name {
			return event, nil
		}
	}
	return domain.Event{}, nil
}

func (repo EventRepo) UpdateEvent(event *domain.Event) error {
	if _, exist := repo.Events[event.ID]; !exist{
		return errors.New("THERE'S NO EVENT WITH THAT ID🤬🤬🤬🚨🚨")
	}
	repo.Events[event.ID] = *event
	return nil
}

func (repo EventRepo) DeleteEvent(id int) error {
	if _, exist := repo.Events[id]; !exist{
		return errors.New("THERE'S NO EVENT WITH THAT ID🤬🚨🤬🚨")
	}
	delete(repo.Events, id)
	return nil
}

func (repo EventRepo) GetAllEvents() ([]domain.Event, error) {
	events := make([]domain.Event, 0, len(repo.Events))
	for _, event := range repo.Events {
		events = append(events, event)
	}	
	return events, nil
}
