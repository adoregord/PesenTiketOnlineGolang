package repository

import (
	"context"
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
	CreateEvent(event *domain.Event, kontek context.Context) (*domain.Event, error)
}
type GetEventByID interface {
	GetEventByID(id int, kontek context.Context) (*domain.Event, error)
}
type GetEventByName interface {
	GetEventByName(name string, kontek context.Context) (*domain.Event, error)
}
type UpdateEvent interface {
	UpdateEvent(event *domain.Event, kontek context.Context) error
}
type DeleteEvent interface {
	DeleteEvent(id int, kontek context.Context) error
}
type GetAllEvents interface {
	GetAllEvents(kontek context.Context) ([]domain.Event, error)
}

func (repo EventRepo) CreateEvent(event *domain.Event, kontek context.Context) (*domain.Event, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, value := range repo.Events {
			if value.Name == event.Name {
				return nil, errors.New("EVENT WITH THAT NAME ALREADY EXIST")
			}
		}

		if repo.Events == nil || len(repo.Events) == 0 {
			event.ID = 1
		} else {
			event.ID = repo.Events[len(repo.Events)].ID + 1
		}
		repo.Events[event.ID] = *event
		return event, nil
	}
}

func (repo EventRepo) GetEventByID(id int, kontek context.Context) (*domain.Event, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, event := range repo.Events {
			if event.ID == id {
				return &event, nil
			}
		}
		return nil, errors.New("THERE'S NO EVENT WITH THAT ID")
	}
}

func (repo EventRepo) GetEventByName(name string, kontek context.Context) (*domain.Event, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		for _, event := range repo.Events {
			if event.Name == name {
				return &event, nil
			}
		}
		return nil, errors.New("THERE'S NO EVENT WITH THAT NAME🤬🚨🤬🚨")
	}
}

func (repo EventRepo) UpdateEvent(event *domain.Event, kontek context.Context) error {
	select {
	case <-kontek.Done():
		return kontek.Err()
	default:
		if _, exist := repo.Events[event.ID]; !exist {
			return errors.New("THERE'S NO EVENT WITH THAT ID🤬🤬🤬🚨🚨")
		}
		repo.Events[event.ID] = *event
		return nil
	}
}

func (repo EventRepo) DeleteEvent(id int, kontek context.Context) error {
	select {
	case <-kontek.Done():
		return kontek.Err()
	default:
		if _, exist := repo.Events[id]; !exist {
			return errors.New("THERE'S NO EVENT WITH THAT ID🤬🚨🤬🚨")
		}
		delete(repo.Events, id)
		return nil
	}
}

func (repo EventRepo) GetAllEvents(kontek context.Context) ([]domain.Event, error) {
	select {
	case <-kontek.Done():
		return nil, kontek.Err()
	default:
		events := make([]domain.Event, 0, len(repo.Events))
		for _, event := range repo.Events {
			events = append(events, event)
		}
		return events, nil
	}
}
