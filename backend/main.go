package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// ==================== MODELS ====================

type InventoryItem struct {
	ID       int     `json:"id"`
	ItemName string  `json:"item_name"`
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	Quantity int     `json:"quantity"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
}

type SocialHistory struct {
	ID                    int    `json:"id"`
	PatientID             int    `json:"patient_id"`
	IsPatientSmoker       string `json:"is_patient_smoker"`
	CigarettePacksPerYear int    `json:"cigarette_packs_per_year"`
	IsAlcoholDrinker      string `json:"is_alcohol_drinker"`
	BottlesPerDay         int    `json:"bottles_per_day"`
	IsIllicitDrugUser     string `json:"is_illicit_drug_user"`
	IsSexuallyActive      string `json:"is_sexually_active"`
}

type PertinentPhysicalExam struct {
	ID                int     `json:"id"`
	PatientID         int     `json:"patient_id"`
	SystolicBP        int     `json:"systolic_bp"`
	DiastolicBP       int     `json:"diastolic_bp"`
	HeartRate         int     `json:"heart_rate"`
	RespiratoryRate   int     `json:"respiratory_rate"`
	Temperature       float64 `json:"temperature"`
	Height            float64 `json:"height"`
	Weight            float64 `json:"weight"`
	BMI               float64 `json:"bmi"`
	PZScore           int     `json:"pzscore"`
	LeftEyeVision     string  `json:"left_eye_vision"`
	RightEyeVision    string  `json:"right_eye_vision"`
	LengthPediatric   float64 `json:"length_pediatric"`
	HeadCircumference float64 `json:"head_circumference"`
	SkinfoldThickness float64 `json:"skinfold_thickness"`
	Waist             float64 `json:"waist"`
	Hip               float64 `json:"hip"`
	Limbs             float64 `json:"limbs"`
	ArmCircumference  float64 `json:"arm_circumference"`
}

type MedicalHistoryItem struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	DiseaseCode string `json:"disease_code"`
	DiseaseName string `json:"disease_name"`
	IsChecked   bool   `json:"is_checked"`
}

type FamilyHistoryItem struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	DiseaseCode string `json:"disease_code"`
	DiseaseName string `json:"disease_name"`
	Notes       string `json:"notes"`
	IsChecked   bool   `json:"is_checked"`
}

type SurgicalHistoryItem struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	SurgeryCode string `json:"surgery_code"`
	SurgeryName string `json:"surgery_name"`
	Notes       string `json:"notes"`
	IsChecked   bool   `json:"is_checked"`
}

type ImmunizationItem struct {
	ID               int    `json:"id"`
	PatientID        int    `json:"patient_id"`
	VaccineCode      string `json:"vaccine_code"`
	VaccineName      string `json:"vaccine_name"`
	Category         string `json:"category"`
	IsChecked        bool   `json:"is_checked"`
	OtherDescription string `json:"other_description"`
}

type Patient struct {
	ID            int    `json:"id"`
	CaseNo        string `json:"caseNo"`
	HospitalNo    string `json:"hospitalNo"`
	Lastname      string `json:"lastname"`
	Firstname     string `json:"firstname"`
	Middlename    string `json:"middlename"`
	Suffix        string `json:"suffix"`
	Birthdate     string `json:"birthdate"`
	Age           string `json:"age"`
	Room          string `json:"room"`
	AdmissionDate string `json:"admissionDate"`
	DischargeDate string `json:"dischargeDate"`
	Sex           string `json:"sex"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	Complaint     string `json:"complaint"`
}

var db *sql.DB

// ==================== MAIN ====================

