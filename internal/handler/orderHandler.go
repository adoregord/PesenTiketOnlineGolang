package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/usecase"
	"strconv"
	"time"
)

// make a connection to usecase
type OrderHandler struct {
	OrderUsecase usecase.OrderUsecaseInterface
}

func NewOrderHandler(orderUsecase usecase.OrderUsecaseInterface) OrderHandlerInterface {
	return OrderHandler{
		OrderUsecase: orderUsecase,
	}
}

type OrderHandlerInterface interface {
	CreateOrder
	GetOrderByID
	GetAllOrders
}
type CreateOrder interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
}
type GetOrderByID interface {
	GetOrderByID(w http.ResponseWriter, r *http.Request)
}
type GetAllOrders interface {
	GetAllOrders(w http.ResponseWriter, r *http.Request)
}

// function for creating Order
func (h OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	// update the kontek to have context timeout in it
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Create Order API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// var Order domain.Order
	//(userID int, eventID int, ticket []domain.Ticket, kontek context.Context)
	var OrderReq domain.OrderRequest

	if err := json.NewDecoder(r.Body).Decode(&OrderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// validate the input
	if err := validate.Struct(OrderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	Orders, err := h.OrderUsecase.CreateOrder(OrderReq, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		} else if err.Error() == "INSUFFICIENT BALANCE" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
			LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		} else if err.Error() == "THERE'S NO USER WITH THAT ID" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error() + ", please make an account before buy a ticket", Status: http.StatusNotFound})
			LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
			return
		} else if err.Error() == "TICKET STOCK NOT ENOUGH" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusConflict})
			LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusConflict)
			return
		}
		json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusInternalServerError})
		LogMethod("Create Order API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{Message: "Order has been created", Status: http.StatusOK, Data: Orders})
	LogMethod("Create Order API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// func for get Order by id
func (h OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get Order By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get query param from url
	OrderIdStr := r.URL.Query().Get("id")
	if OrderIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "Missing Order ID in uri param", Status: http.StatusBadRequest})
		LogMethod("Get Order By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// convert the query param id to int
	OrderId, err := strconv.Atoi(OrderIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(domain.Response{Message: "Invalid Order ID", Status: http.StatusBadRequest})
		LogMethod("Get Order By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	Orders, err := h.OrderUsecase.GetOrderByID(OrderId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Get Order By ID API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		} else if err.Error() == "THIS USER HAVENT BUY A TICKET" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
			LogMethod("Get Order By ID API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: "Order ID not found", Status: http.StatusNotFound})
		LogMethod("Get Order By ID API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}

	// get the data and show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Orders)
	LogMethod("Get Order By ID API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}

// function for get all Orders
func (h OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), domain.Key("waktu"), time.Now())
	kontek, cancel := context.WithTimeout(kontek, 5*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	// check if the method is using get
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get All Orders API Failed", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// send to usecase
	Orders, err := h.OrderUsecase.GetAllOrders(kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			LogMethod("Get All Orders API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusNotFound})
		LogMethod("Get All Orders API Failed "+err.Error(), r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Orders)
	LogMethod("Get All Orders API Success", r.Method, kontek.Value(domain.Key("waktu")).(time.Time), http.StatusOK)
}
