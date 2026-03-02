package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

type LibDisease struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type VaccineLibItem struct {
	VaccineCode string `json:"vaccine_code"`
	VaccineName string `json:"vaccine_name"`
}

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

	// Ensure physical exam tables exist
	createPhysicalExamTables()

	// Ensure patient_medhist summary table exists.
	// This table stores one row per patient with a pipe-separated 0/1 string
	// (e.g. "1|0|1|0|0") representing which diseases were checked in the Medical tab.
	// It is separate from tsekap_tbl_prof_medhist which has FK constraints.
	createMedHistSummaryTable()

	// Setup router
	router := mux.NewRouter()

	// Inventory Routes
	router.HandleFunc("/api/inventory", getInventory).Methods("GET")
	router.HandleFunc("/api/inventory", createInventoryItem).Methods("POST")
	router.HandleFunc("/api/inventory/{id}", updateInventoryItem).Methods("PUT")
	router.HandleFunc("/api/inventory/{id}", deleteInventoryItem).Methods("DELETE")

	// Patient Routes
	router.HandleFunc("/api/items", getPatients).Methods("GET")           // Get all
	router.HandleFunc("/api/items/{id}", getPatient).Methods("GET")       // Get single (FIXED THIS)
	router.HandleFunc("/api/items", createPatient).Methods("POST")        // Create
	router.HandleFunc("/api/items/{id}", updatePatient).Methods("PUT")    // Update
	router.HandleFunc("/api/items/{id}", deletePatient).Methods("DELETE") // Delete

	// Medical History Routes
	router.HandleFunc("/api/patients/{patientId}/social-history", getSocialHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/social-history", saveSocialHistory).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/pertinent-physical-exam", getPertinentPhysicalExam).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/pertinent-physical-exam", savePertinentPhysicalExam).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/medical-history", getMedicalHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/medical-history", saveMedicalHistory).Methods("POST")
	router.HandleFunc("/api/lib/mdiseases", getMedicalDiseases).Methods("GET")

	// Physical Examination Routes (NEW)
	router.HandleFunc("/api/patients/{patientId}/physical-exam/general", getPhysicalExamGeneral).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/physical-exam/general", savePhysicalExamGeneral).Methods("POST")
	router.HandleFunc("/api/patients/{patientId}/physical-exam/findings", getPhysicalExamFindings).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/physical-exam/findings", savePhysicalExamFindings).Methods("POST")
	// Surgical library
	router.HandleFunc("/api/lib/surgery", getSurgicalLib).Methods("GET")
	// Digital rectal library
	router.HandleFunc("/api/lib/digital_rectal", getDigitalRectalLib).Methods("GET")
	router.HandleFunc("/api/lib/digital_rectal", saveDigitalRectalLib).Methods("POST")
	// Genitourinary library
	router.HandleFunc("/api/lib/genitourinary", getGenitourinaryLib).Methods("GET")
	router.HandleFunc("/api/lib/genitourinary", saveGenitourinaryLib).Methods("POST")

	// Immunization library routes
	router.HandleFunc("/api/lib/immchild", getImmChildLib).Methods("GET")
	router.HandleFunc("/api/lib/immyoungw", getImmYoungLib).Methods("GET")
	router.HandleFunc("/api/lib/immpregw", getImmPregLib).Methods("GET")
	router.HandleFunc("/api/lib/immelderly", getImmElderlyLib).Methods("GET")
	// Skin library
	router.HandleFunc("/api/lib/skin", getSkinLib).Methods("GET")
	// Chest library
	router.HandleFunc("/api/lib/chest", getChestLib).Methods("GET")
	router.HandleFunc("/api/lib/chest", saveChestLib).Methods("POST")
	// Abdomen library
	router.HandleFunc("/api/lib/abdomen", getAbdomenLib).Methods("GET")
	router.HandleFunc("/api/lib/abdomen", saveAbdomenLib).Methods("POST")
	// Heart library
	router.HandleFunc("/api/lib/heart", getHeartLib).Methods("GET")
	router.HandleFunc("/api/lib/heart", saveHeartLib).Methods("POST")
	// Neuro library
	router.HandleFunc("/api/lib/neuro", getNeuroLib).Methods("GET")
	router.HandleFunc("/api/lib/neuro", saveNeuroLib).Methods("POST")
	// HEENT library
	router.HandleFunc("/api/lib/heent", getHeentLib).Methods("GET")
	router.HandleFunc("/api/lib/heent", saveHeentLib).Methods("POST")
	// Debug: dump all tables and rows from konsulta database
	router.HandleFunc("/api/debug/dump", dumpDB).Methods("GET")
	// Family Library
	router.HandleFunc("/api/patients/{patientId}/family-history", getFamilyHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/family-history", saveFamilyHistory).Methods("POST")
	// Surgery Library
	router.HandleFunc("/api/patients/{patientId}/surgical-history", getSurgicalHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/surgical-history", saveSurgicalHistory).Methods("POST")
	// Immunization Library
	router.HandleFunc("/api/patients/{patientId}/immunization", getImmunization).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/immunization", saveImmunization).Methods("POST")

	// CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // enable verbose debug logs in dev
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

// getMedicalHistory handles GET /api/patients/{patientId}/medical-history
// It returns the full disease list from tsekap_lib_mdiseases, each with is_checked = true/false
// based on the saved 1|0 string in patient_medhist for this patient.
func getMedicalHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"] // numeric ID from URL e.g. /api/patients/5/medical-history

	// Translate numeric patient ID → case_no (e.g. "C2026-00001") used as patno key
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID // fallback: use numeric ID if case_no not set
	}

	// Load all diseases from the library in sorted order — this defines the bit positions
	// Position 0 = first disease (lowest code), position N = last disease
	libRows, err := db.Query("SELECT mdisease_code, mdisease_desc FROM tsekap_lib_mdiseases ORDER BY mdisease_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer libRows.Close()

	type libEntry struct{ Code, Desc string }
	var lib []libEntry
	for libRows.Next() {
		var e libEntry
		libRows.Scan(&e.Code, &e.Desc)
		lib = append(lib, e)
	}

	// Fetch the saved pipe-separated 0/1 string for this patient
	// Example: "1|0|1|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0" means disease[0] and disease[2] are checked
	var saved string
	db.QueryRow("SELECT mdisease_code FROM patient_medhist WHERE patno = ?", patno).Scan(&saved)

	// Split the saved string into individual bits
	bits := []string{}
	if saved != "" {
		bits = strings.Split(saved, "|")
	}

	// Build the response array — each disease gets its is_checked value from its position in bits[]
	var list []MedicalHistoryItem
	for i, d := range lib {
		isChecked := false
		if i < len(bits) {
			isChecked = bits[i] == "1" // "1" = checked, "0" or missing = unchecked
		}
		list = append(list, MedicalHistoryItem{
			DiseaseCode: d.Code,
			DiseaseName: d.Desc,
			IsChecked:   isChecked,
		})
	}
	json.NewEncoder(w).Encode(list)
}

// saveMedicalHistory handles POST /api/patients/{patientId}/medical-history
// It receives the full disease list with is_checked flags, builds a positional 0/1 string
// ordered by tsekap_lib_mdiseases, and upserts one row into patient_medhist.
func saveMedicalHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"] // numeric ID from URL

	// Decode the JSON body — array of {disease_code, disease_name, is_checked}
	var items []MedicalHistoryItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Translate numeric patient ID → case_no used as the primary key in patient_medhist
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID // fallback to numeric ID
	}

	// Re-fetch the library order from DB — this guarantees the bit positions in the
	// saved string always match the current state of tsekap_lib_mdiseases.
	// If a disease is added/removed, the string will be rebuilt correctly on next save.
	libRows, err := db.Query("SELECT mdisease_code FROM tsekap_lib_mdiseases ORDER BY mdisease_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer libRows.Close()
	var libOrder []string
	for libRows.Next() {
		var code string
		libRows.Scan(&code)
		libOrder = append(libOrder, code)
	}

	// Build a lookup map: disease_code => is_checked from the request body
	checkedMap := map[string]bool{}
	for _, item := range items {
		checkedMap[item.DiseaseCode] = item.IsChecked
	}

	// Build the positional 1|0 string in the same order as the library
	// e.g. if lib has [001,002,003] and 001+003 are checked → "1|0|1"
	bits := make([]string, len(libOrder))
	for i, code := range libOrder {
		if checkedMap[code] {
			bits[i] = "1"
		} else {
			bits[i] = "0"
		}
	}
	mdiseaseCode := strings.Join(bits, "|")

	// UPSERT: INSERT if this patient has no row yet, UPDATE if they do.
	// ON DUPLICATE KEY UPDATE means no need for a separate SELECT COUNT check.
	_, execErr := db.Exec(
		`INSERT INTO patient_medhist (patno, mdisease_code, date_added, added_by)
		 VALUES (?, ?, NOW(), 'system')
		 ON DUPLICATE KEY UPDATE mdisease_code = VALUES(mdisease_code), date_added = NOW()`,
		patno, mdiseaseCode)
	if execErr != nil {
		log.Println("UPSERT error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}

	// Return the saved values so the frontend/SQLyog can verify what was stored
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved", "patno": patno, "mdisease_code": mdiseaseCode})
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

// Library endpoint: list medical disease codes
func getMedicalDiseases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT mdisease_code, mdisease_desc FROM tsekap_lib_mdiseases ORDER BY mdisease_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []LibDisease{}
	for rows.Next() {
		var code, desc string
		rows.Scan(&code, &desc)
		list = append(list, LibDisease{Code: code, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Library endpoint: digital rectal options
func getDigitalRectalLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT RECTAL_ID, RECTAL_DESC FROM tsekap_lib_digital_rectal ORDER BY SORT_NO, RECTAL_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type RectalItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []RectalItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, RectalItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// dumpDB returns all tables and their rows as JSON (for local debugging only)
func dumpDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get table list
	tablesRows, err := db.Query("SHOW TABLES")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tablesRows.Close()

	tables := []string{}
	for tablesRows.Next() {
		var tbl string
		if err := tablesRows.Scan(&tbl); err != nil {
			continue
		}
		tables = append(tables, tbl)
	}

	result := map[string]interface{}{}

	for _, tbl := range tables {
		q := fmt.Sprintf("SELECT * FROM %s", tbl)
		rows, err := db.Query(q)
		if err != nil {
			result[tbl] = map[string]string{"error": err.Error()}
			continue
		}

		cols, err := rows.Columns()
		if err != nil {
			rows.Close()
			result[tbl] = map[string]string{"error": err.Error()}
			continue
		}

		tableData := []map[string]interface{}{}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				continue
			}

			rowMap := map[string]interface{}{}
			for i, colName := range cols {
				val := columns[i]
				// Convert []byte to string for readability
				if b, ok := val.([]byte); ok {
					rowMap[colName] = string(b)
				} else {
					rowMap[colName] = val
				}
			}
			tableData = append(tableData, rowMap)
		}
		rows.Close()
		result[tbl] = tableData
	}

	json.NewEncoder(w).Encode(result)
}

// ==================== IMMUNIZATION LIBRARIES ====================

func getImmChildLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT IMM_CODE AS vaccine_code, IMM_DESC AS vaccine_name FROM tsekap_lib_immchild ORDER BY IMM_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []VaccineLibItem{}
	for rows.Next() {
		var code, name string
		rows.Scan(&code, &name)
		list = append(list, VaccineLibItem{VaccineCode: code, VaccineName: name})
	}
	json.NewEncoder(w).Encode(list)
}

func getImmYoungLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT IMM_CODE AS vaccine_code, IMM_DESC AS vaccine_name FROM tsekap_lib_immyoungw ORDER BY IMM_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []VaccineLibItem{}
	for rows.Next() {
		var code, name string
		rows.Scan(&code, &name)
		list = append(list, VaccineLibItem{VaccineCode: code, VaccineName: name})
	}
	json.NewEncoder(w).Encode(list)
}

func getImmPregLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT IMM_CODE AS vaccine_code, IMM_DESC AS vaccine_name FROM tsekap_lib_immpregw ORDER BY IMM_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []VaccineLibItem{}
	for rows.Next() {
		var code, name string
		rows.Scan(&code, &name)
		list = append(list, VaccineLibItem{VaccineCode: code, VaccineName: name})
	}
	json.NewEncoder(w).Encode(list)
}

func getImmElderlyLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT IMM_CODE AS vaccine_code, IMM_DESC AS vaccine_name FROM tsekap_lib_immelderly ORDER BY IMM_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := []VaccineLibItem{}
	for rows.Next() {
		var code, name string
		rows.Scan(&code, &name)
		list = append(list, VaccineLibItem{VaccineCode: code, VaccineName: name})
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

// generateCaseNo creates a unique patient case number in the format CYYYY-NNNNN.
// Examples: C2026-00001, C2026-00002, C2027-00001 (resets each year)
// Logic:
//  1. Gets the current year from the DB server (not the app server) to avoid clock issues.
//  2. Queries the max existing sequence for that year prefix.
//  3. Returns prefix + (maxSeq + 1) zero-padded to 5 digits.
func generateCaseNo() string {
	// Use DB server time so case numbers are consistent even if backend clock is wrong
	var dbYear int
	if err := db.QueryRow("SELECT YEAR(NOW())").Scan(&dbYear); err != nil || dbYear == 0 {
		dbYear = 2026 // safe default
	}
	prefix := fmt.Sprintf("C%d-", dbYear) // e.g. "C2026-"

	// Find the highest sequence number already used this year (e.g. 3 if C2026-00003 exists)
	var maxSeq int
	db.QueryRow(
		`SELECT COALESCE(MAX(CAST(SUBSTRING(case_no, 7) AS UNSIGNED)), 0)
		 FROM patients WHERE case_no LIKE ?`, prefix+"%").Scan(&maxSeq)

	// Return next sequence zero-padded to 5 digits: C2026-00004
	return fmt.Sprintf("%s%05d", prefix, maxSeq+1)
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	var p Patient
	json.NewDecoder(r.Body).Decode(&p)

	// Always auto-generate the case number on the backend
	p.CaseNo = generateCaseNo()

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
	newID, _ := result.LastInsertId()
	p.ID = int(newID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func updatePatient(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	var p Patient
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(`UPDATE patients SET case_no=?, hospital_no=?, lastname=?, firstname=?, middlename=?, suffix=?, 
			birthdate=?, age=?, room=?, admission_date=?, discharge_date=?, sex=?, height=?, weight=?, complaint=? WHERE id=?`,
		p.CaseNo, p.HospitalNo, p.Lastname, p.Firstname, p.Middlename, p.Suffix,
		p.Birthdate, p.Age, p.Room, p.AdmissionDate, p.DischargeDate, p.Sex,
		p.Height, p.Weight, p.Complaint, idStr)
	if err != nil {
		log.Println("Error updating patient:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Confirm rows affected
	if ra, _ := res.RowsAffected(); ra == 0 {
		http.Error(w, "No rows updated", http.StatusNotFound)
		return
	}

	// Read back the saved row and return it (authoritative source)
	err = db.QueryRow(`SELECT id, case_no, hospital_no, lastname, firstname, middlename, 
			suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint 
			FROM patients WHERE id = ?`, idStr).Scan(
		&p.ID, &p.CaseNo, &p.HospitalNo, &p.Lastname, &p.Firstname, &p.Middlename, &p.Suffix,
		&p.Birthdate, &p.Age, &p.Room, &p.AdmissionDate, &p.DischargeDate, &p.Sex, &p.Height, &p.Weight, &p.Complaint)
	if err != nil {
		log.Println("Error fetching updated patient:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

// Library endpoint: surgical options
func getSurgicalLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Use the actual column names present in the database. Some installs
	// name the columns `surgery_code`, `surgery_name` and `description`.
	rows, err := db.Query("SELECT surgery_code, surgery_name, description FROM tsekap_lib_surgical ORDER BY sort_order, surgery_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SurgItem struct {
		SurgeryCode string `json:"surgery_code"`
		SurgeryName string `json:"surgery_name"`
		Description string `json:"description"`
	}

	list := []SurgItem{}
	for rows.Next() {
		var code, name, desc sql.NullString
		if err := rows.Scan(&code, &name, &desc); err != nil {
			continue
		}
		// Normalize and skip empty/null-like values coming from different DB installs
		codeStr := strings.TrimSpace(code.String)
		nameStr := strings.TrimSpace(name.String)
		descStr := strings.TrimSpace(desc.String)
		if nameStr == "" || strings.EqualFold(nameStr, "(null)") {
			// skip rows without a valid name
			continue
		}
		if codeStr == "" || strings.EqualFold(codeStr, "(null)") {
			// allow missing code but if it's the literal '(NULL)' skip
			if codeStr == "(NULL)" {
				continue
			}
		}
		list = append(list, SurgItem{SurgeryCode: codeStr, SurgeryName: nameStr, Description: descStr})
	}
	json.NewEncoder(w).Encode(list)
}

// createMedHistSummaryTable auto-creates the patient_medhist table on backend startup.
// Schema design:
//
//	patno         - the patient's case_no (e.g. C2026-00001), serves as PRIMARY KEY
//	mdisease_code - pipe-separated 0/1 string, one bit per disease in tsekap_lib_mdiseases order
//	                e.g. "1|0|1|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0|0" means disease 1 and 3 are checked
//	date_added    - timestamp of last save
//	added_by      - who saved (currently always 'system')
//
// Note: no foreign keys, so diseases can be freely added/deleted in tsekap_lib_mdiseases.
func createMedHistSummaryTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_medhist (
		patno        VARCHAR(20)  NOT NULL PRIMARY KEY,
		mdisease_code VARCHAR(500) NOT NULL DEFAULT '',
		date_added   DATETIME,
		added_by     VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_medhist table ready")
}

// ensureMedHistColumns — kept as no-op for compatibility
func ensureMedHistColumns() {}

func createPhysicalExamTables() {
	// General info table (general survey, remarks, blood type)
	generalQuery := `CREATE TABLE IF NOT EXISTS tsekap_tbl_prof_pe_general (
		id INT AUTO_INCREMENT PRIMARY KEY,
		patient_id INT NOT NULL UNIQUE,
		general_survey VARCHAR(50) DEFAULT 'awake',
		remarks TEXT,
		blood_type VARCHAR(10) DEFAULT 'A+',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY (patient_id)
	)`
	if _, err := db.Exec(generalQuery); err != nil {
		log.Printf("Warning createPhysicalExamTables (general): %v", err)
	}

	// Findings table (one row per checked finding per patient)
	findingsQuery := `CREATE TABLE IF NOT EXISTS tsekap_tbl_prof_pe_findings (
		id INT AUTO_INCREMENT PRIMARY KEY,
		patient_id INT NOT NULL,
		category VARCHAR(50) NOT NULL,
		finding_code VARCHAR(100) NOT NULL,
		finding_desc VARCHAR(255) NOT NULL,
		is_checked BOOLEAN DEFAULT TRUE,
		others_text TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		KEY (patient_id),
		UNIQUE KEY unique_finding (patient_id, category, finding_code)
	)`
	if _, err := db.Exec(findingsQuery); err != nil {
		log.Printf("Warning createPhysicalExamTables (findings): %v", err)
	}

	log.Println("✓ Physical exam tables ready")
}

// ==================== PHYSICAL EXAM HANDLERS ====================

func getPhysicalExamGeneral(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var g PhysicalExamGeneral
	err := db.QueryRow(`SELECT id, patient_id, general_survey, remarks, blood_type 
		FROM tsekap_tbl_prof_pe_general WHERE patient_id = ?`, patientID).Scan(
		&g.ID, &g.PatientID, &g.GeneralSurvey, &g.Remarks, &g.BloodType)

	if err == sql.ErrNoRows {
		// Return defaults if nothing saved yet
		json.NewEncoder(w).Encode(PhysicalExamGeneral{
			GeneralSurvey: "awake",
			BloodType:     "A+",
		})
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(g)
}

func savePhysicalExamGeneral(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var g PhysicalExamGeneral
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`INSERT INTO tsekap_tbl_prof_pe_general 
		(patient_id, general_survey, remarks, blood_type)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		general_survey = VALUES(general_survey),
		remarks = VALUES(remarks),
		blood_type = VALUES(blood_type)`,
		patientID, g.GeneralSurvey, g.Remarks, g.BloodType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// FindingPayload for batch save
type FindingPayload struct {
	Category    string `json:"category"`
	FindingCode string `json:"finding_code"`
	FindingDesc string `json:"finding_desc"`
	IsChecked   bool   `json:"is_checked"`
	OthersText  string `json:"others_text"`
}

func getPhysicalExamFindings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	rows, err := db.Query(`SELECT id, patient_id, category, finding_code, finding_desc, is_checked, others_text
		FROM tsekap_tbl_prof_pe_findings WHERE patient_id = ?`, patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type FindingRow struct {
		ID          int    `json:"id"`
		PatientID   int    `json:"patient_id"`
		Category    string `json:"category"`
		FindingCode string `json:"finding_code"`
		FindingDesc string `json:"finding_desc"`
		IsChecked   bool   `json:"is_checked"`
		OthersText  string `json:"others_text"`
	}

	list := []FindingRow{}
	for rows.Next() {
		var f FindingRow
		rows.Scan(&f.ID, &f.PatientID, &f.Category, &f.FindingCode, &f.FindingDesc, &f.IsChecked, &f.OthersText)
		list = append(list, f)
	}
	json.NewEncoder(w).Encode(list)
}

func savePhysicalExamFindings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var payload struct {
		Category string           `json:"category"`
		Findings []FindingPayload `json:"findings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Decode error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use DELETE + INSERT for simplicity and reliability
	_, err := db.Exec(`DELETE FROM tsekap_tbl_prof_pe_findings WHERE patient_id = ? AND category = ?`,
		patientID, payload.Category)
	if err != nil {
		log.Println("Delete error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only insert checked findings (or "others" with text)
	for _, f := range payload.Findings {
		if f.IsChecked {
			_, err := db.Exec(`INSERT INTO tsekap_tbl_prof_pe_findings 
				(patient_id, category, finding_code, finding_desc, is_checked, others_text)
				VALUES (?, ?, ?, ?, ?, ?)`,
				patientID, payload.Category, f.FindingCode, f.FindingDesc, true, f.OthersText)
			if err != nil {
				log.Println("Insert finding error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// PhysicalExamFinding - one row per checked finding
type PhysicalExamFinding struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	Category    string `json:"category"`     // e.g. "skin", "heent", "chest"
	FindingCode string `json:"finding_code"` // e.g. "essentiallyNormal"
	FindingDesc string `json:"finding_desc"` // e.g. "Essentially normal"
	IsChecked   bool   `json:"is_checked"`
}

// PhysicalExamGeneral - stores general survey, remarks, blood type
type PhysicalExamGeneral struct {
	ID            int    `json:"id"`
	PatientID     int    `json:"patient_id"`
	GeneralSurvey string `json:"general_survey"`
	Remarks       string `json:"remarks"`
	BloodType     string `json:"blood_type"`
}

// Library endpoint: skin/extremities options
func getSkinLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT SKIN_ID, SKIN_DESC FROM tsekap_lib_skin_extremities ORDER BY SORT_NO, SKIN_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SkinItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []SkinItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, SkinItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Library endpoint: HEENT options
func getHeentLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT HEENT_ID, HEENT_DESC FROM tsekap_lib_heent ORDER BY SORT_NO, HEENT_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HeentItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []HeentItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, HeentItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert HEENT library rows (accepts array of {id, desc, sort_no})
func saveHeentLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			// try update
			if _, err := tx.Exec(`UPDATE tsekap_lib_heent SET HEENT_DESC = ?, SORT_NO = ? WHERE HEENT_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_heent (HEENT_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Library endpoint: Chest options
func getChestLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT CHEST_ID, CHEST_DESC FROM tsekap_lib_chest ORDER BY SORT_NO, CHEST_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ChestItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []ChestItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, ChestItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert Chest library rows (accepts array of {id, desc, sort_no})
func saveChestLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_chest SET CHEST_DESC = ?, SORT_NO = ? WHERE CHEST_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_chest (CHEST_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Library endpoint: Heart options
func getHeartLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT HEART_ID, HEART_DESC FROM tsekap_lib_heart ORDER BY SORT_NO, HEART_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HeartItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []HeartItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, HeartItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert Heart library rows (accepts array of {id, desc, sort_no})
func saveHeartLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_heart SET HEART_DESC = ?, SORT_NO = ? WHERE HEART_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_heart (HEART_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Library endpoint: Abdomen options
func getAbdomenLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT ABDOMEN_ID, ABDOMEN_DESC FROM tsekap_lib_abdomen ORDER BY SORT_NO, ABDOMEN_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AbdomenItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []AbdomenItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, AbdomenItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert Abdomen library rows (accepts array of {id, desc, sort_no})
func saveAbdomenLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_abdomen SET ABDOMEN_DESC = ?, SORT_NO = ? WHERE ABDOMEN_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_abdomen (ABDOMEN_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Library endpoint: Neuro options
func getNeuroLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT NEURO_ID, NEURO_DESC FROM tsekap_lib_neuro ORDER BY SORT_NO, NEURO_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type NeuroItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []NeuroItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, NeuroItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert Neuro library rows (accepts array of {id, desc, sort_no})
func saveNeuroLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_neuro SET NEURO_DESC = ?, SORT_NO = ? WHERE NEURO_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_neuro (NEURO_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Library endpoint: Genitourinary options
func getGenitourinaryLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT GU_ID, GU_DESC FROM tsekap_lib_genitourinary ORDER BY SORT_NO, GU_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type GUItem struct {
		ID   int    `json:"id"`
		Desc string `json:"desc"`
	}

	list := []GUItem{}
	for rows.Next() {
		var id int
		var desc string
		rows.Scan(&id, &desc)
		list = append(list, GUItem{ID: id, Desc: desc})
	}
	json.NewEncoder(w).Encode(list)
}

// Save/Upsert Genitourinary library rows (accepts array of {id, desc, sort_no})
func saveGenitourinaryLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_genitourinary SET GU_DESC = ?, SORT_NO = ? WHERE GU_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_genitourinary (GU_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

// Save/Upsert Digital Rectal library rows (accepts array of {id, desc, sort_no})
func saveDigitalRectalLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []struct {
		ID     int    `json:"id"`
		Desc   string `json:"desc"`
		SortNo int    `json:"sort_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range items {
		if it.ID > 0 {
			if _, err := tx.Exec(`UPDATE tsekap_lib_digital_rectal SET RECTAL_DESC = ?, SORT_NO = ? WHERE RECTAL_ID = ?`, it.Desc, it.SortNo, it.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := tx.Exec(`INSERT INTO tsekap_lib_digital_rectal (RECTAL_DESC, SORT_NO) VALUES (?, ?)`, it.Desc, it.SortNo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}
