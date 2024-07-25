package main

import (
	"fmt"
	"net/http"
	"pemesananTiketOnlineGo/internal/domain"
	"pemesananTiketOnlineGo/internal/handler"
	"pemesananTiketOnlineGo/internal/repository"
	"pemesananTiketOnlineGo/internal/usecase"
)

func main() {
	// event connection
	eventRepo := repository.NewEventRepo()
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	eventHandler := handler.NewEventHandler(eventUsecase)

	events := []domain.Event{
		{ID: 1, Name: "Concert1", Date: "02-Jan-2006 15:04:05", Description: "Awokwok1", Location: "Location1"},
	}

	for _, value := range events {
		eventUsecase.CreateEvent(value)
	}

	routes := http.NewServeMux()
	routes.HandleFunc("/event", eventHandler.CreateEvent)
	routes.HandleFunc("/eventGet", eventHandler.GetAllEvents)
	routes.HandleFunc("/eventGetById", eventHandler.GetEventByID)
	routes.HandleFunc("/eventGetByName", eventHandler.GetEventByName)
	routes.HandleFunc("/eventUpdate", eventHandler.UpdateEvent)
	routes.HandleFunc("/eventDelete", eventHandler.DeleteEvent)

	server := http.Server{}
	server.Handler = routes
	server.Addr = ":8080"

	fmt.Println("Server berjalan di http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}