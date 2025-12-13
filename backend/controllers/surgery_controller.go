package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"backend/models"
)

// GetAllSurgeriesHandler handles GET /api/surgeries
func GetAllSurgeriesHandler(w http.ResponseWriter, r *http.Request) {
	surgeries, err := models.GetAllSurgeries()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch surgeries: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surgeries)
}

// GetSurgeryByIDHandler handles GET /api/surgeries/:id
func GetSurgeryByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid surgery ID"}`, http.StatusBadRequest)
		return
	}

	surgery, err := models.GetSurgeryByID(id)
	if err != nil {
		http.Error(w, `{"error": "Surgery not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surgery)
}

// CreateSurgeryHandler handles POST /api/surgeries
func CreateSurgeryHandler(w http.ResponseWriter, r *http.Request) {
	var surgery models.Surgery
	if err := json.NewDecoder(r.Body).Decode(&surgery); err != nil {
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if surgery.PatientName == "" || surgery.SurgeryType == "" || surgery.SurgeonName == "" {
		http.Error(w, `{"error": "Patient name, surgery type, and surgeon name are required"}`, http.StatusBadRequest)
		return
	}

	if err := models.CreateSurgery(&surgery); err != nil {
		http.Error(w, `{"error": "Failed to create surgery: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(surgery)
}

// UpdateSurgeryHandler handles PUT /api/surgeries/:id
func UpdateSurgeryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid surgery ID"}`, http.StatusBadRequest)
		return
	}

	var surgery models.Surgery
	if err := json.NewDecoder(r.Body).Decode(&surgery); err != nil {
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if surgery.PatientName == "" || surgery.SurgeryType == "" || surgery.SurgeonName == "" {
		http.Error(w, `{"error": "Patient name, surgery type, and surgeon name are required"}`, http.StatusBadRequest)
		return
	}

	if err := models.UpdateSurgery(id, &surgery); err != nil {
		http.Error(w, `{"error": "Failed to update surgery: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	surgery.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surgery)
}

// DeleteSurgeryHandler handles DELETE /api/surgeries/:id
func DeleteSurgeryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid surgery ID"}`, http.StatusBadRequest)
		return
	}

	if err := models.DeleteSurgery(id); err != nil {
		http.Error(w, `{"error": "Failed to delete surgery: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Surgery deleted successfully"})
}
