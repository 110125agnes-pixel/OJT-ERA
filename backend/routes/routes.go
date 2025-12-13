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

	// Item routes (Employee profiling)
	router.HandleFunc("/api/items", controllers.GetItems).Methods("GET")
	router.HandleFunc("/api/items", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", controllers.UpdateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", controllers.DeleteItem).Methods("DELETE")

	// Surgery routes
	router.HandleFunc("/api/surgeries", controllers.GetAllSurgeriesHandler).Methods("GET")
	router.HandleFunc("/api/surgeries/{id}", controllers.GetSurgeryByIDHandler).Methods("GET")
	router.HandleFunc("/api/surgeries", controllers.CreateSurgeryHandler).Methods("POST")
	router.HandleFunc("/api/surgeries/{id}", controllers.UpdateSurgeryHandler).Methods("PUT")
	router.HandleFunc("/api/surgeries/{id}", controllers.DeleteSurgeryHandler).Methods("DELETE")

	return router
}
