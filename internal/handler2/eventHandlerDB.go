package handler2

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/usecase2"
	"pemesananTiketOnlineGo/internal/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// make a connection to usecase
type EventHandler struct {
	EventUsecase usecase2.EventUsecaseInterface
}

func NewEventHandler(eventUsecase usecase2.EventUsecaseInterface) EventHandlerInterface {
	return EventHandler{
		EventUsecase: eventUsecase,
	}
}

type EventHandlerInterface interface {
	CreateEventDB
	ViewAllEvents
	ViewEventByIdDB
	DeleteEventDB
	UpdateEventTicketDB
}

type CreateEventDB interface {
	CreateEventDB(c *gin.Context)
}
type ViewAllEvents interface {
	ViewAllEvents(c *gin.Context)
}
type ViewEventByIdDB interface {
	ViewEventByIdDB(c *gin.Context)
}
type DeleteEventDB interface {
	DeleteEventDB(c *gin.Context)
}
type UpdateEventTicketDB interface {
	UpdateEventTicketDB(c *gin.Context)
}

func (h EventHandler) ViewEventByIdDB(c *gin.Context) {
	kontek := context.WithValue(c.Request.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)

	var logError error
	var logMessage string
	var logStatus int

	defer func() {
		cancel()
		if logError != nil {
			util.LogFailed(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus, logError)
		} else {
			util.LogSuccess(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus)
		}
	}()

	c.Writer.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Get Event By ID API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	eventIdStr := c.Request.URL.Query().Get("id")
	if eventIdStr == "" {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Missing event ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing event ID in uri param")
		logMessage = "Get Event By ID API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id to int
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Invalid event ID", Status: http.StatusBadRequest})
		logError = err
		logMessage = "Get Event By ID API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	events, err := h.EventUsecase.ViewEventByIdDB(eventId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Get Event By ID API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.JSON(http.StatusNotFound, domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Get Event By ID API Failed"
		logStatus = http.StatusNotFound
		return
	}
	// get the data and show it on response body
	c.JSON(http.StatusOK, domain.Response{Message: "Success", Status: http.StatusOK, Data: events})
	logMessage = "Get Event By ID API Success"
	logStatus = http.StatusOK
}

func (h EventHandler) ViewAllEvents(c *gin.Context) {
	kontek := context.WithValue(c.Request.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)

	var logError error
	var logMessage string
	var logStatus int

	defer func() {
		if logError != nil {
			util.LogFailed(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus, logError)
		}
		util.LogSuccess(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus)
		cancel()
	}()

	c.Writer.Header().Set("Content-Type", "application/json")

	// check if the method is using get
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, domain.Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Get All Events API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// send to usecase
	events, err := h.EventUsecase.ViewAllEvents(kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Get All Events API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.JSON(http.StatusNotFound, domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Get All Events API Failed"
		logStatus = http.StatusNotFound
		return
	}
	// show it on response body
	c.JSON(http.StatusOK, domain.Response{Message: "Success", Status: http.StatusOK, Data: events})
	logMessage = "Get All Events API Success"
	logStatus = http.StatusOK
}

// function for creating event
func (h EventHandler) CreateEventDB(c *gin.Context) {
	kontek := context.WithValue(c.Request.Context(), domain.Key("waktu"), time.Now())
	// update the kontek to have context timeout in it
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)

	var logError error
	var logMessage string
	var logStatus int

	defer func() {
		cancel()
		if logError != nil {
			util.LogFailed(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus, logError)
		} else {
			util.LogSuccess(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus)
		}
	}()

	c.Writer.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Create Event API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	var event domain.Event

	if err := json.NewDecoder(c.Request.Body).Decode(&event); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create Event API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// validate the input
	if err := validate.Struct(event); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create Event API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	err := h.EventUsecase.CreateEventDB(&event, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Create Event API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{Message: err.Error(), Status: http.StatusInternalServerError})
		logError = err
		logMessage = "Create Event API Failed"
		logStatus = http.StatusInternalServerError
		return
	}
	c.JSON(http.StatusOK, domain.Response{Message: "Event has been created", Status: http.StatusOK})
	logMessage = "Create Event API Success"
	logStatus = http.StatusOK
}

// function for deleting event
func (h EventHandler) DeleteEventDB(c *gin.Context) {
	kontek := context.WithValue(c.Request.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)

	var logError error
	var logMessage string
	var logStatus int

	defer func() {
		cancel()
		if logError != nil {
			util.LogFailed(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus, logError)
		} else {
			util.LogSuccess(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus)
		}
	}()

	c.Writer.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if c.Request.Method != "DELETE" {
		c.JSON(http.StatusMethodNotAllowed, domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Delete Event API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	eventIdStr := c.Request.URL.Query().Get("id")
	if eventIdStr == "" {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Missing event ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing event ID in uri param")
		logMessage = "Delete Event API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id to int
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Invalid event ID", Status: http.StatusBadRequest})
		logError = err
		logMessage = "Delete Event API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	err = h.EventUsecase.DeleteEventDB(eventId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Delete Event API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.JSON(http.StatusNotFound, domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Delete Event API Failed"
		logStatus = http.StatusNotFound
		return
	}
	// get the data and show it on response body
	c.JSON(http.StatusOK, domain.Response{Message: "Success deleting event with id: " + eventIdStr, Status: http.StatusOK})
	logMessage = "Delete Event API Success"
	logStatus = http.StatusOK
}

// function add ticket event
func (h EventHandler) UpdateEventTicketDB(c *gin.Context) {
	kontek := context.WithValue(c.Request.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)

	var logError error
	var logMessage string
	var logStatus int

	defer func() {
		cancel()
		if logError != nil {
			util.LogFailed(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus, logError)
		} else {
			util.LogSuccess(logMessage, c.Request.Method, kontek.Value(domain.Key("waktu")).(time.Time), logStatus)
		}
	}()

	c.Writer.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Update Event Ticket API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	eventIdStr := c.Request.URL.Query().Get("id")
	if eventIdStr == "" {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Missing event ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing event ID in uri param")
		logMessage = "Update Event Ticket API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id to int
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: "Invalid event ID", Status: http.StatusBadRequest})
		logError = err
		logMessage = "Update Event Ticket API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	type tickets struct{
		Ticket []domain.Ticket `json:"ticket" validate:"dive"`
	}

	var ticketsInput tickets

	if err := json.NewDecoder(c.Request.Body).Decode(&ticketsInput); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Update Event Ticket API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	err = h.EventUsecase.UpdateEventTicketDB(eventId, &ticketsInput.Ticket, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Update Event Ticket API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.JSON(http.StatusNotFound, domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Update Event Ticket API Failed"
		logStatus = http.StatusNotFound
		return
	}
	// get the data and show it on response body
	c.JSON(http.StatusOK, domain.Response{Message: "Success adding ticket event with id: " + eventIdStr, Status: http.StatusOK})
	logMessage = "Update Event Ticket API Success"
	logStatus = http.StatusOK
}
