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
type UserHandler struct {
	UserUsecase usecase2.UserUsecaseInterface
}

func NewUserHandler(userUsecase usecase2.UserUsecaseInterface) UserHandlerInterface {
	return UserHandler{
		UserUsecase: userUsecase,
	}
}

type UserHandlerInterface interface {
	CreateUserDB
	GetAllUserDB
	GetUserIDDB
	AddBalanceDB
	DeleteUserDB
}
type CreateUserDB interface {
	CreateUserDB(c *gin.Context)
}
type GetAllUserDB interface {
	GetAllUserDB(c *gin.Context)
}
type GetUserIDDB interface {
	GetUserIDDB(c *gin.Context)
}
type AddBalanceDB interface {
	AddBalanceDB(c *gin.Context)
}
type DeleteUserDB interface {
	DeleteUserDB(c *gin.Context)
}

// function for creating User
func (h UserHandler) CreateUserDB(c *gin.Context) {
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
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Create User API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	var User domain.User

	if err := json.NewDecoder(c.Request.Body).Decode(&User); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create User API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// validate the input
	if err := validate.Struct(User); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Create User API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	user, err := h.UserUsecase.CreateUserDB(&User, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Create User API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusInternalServerError})
		logError = err
		logMessage = "Create User API Failed"
		logStatus = http.StatusInternalServerError
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(domain.Response{Message: "User has been created", Status: http.StatusOK, Data: user})
	logMessage = "Create User API Success"
	logStatus = http.StatusOK
}

// func for get User by id
func (h UserHandler) GetUserIDDB(c *gin.Context) {
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
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Get User By ID API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	UserIdStr := c.Request.URL.Query().Get("id")
	if UserIdStr == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Missing User ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing sser ID in uri param")
		logMessage = "Get User By ID API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id to int
	UserId, err := strconv.Atoi(UserIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Get User By ID API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	user, err := h.UserUsecase.GetUserIDDB(UserId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Get User By ID API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Get User By ID API Failed"
		logStatus = http.StatusNotFound
		return
	}
	c.JSON(http.StatusOK, domain.Response{Message: "Success", Status: http.StatusOK, Data: user})
	logMessage = "Get User By ID API Success"
	logStatus = http.StatusOK
}

// function for get all Users
func (h UserHandler) GetAllUserDB(c *gin.Context) {
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

	// check if the method is using get
	if c.Request.Method != "GET" {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Method Not Allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Get All Users API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// send to usecase
	users, err := h.UserUsecase.GetAllUserDB(kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Get All Users API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Get All Users API Failed"
		logStatus = http.StatusNotFound
		return
	}
	// show it on response body
	c.JSON(http.StatusOK, domain.Response{Message: "Success", Status: http.StatusOK, Data: users})
	logMessage = "Get All Users API Success"
	logStatus = http.StatusOK
}

// function for add balance to user
func (h UserHandler) AddBalanceDB(c *gin.Context) {
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
	if c.Request.Method != "PATCH" {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	UserIdStr := c.Request.URL.Query().Get("id")
	TotalStr := c.Request.URL.Query().Get("total")
	if UserIdStr == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Missing User ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing user ID in uri param")
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusBadRequest
		return
	}
	if TotalStr == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Missing total money in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing total money in uri param")
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id and total to int and float
	UserId, err := strconv.Atoi(UserIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusBadRequest
		return
	}
	Total, err := strconv.ParseFloat(TotalStr, 64)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	user, err := h.UserUsecase.AddBalanceDB(UserId, Total, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Add Balance API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Add Balance API Failed"
		logStatus = http.StatusNotFound
		return
	}
	c.JSON(http.StatusOK, domain.Response{Message: "Success", Status: http.StatusOK, Data: user})
	logMessage = "Add Balance API Success"
	logStatus = http.StatusOK
}

// function for delete user
func (h UserHandler) DeleteUserDB(c *gin.Context) {
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
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Method not allowed", Status: http.StatusMethodNotAllowed})
		logError = errors.New("method not allowed")
		logMessage = "Delete User API Failed"
		logStatus = http.StatusMethodNotAllowed
		return
	}

	// get query param from url
	UserIdStr := c.Request.URL.Query().Get("id")
	if UserIdStr == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: "Missing User ID in uri param", Status: http.StatusBadRequest})
		logError = errors.New("missing user ID in uri param")
		logMessage = "Delete User API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// convert the query param id and total to int and float
	UserId, err := strconv.Atoi(UserIdStr)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusBadRequest})
		logError = err
		logMessage = "Delete User API Failed"
		logStatus = http.StatusBadRequest
		return
	}

	// send the data to usecase
	err = h.UserUsecase.DeleteUserDB(UserId, kontek)
	if err != nil {
		if err.Error() == "context deadline exceeded" {
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusGatewayTimeout})
			logError = err
			logMessage = "Add Balance API Failed"
			logStatus = http.StatusGatewayTimeout
			return
		}
		c.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(c.Writer).Encode(domain.Response{Message: err.Error(), Status: http.StatusNotFound})
		logError = err
		logMessage = "Delete User API Failed"
		logStatus = http.StatusNotFound
		return
	}
	c.JSON(http.StatusOK, domain.Response{Message: "Success deleting user with id: " + UserIdStr, Status: http.StatusOK})
	logMessage = "Delete User API Success"
	logStatus = http.StatusOK
}