func main() {
	var err error

	// CHANGE THIS TO YOUR DATABASE PASSWORD!
	dsn := "root:root@tcp(localhost:3306)/konsulta?parseTime=true"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error opening database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Error connecting to database: ", err)
	}
	log.Println("✓ Successfully connected to konsulta database!")

	// Ensure the patients table exists on startup
	createPatientsTable()

	// Setup router
	router := mux.NewRouter()

	// Inventory Routes
	router.HandleFunc("/api/inventory", getInventory).Methods("GET")
	router.HandleFunc("/api/inventory", createInventoryItem).Methods("POST")
	router.HandleFunc("/api/inventory/{id}", updateInventoryItem).Methods("PUT")
	router.HandleFunc("/api/inventory/{id}", deleteInventoryItem).Methods("DELETE")

	// Patient Routes
	router.HandleFunc("/api/items", getPatients).Methods("GET")          // Get all
	router.HandleFunc("/api/items/{id}", getPatient).Methods("GET")      // Get single (FIXED THIS)
	router.HandleFunc("/api/items", createPatient).Methods("POST")       // Create
	router.HandleFunc("/api/items/{id}", updatePatient).Methods("PUT")   // Update
	router.HandleFunc("/api/items/{id}", deletePatient).Methods("DELETE")// Delete

	// Medical History Routes
	router.HandleFunc("/api/patients/{patientId}/social-history", getSocialHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/social-history", saveSocialHistory).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/pertinent-physical-exam", getPertinentPhysicalExam).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/pertinent-physical-exam", savePertinentPhysicalExam).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/medical-history", getMedicalHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/medical-history", saveMedicalHistory).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/family-history", getFamilyHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/family-history", saveFamilyHistory).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/surgical-history", getSurgicalHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/surgical-history", saveSurgicalHistory).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/immunization", getImmunization).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/immunization", saveImmunization).Methods("POST")

	// CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ==================== INVENTORY HANDLERS ====================

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
		rows.Scan(&item.ID, &item.ItemName, &item.Category, &item.Brand, &item.Quantity, &item.Unit, &item.Price)
		items = append(items, item)
	}
	json.NewEncoder(w).Encode(items)
}

func createInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO inventory (item_name, category, brand, quantity, unit, price) VALUES (?, ?, ?, ?, ?, ?)",
		item.ItemName, item.Category, item.Brand, item.Quantity, item.Unit, item.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	item.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func updateInventoryItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var item InventoryItem
	json.NewDecoder(r.Body).Decode(&item)
	_, err := db.Exec("UPDATE inventory SET item_name=?, category=?, brand=?, quantity=?, unit=?, price=? WHERE id=?",
		item.ItemName, item.Category, item.Brand, item.Quantity, item.Unit, item.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.ID = id
	json.NewEncoder(w).Encode(item)
}

func deleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db.Exec("DELETE FROM inventory WHERE id = ?", id)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}

// ==================== SOCIAL HISTORY HANDLERS ====================

func getSocialHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var sh SocialHistory
	err := db.QueryRow(`SELECT id, patient_id, is_patient_smoker, cigarette_packs_per_year, 
                        is_alcohol_drinker, bottles_per_day, is_illicit_drug_user, is_sexually_active 
                        FROM tsekap_tbl_prof_sochist WHERE patient_id = ?`, patientID).Scan(
		&sh.ID, &sh.PatientID, &sh.IsPatientSmoker, &sh.CigarettePacksPerYear,
		&sh.IsAlcoholDrinker, &sh.BottlesPerDay, &sh.IsIllicitDrugUser, &sh.IsSexuallyActive)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(sh)
}

func saveSocialHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var sh SocialHistory
	json.NewDecoder(r.Body).Decode(&sh)

	var exists int
	db.QueryRow("SELECT COUNT(*) FROM tsekap_tbl_prof_sochist WHERE patient_id = ?", patientID).Scan(&exists)

	if exists > 0 {
		db.Exec(`UPDATE tsekap_tbl_prof_sochist SET is_patient_smoker=?, cigarette_packs_per_year=?, 
                 is_alcohol_drinker=?, bottles_per_day=?, is_illicit_drug_user=?, is_sexually_active=? 
                 WHERE patient_id=?`, sh.IsPatientSmoker, sh.CigarettePacksPerYear, sh.IsAlcoholDrinker,
			sh.BottlesPerDay, sh.IsIllicitDrugUser, sh.IsSexuallyActive, patientID)
	} else {
		db.Exec(`INSERT INTO tsekap_tbl_prof_sochist (patient_id, is_patient_smoker, cigarette_packs_per_year, 
                 is_alcohol_drinker, bottles_per_day, is_illicit_drug_user, is_sexually_active) 
                 VALUES (?, ?, ?, ?, ?, ?, ?)`, patientID, sh.IsPatientSmoker, sh.CigarettePacksPerYear,
			sh.IsAlcoholDrinker, sh.BottlesPerDay, sh.IsIllicitDrugUser, sh.IsSexuallyActive)
	}
	json.NewEncoder(w).Encode(sh)
}

