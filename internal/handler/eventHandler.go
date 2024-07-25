package handler

import (
	"encoding/json"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/usecase"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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

type Response struct {
	Message string `json:"message"`
	Status  any    `json:"status,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (h EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var event domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Bad Request")
		return
	}

	// validate the input

	if err := validate.Struct(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Bad Request")
		return
	}
	if err := h.EventUsecase.CreateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Info().
			Int("httpStatusCode", http.StatusInternalServerError).
			Str("httpMethod", r.Method).
			Msg("Internal Server Error")
		return
	}
	events, _ := h.EventUsecase.GetEventByName(event.Name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "Event has been created", Status: http.StatusOK, Data: events})
}

func (h EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// get id from url param
	eventIdStr := r.URL.Query().Get("id")
	if eventIdStr == "" {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Missing event ID")
		return
	}

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Invalid event ID")
		return
	}
	events, err := h.EventUsecase.GetEventByID(eventId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Info().
			Int("httpStatusCode", http.StatusInternalServerError).
			Str("httpMethod", r.Method).
			Msg("Internal server error")
		return
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(http.Response{Msg: err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h EventHandler) GetEventByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// get id from url param
	eventName := r.URL.Query().Get("name")
	if strings.TrimSpace(eventName) == "" {
		http.Error(w, "Missing Event Name", http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Missing event ID")
		return
	}

	events, err := h.EventUsecase.GetEventByName(eventName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		log.Info().
			Int("httpStatusCode", http.StatusNotFound).
			Str("httpMethod", r.Method).
			Msg("Internal server error")
		return
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(http.Response{Msg: err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var event domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Bad Request")
		return
	}

	// validate the input

	if err := validate.Struct(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}
	if err := h.EventUsecase.UpdateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Info().
			Int("httpStatusCode", http.StatusNotFound).
			Str("httpMethod", r.Method).
			Msg("Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// get id from url param
	eventIdStr := r.URL.Query().Get("id")
	if eventIdStr == "" {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Missing event ID")
		return
	}

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().
			Int("httpStatusCode", http.StatusBadRequest).
			Str("httpMethod", r.Method).
			Msg("Invalid event ID")
		return
	} else if err := h.EventUsecase.DeleteEvent(eventId); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		log.Info().
			Int("httpStatusCode", http.StatusNotFound).
			Str("httpMethod", r.Method).
			Msg("Internal server error")
		return
		// json.NewEncoder(w).Encode(http.Response{Msg: err.Error()})
	}
	w.WriteHeader(http.StatusOK)
}

func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Info().
			Int("httpStatusCode", http.StatusMethodNotAllowed).
			Str("httpMethod", r.Method).
			Msg("Method not allowed")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	events, err := h.EventUsecase.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		log.Info().
			Int("httpStatusCode", http.StatusNotFound).
			Str("httpMethod", r.Method).
			Msg("Internal server error")
		return
		// json.NewEncoder(w).Encode(http.Response{Msg: err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
