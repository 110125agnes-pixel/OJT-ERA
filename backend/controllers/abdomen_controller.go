package controllers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func GetAbdomens(w http.ResponseWriter, r *http.Request) {
	abdomens, err := models.GetAllAbdomens(DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(abdomens)
}
