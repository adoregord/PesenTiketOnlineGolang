package main

import (
	"context"
	"fmt"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/handler"
	"pemesananTiketOnlineGo/internal/repository"
	"pemesananTiketOnlineGo/internal/usecase"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup

	// event connection
	eventRepo := repository.NewEventRepo()
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	eventHandler := handler.NewEventHandler(eventUsecase)

	// user connection
	userRepo := repository.NewUserRepo()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// order connection
	orderRepo := repository.NewOrderRepo()
	orderUsecase := usecase.NewOrderUsecase(orderRepo, eventRepo, userRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	// create event

	tickets := []domain.Ticket{
		{ID: 1, Type: "VIP", Price: 5000.0, Quantity: 10},
		{ID: 2, Type: "CAT 1", Price: 250.0, Quantity: 100},
	}

	events := []domain.Event{
		{ID: 1, Name: "Concert1", Date: "02-Jan-2006 15:04:05", Description: "Awokwok1", Location: "Location1", Ticket: tickets},
		{ID: 2, Name: "Concert2", Date: "03-Jan-2006 15:04:05", Description: "Awokwok2", Location: "Location2", Ticket: tickets},
		{ID: 3, Name: "Concert3", Date: "04-Jan-2006 15:04:05", Description: "Awokwok3", Location: "Location3", Ticket: tickets},
		{ID: 4, Name: "Concert4", Date: "03-Jan-2006 15:04:05", Description: "Awokwok4", Location: "Location4", Ticket: tickets},
		{ID: 5, Name: "Concert5", Date: "03-Jan-2006 15:04:05", Description: "Awokwok5", Location: "Location5", Ticket: tickets},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, value := range events {
			eventUsecase.CreateEvent(value, context.Background())
		}
	}()

	routes := http.NewServeMux()
	routes.HandleFunc("/event", eventHandler.CreateEvent)
	routes.HandleFunc("/eventGet", eventHandler.GetAllEvents) //check all ticket and event
	routes.HandleFunc("/eventGetById", eventHandler.GetEventByID)
	routes.HandleFunc("/eventGetByName", eventHandler.GetEventByName)
	routes.HandleFunc("/eventUpdate", eventHandler.UpdateEvent)
	routes.HandleFunc("/eventDelete", eventHandler.DeleteEvent)

	routes.HandleFunc("/userPost", userHandler.CreateUser)
	routes.HandleFunc("/userGetAll", userHandler.GetAllUsers)
	routes.HandleFunc("/userGetById", userHandler.GetUserByID)
	routes.HandleFunc("/userGetByName", userHandler.GetUserByName)
	routes.HandleFunc("/userUpdate", userHandler.UpdateUser)
	routes.HandleFunc("/userDelete", userHandler.DeleteUser)

	routes.HandleFunc("/buyTicket", orderHandler.CreateOrder) // buy the ticket
	routes.HandleFunc("/orderGetAll", orderHandler.GetAllOrders)
	routes.HandleFunc("/orderGetByUserId", orderHandler.GetOrderByID) // list all orders from that one user

	server := http.Server{}
	server.Handler = routes
	server.Addr = ":8080"

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Server berjalan di http://localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	wg.Wait()
}
