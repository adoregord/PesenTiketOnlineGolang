package routes

import (
	"database/sql"
	"pemesananTiketOnlineGo/internal/handler2"
	"pemesananTiketOnlineGo/internal/repository2"
	"pemesananTiketOnlineGo/internal/usecase2"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(database *sql.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	userRepo2 := repository2.NewUserRepo(database)
	userUsecase2 := usecase2.NewUserUsecase(userRepo2)
	userHandler2 := handler2.NewUserHandler(userUsecase2)

	ticketRepo2 := repository2.NewTicketRepo(database)

	eventRepo2 := repository2.NewEventRepo(database)
	eventUsecase2 := usecase2.NewEventUsecase(eventRepo2, ticketRepo2)
	eventHandler2 := handler2.NewEventHandler(eventUsecase2)

	orderRepo2 := repository2.NewOrderRepo(database)
	orderUsecase2 := usecase2.NewOrderUsecase(orderRepo2, eventRepo2, userRepo2)
	orderHandler2 := handler2.NewOrderHandler(orderUsecase2)

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/get", userHandler2.GetAllUserDB)
		userRoutes.GET("/getByID", userHandler2.GetUserIDDB)
		userRoutes.POST("/create", userHandler2.CreateUserDB)
		userRoutes.PATCH("addBalance", userHandler2.AddBalanceDB)
		userRoutes.DELETE("/delete", userHandler2.DeleteUserDB)
	}

	eventRoutes := router.Group("/event")
	{
		eventRoutes.GET("/get", eventHandler2.ViewAllEvents)
		eventRoutes.GET("/getByID", eventHandler2.ViewEventByIdDB)
		eventRoutes.POST("/create", eventHandler2.CreateEventDB)
		eventRoutes.POST("/addTicket", eventHandler2.UpdateEventTicketDB)
		eventRoutes.DELETE("/delete", eventHandler2.DeleteEventDB)
	}

	orderRoutes := router.Group("/order")
	{
		orderRoutes.GET("/view", orderHandler2.ViewAllOrdersDB)
		orderRoutes.GET("/userView", orderHandler2.ViewUsersOrder)
		orderRoutes.POST("/create", orderHandler2.CreateOrderDB)
	}

	return router
}
