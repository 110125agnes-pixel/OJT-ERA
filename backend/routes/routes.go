package routes

import (
	"github.com/gorilla/mux"
	"backend/controllers"
)

// SetupRoutes configures all application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/api/health", controllers.HealthCheck).Methods("GET")

	// Item routes
	router.HandleFunc("/api/items", controllers.GetItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", controllers.GetItem).Methods("GET")
	router.HandleFunc("/api/items", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", controllers.UpdateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", controllers.DeleteItem).Methods("DELETE")

	return router
}
