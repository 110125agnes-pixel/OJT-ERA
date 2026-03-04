package routes

import (
	"backend/controllers"

	"github.com/gorilla/mux"
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

	// Abdomen routes
	router.HandleFunc("/api/abdomens", controllers.GetAbdomens).Methods("GET")

	return router
}
