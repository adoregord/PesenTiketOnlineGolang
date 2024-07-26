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
type UserHandler struct {
	UserUsecase usecase.UserUsecaseInterface
}

func NewUserHandler(userUsecase usecase.UserUsecaseInterface) UserHandlerInterface {
	return UserHandler{
		UserUsecase: userUsecase,
	}
}

type UserHandlerInterface interface {
	CreateUser
	GetUserByID
	GetUserByName
	UpdateUser
	DeleteUser
	GetAllUsers
}
type CreateUser interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}
type GetUserByID interface {
	GetUserByID(w http.ResponseWriter, r *http.Request)
}
type GetUserByName interface {
	GetUserByName(w http.ResponseWriter, r *http.Request)
}
type UpdateUser interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
}
type DeleteUser interface {
	DeleteUser(w http.ResponseWriter, r *http.Request)
}
type GetAllUsers interface {
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type Response struct {
	Message string `json:"message"`
	Status  any    `json:"status,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type key string

// function for creating User
func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Create User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	var User domain.User

	if err := json.NewDecoder(r.Body).Decode(&User); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// validate the input
	if err := validate.Struct(User); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Create User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	if err := h.UserUsecase.CreateUser(User); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusInternalServerError})
		LogMethod("Create User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusInternalServerError)
		return
	}
	// get the data and show it on response body
	user, _ := h.UserUsecase.GetUserByName(User.Name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "User has been created", Status: http.StatusOK, Data: user})
	LogMethod("Create User API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}

// func for get User by id
func (h UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is post
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get User By ID API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get query param from url
	UserIdStr := r.URL.Query().Get("id")
	if UserIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Missing User ID in uri param", Status: http.StatusBadRequest})
		LogMethod("Get User By ID API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// convert the query param id to int
	UserId, err := strconv.Atoi(UserIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid User ID", Status: http.StatusBadRequest})
		LogMethod("Get User By ID API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	user, err := h.UserUsecase.GetUserByID(UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Message: "User ID not found", Status: http.StatusNotFound})
		LogMethod("Get User By ID API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusNotFound)
		return
	}

	// get the data and show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	LogMethod("Get User By ID API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}

// function for getting User by id
func (h UserHandler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is get
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get User By Name API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get id from url param
	UserName := r.URL.Query().Get("name")
	if strings.TrimSpace(UserName) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Missing User name", Status: http.StatusBadRequest})
		LogMethod("Get User By Name API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send the data to usecase
	user, err := h.UserUsecase.GetUserByName(UserName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Message: "Can't find User By Name", Status: http.StatusNotFound})
		LogMethod("Get User By Name API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	LogMethod("Get User By Name API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}

// function for updating User
func (h UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is put
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Update User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Update User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// validate the input
	if err := validate.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Update User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// send it to usecase
	if err := h.UserUsecase.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusNotFound})
		LogMethod("Update User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "User has been updated", Status: http.StatusOK, Data: user})
	LogMethod("Update User API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}

// function for deleting User
func (h UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is delete
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Delete User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// get id from url param
	UserIdStr := r.URL.Query().Get("id")
	if UserIdStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "ID param is required", Status: http.StatusBadRequest})
		LogMethod("Delete User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}

	// convert the query param id to int
	UserId, err := strconv.Atoi(UserIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusBadRequest})
		LogMethod("Delete User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusBadRequest)
		return
	}
	// send to usecase
	if err := h.UserUsecase.DeleteUser(UserId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Message: err.Error(), Status: http.StatusNotFound})
		LogMethod("Delete User API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "User has been deleted", Status: http.StatusOK})
	LogMethod("Delete User API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}

// function for get all Users
func (h UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	kontek := context.WithValue(r.Context(), key("waktu"), time.Now())
	w.Header().Set("Content-Type", "application/json")

	// check if the method is using get
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		LogMethod("Get All Users API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusMethodNotAllowed)
		return
	}

	// send to usecase
	users, err := h.UserUsecase.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Message: "Method Not Allowed", Status: http.StatusNotFound})
		LogMethod("Get All Users API Failed", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusNotFound)
		return
	}
	// show it on response body
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	LogMethod("Get All Users API Success", r.Method, kontek.Value(key("waktu")).(time.Time), http.StatusOK)
}