// ==================== PHYSICAL EXAM HANDLERS ====================

func getPertinentPhysicalExam(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var ppe PertinentPhysicalExam
	err := db.QueryRow(`SELECT id, patient_id, systolic_bp, diastolic_bp, heart_rate, respiratory_rate,
                        temperature, height, weight, bmi, pzscore, left_eye_vision, right_eye_vision,
                        length_pediatric, head_circumference, skinfold_thickness, waist, hip, limbs, arm_circumference
                        FROM tsekap_tbl_prof_pespecific WHERE patient_id = ?`, patientID).Scan(
		&ppe.ID, &ppe.PatientID, &ppe.SystolicBP, &ppe.DiastolicBP, &ppe.HeartRate, &ppe.RespiratoryRate,
		&ppe.Temperature, &ppe.Height, &ppe.Weight, &ppe.BMI, &ppe.PZScore, &ppe.LeftEyeVision,
		&ppe.RightEyeVision, &ppe.LengthPediatric, &ppe.HeadCircumference, &ppe.SkinfoldThickness,
		&ppe.Waist, &ppe.Hip, &ppe.Limbs, &ppe.ArmCircumference)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(ppe)
}

func savePertinentPhysicalExam(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var ppe PertinentPhysicalExam
	json.NewDecoder(r.Body).Decode(&ppe)

	var exists int
	db.QueryRow("SELECT COUNT(*) FROM tsekap_tbl_prof_pespecific WHERE patient_id = ?", patientID).Scan(&exists)

	if exists > 0 {
		db.Exec(`UPDATE tsekap_tbl_prof_pespecific SET systolic_bp=?, diastolic_bp=?, heart_rate=?, 
                 respiratory_rate=?, temperature=?, height=?, weight=?, bmi=?, pzscore=?, left_eye_vision=?, 
                 right_eye_vision=?, length_pediatric=?, head_circumference=?, skinfold_thickness=?, 
                 waist=?, hip=?, limbs=?, arm_circumference=? WHERE patient_id=?`,
			ppe.SystolicBP, ppe.DiastolicBP, ppe.HeartRate, ppe.RespiratoryRate, ppe.Temperature,
			ppe.Height, ppe.Weight, ppe.BMI, ppe.PZScore, ppe.LeftEyeVision, ppe.RightEyeVision,
			ppe.LengthPediatric, ppe.HeadCircumference, ppe.SkinfoldThickness, ppe.Waist, ppe.Hip,
			ppe.Limbs, ppe.ArmCircumference, patientID)
	} else {
		db.Exec(`INSERT INTO tsekap_tbl_prof_pespecific (patient_id, systolic_bp, diastolic_bp, heart_rate, 
                 respiratory_rate, temperature, height, weight, bmi, pzscore, left_eye_vision, 
                 right_eye_vision, length_pediatric, head_circumference, skinfold_thickness, 
                 waist, hip, limbs, arm_circumference) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			patientID, ppe.SystolicBP, ppe.DiastolicBP, ppe.HeartRate, ppe.RespiratoryRate, ppe.Temperature,
			ppe.Height, ppe.Weight, ppe.BMI, ppe.PZScore, ppe.LeftEyeVision, ppe.RightEyeVision,
			ppe.LengthPediatric, ppe.HeadCircumference, ppe.SkinfoldThickness, ppe.Waist, ppe.Hip,
			ppe.Limbs, ppe.ArmCircumference)
	}
	json.NewEncoder(w).Encode(ppe)
}

// ==================== CHECKBOX HISTORY HANDLERS ====================

func getMedicalHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	rows, _ := db.Query("SELECT id, patient_id, disease_code, disease_name, is_checked FROM tsekap_tbl_prof_medhist WHERE patient_id = ?", patientID)
	defer rows.Close()
	var list []MedicalHistoryItem
	for rows.Next() {
		var h MedicalHistoryItem
		rows.Scan(&h.ID, &h.PatientID, &h.DiseaseCode, &h.DiseaseName, &h.IsChecked)
		list = append(list, h)
	}
	json.NewEncoder(w).Encode(list)
}

func saveMedicalHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var items []MedicalHistoryItem
	json.NewDecoder(r.Body).Decode(&items)
	db.Exec("DELETE FROM tsekap_tbl_prof_medhist WHERE patient_id = ?", patientID)
	for _, item := range items {
		if item.IsChecked {
			db.Exec("INSERT INTO tsekap_tbl_prof_medhist (patient_id, disease_code, disease_name, is_checked) VALUES (?, ?, ?, ?)",
				patientID, item.DiseaseCode, item.DiseaseName, true)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func getFamilyHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	rows, _ := db.Query("SELECT id, patient_id, disease_code, disease_name, notes, is_checked FROM tsekap_tbl_prof_famhist WHERE patient_id = ?", patientID)
	defer rows.Close()
	var list []FamilyHistoryItem
	for rows.Next() {
		var h FamilyHistoryItem
		rows.Scan(&h.ID, &h.PatientID, &h.DiseaseCode, &h.DiseaseName, &h.Notes, &h.IsChecked)
		list = append(list, h)
	}
	json.NewEncoder(w).Encode(list)
}

func saveFamilyHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var items []FamilyHistoryItem
	json.NewDecoder(r.Body).Decode(&items)
	db.Exec("DELETE FROM tsekap_tbl_prof_famhist WHERE patient_id = ?", patientID)
	for _, item := range items {
		if item.IsChecked {
			db.Exec("INSERT INTO tsekap_tbl_prof_famhist (patient_id, disease_code, disease_name, notes, is_checked) VALUES (?, ?, ?, ?, ?)",
				patientID, item.DiseaseCode, item.DiseaseName, item.Notes, true)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func getSurgicalHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	rows, _ := db.Query("SELECT id, patient_id, surgery_code, surgery_name, notes, is_checked FROM tsekap_tbl_prof_surghist WHERE patient_id = ?", patientID)
	defer rows.Close()
	var list []SurgicalHistoryItem
	for rows.Next() {
		var h SurgicalHistoryItem
		rows.Scan(&h.ID, &h.PatientID, &h.SurgeryCode, &h.SurgeryName, &h.Notes, &h.IsChecked)
		list = append(list, h)
	}
	json.NewEncoder(w).Encode(list)
}

func saveSurgicalHistory(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var items []SurgicalHistoryItem
	json.NewDecoder(r.Body).Decode(&items)
	db.Exec("DELETE FROM tsekap_tbl_prof_surghist WHERE patient_id = ?", patientID)
	for _, item := range items {
		if item.IsChecked {
			db.Exec("INSERT INTO tsekap_tbl_prof_surghist (patient_id, surgery_code, surgery_name, notes, is_checked) VALUES (?, ?, ?, ?, ?)",
				patientID, item.SurgeryCode, item.SurgeryName, item.Notes, true)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func getImmunization(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	rows, _ := db.Query("SELECT id, patient_id, vaccine_code, vaccine_name, category, is_checked, other_description FROM tsekap_tbl_prof_immunization WHERE patient_id = ?", patientID)
	defer rows.Close()
	var list []ImmunizationItem
	for rows.Next() {
		var i ImmunizationItem
		rows.Scan(&i.ID, &i.PatientID, &i.VaccineCode, &i.VaccineName, &i.Category, &i.IsChecked, &i.OtherDescription)
		list = append(list, i)
	}
	json.NewEncoder(w).Encode(list)
}

func saveImmunization(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var items []ImmunizationItem
	json.NewDecoder(r.Body).Decode(&items)
	db.Exec("DELETE FROM tsekap_tbl_prof_immunization WHERE patient_id = ?", patientID)
	for _, item := range items {
		if item.IsChecked {
			db.Exec("INSERT INTO tsekap_tbl_prof_immunization (patient_id, vaccine_code, vaccine_name, category, is_checked, other_description) VALUES (?, ?, ?, ?, ?, ?)",
				patientID, item.VaccineCode, item.VaccineName, item.Category, true, item.OtherDescription)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// ==================== PATIENT HANDLERS ====================

// THIS FUNCTION WAS MISSING, CAUSING THE 405 ERROR
func getPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var p Patient
	err := db.QueryRow(`SELECT id, case_no, hospital_no, lastname, firstname, middlename, 
        suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint 
        FROM patients WHERE id = ?`, id).Scan(
		&p.ID, &p.CaseNo, &p.HospitalNo, &p.Lastname, &p.Firstname, &p.Middlename, &p.Suffix,
		&p.Birthdate, &p.Age, &p.Room, &p.AdmissionDate, &p.DischargeDate, &p.Sex, &p.Height, &p.Weight, &p.Complaint)

	if err == sql.ErrNoRows {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error fetching patient:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func getPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT id, case_no, hospital_no, lastname, firstname, middlename, 
        suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint 
        FROM patients ORDER BY id DESC`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	patients := []Patient{}
	for rows.Next() {
		var p Patient
		rows.Scan(&p.ID, &p.CaseNo, &p.HospitalNo, &p.Lastname, &p.Firstname, &p.Middlename, &p.Suffix,
			&p.Birthdate, &p.Age, &p.Room, &p.AdmissionDate, &p.DischargeDate, &p.Sex, &p.Height, &p.Weight, &p.Complaint)
		patients = append(patients, p)
	}
	json.NewEncoder(w).Encode(patients)
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	var p Patient
	json.NewDecoder(r.Body).Decode(&p)
	result, err := db.Exec(`INSERT INTO patients 
        (case_no, hospital_no, lastname, firstname, middlename, suffix, birthdate, age, 
         room, admission_date, discharge_date, sex, height, weight, complaint) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.CaseNo, p.HospitalNo, p.Lastname, p.Firstname, p.Middlename, p.Suffix,
		p.Birthdate, p.Age, p.Room, p.AdmissionDate, p.DischargeDate, p.Sex,
		p.Height, p.Weight, p.Complaint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	json.NewEncoder(w).Encode(p)
}

func updatePatient(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p Patient
	json.NewDecoder(r.Body).Decode(&p)
	db.Exec(`UPDATE patients SET case_no=?, hospital_no=?, lastname=?, firstname=?, middlename=?, suffix=?, 
        birthdate=?, age=?, room=?, admission_date=?, discharge_date=?, sex=?, height=?, weight=?, complaint=? WHERE id=?`,
		p.CaseNo, p.HospitalNo, p.Lastname, p.Firstname, p.Middlename, p.Suffix,
		p.Birthdate, p.Age, p.Room, p.AdmissionDate, p.DischargeDate, p.Sex,
		p.Height, p.Weight, p.Complaint, id)
	json.NewEncoder(w).Encode(p)
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db.Exec("DELETE FROM patients WHERE id = ?", id)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted"})
}

func createPatientsTable() {
	query := `CREATE TABLE IF NOT EXISTS patients (
        id INT AUTO_INCREMENT PRIMARY KEY,
        case_no VARCHAR(50), hospital_no VARCHAR(50), lastname VARCHAR(100) NOT NULL,
        firstname VARCHAR(100) NOT NULL, middlename VARCHAR(100), suffix VARCHAR(20),
        birthdate DATE, age VARCHAR(10), room VARCHAR(50), admission_date DATETIME,
        discharge_date DATETIME, sex VARCHAR(20), height VARCHAR(20), weight VARCHAR(20),
        complaint TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Warning: %v", err)
	}
}