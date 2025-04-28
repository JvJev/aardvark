package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"infocenter-service/handlers"
	"infocenter-service/services"
)

func main() {
	// Create service with 30 second timeout
	infocenterService := services.NewInfocenterService(30 * time.Second)
	
	// Create handlers
	infocenterHandlers := handlers.NewHandlers(infocenterService)
	
	// Define routes
	router := setupRouter()
	
	router.Route("/infocenter", func(routeGroup chi.Router) {
		routeGroup.Get("/{topic}", infocenterHandlers.ReceiveMessages)
		routeGroup.Post("/{topic}", infocenterHandlers.SendMessage)
		routeGroup.Options("/{topic}", func(responseWriter http.ResponseWriter, request *http.Request) {
			responseWriter.WriteHeader(http.StatusOK)
		})
	})
	
	// Start server
	serverPort := ":8080"
	fmt.Printf("Starting Infocenter Service on port %s...\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, router))
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()
	
	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	
	return router
}