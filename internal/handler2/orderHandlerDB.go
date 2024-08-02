package handler2

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/usecase2"
	"pemesananTiketOnlineGo/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

// make a connection to usecase
type OrderHandler struct {
	OrderUsecase usecase2.OrderUsecaseInterface
}

func NewOrderHandler(orderUsecase usecase2.OrderUsecaseInterface) OrderHandlerInterface {
	return OrderHandler{
		OrderUsecase: orderUsecase,
	}
}

type OrderHandlerInterface interface {
	CreateOrderDB
}

type CreateOrderDB interface {
	CreateOrderDB(c *gin.Context)
}

// function for creating order in db
func (h OrderHandler) CreateOrderDB(c *gin.Context) {
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
		logMessage = "Create Order API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// var Order domain.Order
	//(userID int, eventID int, ticket []domain.Ticket, kontek context.Context)
	var OrderReq domain.OrderRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&OrderReq); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create Order API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// validate the input
	if err := validate.Struct(OrderReq); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create Order API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	Orders, err := h.OrderUsecase.CreateOrderDB(&OrderReq, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.JSON(http.StatusGatewayTimeout, domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Create Order API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		} else if err.Error() == "INSUFFICIENT BALANCE🤬🤬🤬🚨🚨" {
			c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
			logError = err
			logMessage = "Create Order API Failed"
			logStatus = http.StatusBadRequest
			return
		} else if err.Error() == "THERE'S NO USER WITH THAT ID🤬🚨🤬🚨" {
			c.JSON(http.StatusNotFound, domain.Response{Message: err.Error() + ", please make an account before buy a ticket", Status: http.StatusNotFound})
			logError = err
			logMessage = "Create Order API Failed"
			logStatus = http.StatusNotFound
			return
		} else if err.Error() == "NOT ENOUGH TICKET STOCK🤬🚨🤬🚨" {
			c.JSON(http.StatusConflict, domain.Response{Message: err.Error(), Status: http.StatusConflict})
			logError = err
			logMessage = "Create Order API Failed"
			logStatus = http.StatusConflict
			return
		} else if err.Error() == "THERE'S NO EVENT WITH THAT ID🤬🚨🤬🚨" {
			c.JSON(http.StatusNotFound, domain.Response{Message: err.Error(), Status: http.StatusNotFound})
			logError = err
			logMessage = "Create Order API Failed"
			logStatus = http.StatusNotFound
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{Message: err.Error(), Status: http.StatusInternalServerError})
		logError = err
		logMessage = "Create Order API Failed"
		logStatus = http.StatusInternalServerError
		return
	}
	c.JSON(http.StatusOK, domain.Response{Message: "Success creating order!", Status: http.StatusOK, Data: Orders})
	logMessage = "Create Order API Success"
	logStatus = http.StatusOK
}
