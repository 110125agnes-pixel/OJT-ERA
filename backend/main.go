package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

type Item struct {
	ID          int    `json:"id"`
	Lastname    string `json:"lastname"`
	Firstname   string `json:"firstname"`
	Middlename  string `json:"middlename"`
	Suffix      string `json:"suffix"`
	Birthdate   string `json:"birthdate"`
	Sex         string `json:"sex"`
	CivilStatus string `json:"civil_status"`
}

type InventoryItem struct {
	ID       int     `json:"id"`
	ItemName string  `json:"item_name"`
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	Quantity int     `json:"quantity"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
}

var db *sql.DB

func main() {
	// Database connection
	var err error
	dbPath := getEnv("DB_PATH", "./app.db")

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	log.Println("Successfully connected to SQLite database!")

	// Create table if not exists
	createTable()

	// Setup router
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/items", getItems).Methods("GET")
	router.HandleFunc("/api/items", createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")
	
	// Inventory routes
	router.HandleFunc("/api/inventory", getInventory).Methods("GET")
	router.HandleFunc("/api/inventory", createInventoryItem).Methods("POST")
	router.HandleFunc("/api/inventory/{id}", updateInventoryItem).Methods("PUT")
	router.HandleFunc("/api/inventory/{id}", deleteInventoryItem).Methods("DELETE")
	
	router.HandleFunc("/api/health", healthCheck).Methods("GET")

	// CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func createTable() {
	// Create employee table
	query1 := `
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			lastname TEXT NOT NULL,
			firstname TEXT NOT NULL,
			middlename TEXT,
			suffix TEXT,
			birthdate TEXT,
			sex TEXT,
			civil_status TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query1)
	if err != nil {
		log.Fatal("Error creating items table: ", err)
	}
	log.Println("Table 'items' ready")

	// Create inventory table
	query2 := `
		CREATE TABLE IF NOT EXISTS inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			item_name TEXT NOT NULL,
			category TEXT NOT NULL,
			brand TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 0,
			unit TEXT NOT NULL,
			price REAL NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(query2)
	if err != nil {
		log.Fatal("Error creating inventory table: ", err)
	}
	log.Println("Table 'inventory' ready")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, lastname, firstname, middlename, suffix, birthdate, sex, civil_status FROM items ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Lastname, &item.Firstname, &item.Middlename, &item.Suffix, &item.Birthdate, &item.Sex, &item.CivilStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Lastname == "" || item.Firstname == "" {
		http.Error(w, "Lastname and Firstname are required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO items (lastname, firstname, middlename, suffix, birthdate, sex, civil_status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Sex, item.CivilStatus,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item Item
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Lastname == "" || item.Firstname == "" {
		http.Error(w, "Lastname and Firstname are required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"UPDATE items SET lastname = ?, firstname = ?, middlename = ?, suffix = ?, birthdate = ?, sex = ?, civil_status = ? WHERE id = ?",
		item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Sex, item.CivilStatus, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item.ID = id
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Inventory handlers
func getInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, item_name, category, brand, quantity, unit, price FROM inventory ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []InventoryItem{}
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(&item.ID, &item.ItemName, &item.Category, &item.Brand, &item.Quantity, &item.Unit, &item.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func createInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item InventoryItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.ItemName == "" || item.Category == "" || item.Brand == "" {
		http.Error(w, "Item name, category, and brand are required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO inventory (item_name, category, brand, quantity, unit, price) VALUES (?, ?, ?, ?, ?, ?)",
		item.ItemName, item.Category, item.Brand, item.Quantity, item.Unit, item.Price,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func updateInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item InventoryItem
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.ItemName == "" || item.Category == "" || item.Brand == "" {
		http.Error(w, "Item name, category, and brand are required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"UPDATE inventory SET item_name = ?, category = ?, brand = ?, quantity = ?, unit = ?, price = ? WHERE id = ?",
		item.ItemName, item.Category, item.Brand, item.Quantity, item.Unit, item.Price, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item.ID = id
	json.NewEncoder(w).Encode(item)
}

func deleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM inventory WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Inventory item deleted successfully"})
}
