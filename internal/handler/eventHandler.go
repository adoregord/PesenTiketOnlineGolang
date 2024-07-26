package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/usecase"
	"strconv"
	"strings"
	"time"
)

// make a connection to usecase
type EventHandler struct {
	EventUsecase usecase.EventUsecaseInterface
}

func NewEventHandler(eventUsecase usecase.EventUsecaseInterface) EventHandlerInterface {
	return EventHandler{
		EventUsecase: eventUsecase,
	}
}

type EventHandlerInterface interface {
	CreateEvent
	GetEventByID
	GetEventByName
	UpdateEvent
	DeleteEvent
	GetAllEvents
}
type CreateEvent interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
}
type GetEventByID interface {
	GetEventByID(w http.ResponseWriter, r *http.Request)
}
type GetEventByName interface {
	GetEventByName(w http.ResponseWriter, r *http.Request)
}
type UpdateEvent interface {
	UpdateEvent(w http.ResponseWriter, r *http.Request)
}
type DeleteEvent interface {
	DeleteEvent(w http.ResponseWriter, r *http.Request)
}
type GetAllEvents interface {
	GetAllEvents(w http.ResponseWriter, r *http.Request)
}

// function for creating event
func (h EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	// update the kontek to have context timeout in it
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Create Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	var event domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// validate the input
	if err := validate.Struct(event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	events, err := h.EventUsecase.CreateEvent(event, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Create Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusInternalServerError})
		LogMethod("Create Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{Message: "Event has been created", Status: http.StatusOK, Data: events})
	LogMethod("Create Event API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// func for get event by id
func (h EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get Event By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get query param from url
	eventIdStr := r.URL.Query().Get("id")
	if eventIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "Missing event ID in uri param", Status: http.StatusBadRequest})
		LogMethod("Get Event By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// convert the query param id to int
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "Invalid event ID", Status: http.StatusBadRequest})
		LogMethod("Get Event By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	events, err := h.EventUsecase.GetEventByID(eventId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Get Event By ID API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: "Event ID not found", Status: http.StatusNotFound})
		LogMethod("Get Event By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}

	// get the data and show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
	LogMethod("Get Event By ID API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// function for getting event by id
func (h EventHandler) GetEventByName(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is get
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get Event By Name API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get id from url param
	eventName := r.URL.Query().Get("name")
	if strings.TrimSpace(eventName) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "Missing event Name", Status: http.StatusBadRequest})
		LogMethod("Get Event By Name API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	events, err := h.EventUsecase.GetEventByName(eventName, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Get Event By Name API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: "Can't find Event By Name", Status: http.StatusNotFound})
		LogMethod("Get Event By Name API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
	LogMethod("Get Event By Name API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// function for updating event
func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is put
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Update Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Update Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// validate the input
	if err := validate.Struct(event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Update Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send it to usecase
	if err := h.EventUsecase.UpdateEvent(event, kontek); err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Update Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		LogMethod("Update Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{Message: "Event has been updated", Status: http.StatusOK, Data: event})
	LogMethod("Update Event API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// function for deleting event
func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is delete
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Delete Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get id from url param
	eventIdStr := r.URL.Query().Get("id")
	if eventIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "ID param is required", Status: http.StatusBadRequest})
		LogMethod("Delete Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// convert the query param id to int
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Delete Event API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}
	// send to usecase
	if err := h.EventUsecase.DeleteEvent(eventId, kontek); err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Delete Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		LogMethod("Delete Event API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{Message: "Event has been deleted", Status: http.StatusOK})
	LogMethod("Delete Event API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// function for get all events
func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is using get
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get All Events API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// send to usecase
	events, err := h.EventUsecase.GetAllEvents(kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Get All Events API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusNotFound})
		LogMethod("Get All Events API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
	LogMethod("Get All Events API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}
