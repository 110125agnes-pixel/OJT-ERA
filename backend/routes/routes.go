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
	router.HandleFunc("/api/items", controllers.CreateItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", controllers.UpdateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", controllers.DeleteItem).Methods("DELETE")

	// Disease routes
	router.HandleFunc("/api/diseases", controllers.GetDiseases).Methods("GET")
	router.HandleFunc("/api/diseases", controllers.CreateDisease).Methods("POST")
	router.HandleFunc("/api/diseases/{id}", controllers.UpdateDisease).Methods("PUT")
	router.HandleFunc("/api/diseases/{id}", controllers.DeleteDisease).Methods("DELETE")

	// Employee-Disease routes
	router.HandleFunc("/api/employees/{employee_id}/diseases", controllers.GetEmployeeDiseases).Methods("GET")
	router.HandleFunc("/api/employees/{employee_id}/diseases", controllers.AddEmployeeDisease).Methods("POST")
	router.HandleFunc("/api/employees/{employee_id}/diseases/{disease_id}", controllers.RemoveEmployeeDisease).Methods("DELETE")

	return router
}
