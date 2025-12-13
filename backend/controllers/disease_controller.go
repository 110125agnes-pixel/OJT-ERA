package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"backend/models"
)

// GetDiseases handles GET requests to retrieve all diseases
func GetDiseases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	diseases, err := models.GetAllDiseases()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if diseases == nil {
		diseases = []models.Disease{}
	}
	json.NewEncoder(w).Encode(diseases)
}

// CreateDisease handles POST requests to create a new disease
func CreateDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var disease models.Disease
	err := json.NewDecoder(r.Body).Decode(&disease)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if disease.Name == "" || disease.Code == "" || disease.Barcode == "" {
		http.Error(w, "Name, Code, and Barcode are required", http.StatusBadRequest)
		return
	}

	err = models.CreateDisease(&disease)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(disease)
}

// UpdateDisease handles PUT requests to update a disease
func UpdateDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var disease models.Disease
	err = json.NewDecoder(r.Body).Decode(&disease)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if disease.Name == "" || disease.Code == "" || disease.Barcode == "" {
		http.Error(w, "Name, Code, and Barcode are required", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.UpdateDisease(id, &disease)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Disease not found", http.StatusNotFound)
		return
	}

	disease.ID = id
	json.NewEncoder(w).Encode(disease)
}

// DeleteDisease handles DELETE requests to remove a disease
func DeleteDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.DeleteDisease(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Disease not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Disease deleted successfully"})
}

// GetEmployeeDiseases handles GET requests to retrieve diseases for a specific employee
func GetEmployeeDiseases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["employee_id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	diseases, err := models.GetEmployeeDiseases(employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if diseases == nil {
		diseases = []models.EmployeeDisease{}
	}
	json.NewEncoder(w).Encode(diseases)
}

// AddEmployeeDisease handles POST requests to add a disease to an employee
func AddEmployeeDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["employee_id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DiseaseID     int    `json:"disease_id"`
		DateDiagnosed string `json:"date_diagnosed"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.DiseaseID <= 0 {
		http.Error(w, "Valid disease ID is required", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.AddEmployeeDisease(employeeID, req.DiseaseID, req.DateDiagnosed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Failed to add disease to employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Disease added to employee successfully"})
}

// RemoveEmployeeDisease handles DELETE requests to remove a disease from an employee
func RemoveEmployeeDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["employee_id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	diseaseID, err := strconv.Atoi(params["disease_id"])
	if err != nil {
		http.Error(w, "Invalid disease ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.RemoveEmployeeDisease(employeeID, diseaseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Disease not found for this employee", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Disease removed from employee successfully"})
}
