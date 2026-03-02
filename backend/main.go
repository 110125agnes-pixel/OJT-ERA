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

type InventoryItem struct {
	ID       int     `json:"id"`
	ItemName string  `json:"item_name"`
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	Quantity int     `json:"quantity"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
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
	// Ensure patient_femalehistory table exists (stores Female tab data)
	createPatientFemaleHistoryTable()
	// Ensure tsekap_lib_femalehistory library table exists
	createFemaleHistoryLibTable()

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
	router.HandleFunc("/api/patients/{patientId}/female-history", getFemaleHistory).Methods("GET")
	router.HandleFunc("/api/patients/{patientId}/female-history", saveFemaleHistory).Methods("POST")
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
	// Female History Library
	router.HandleFunc("/api/lib/femalehistory", getFemaleHistoryLib).Methods("GET")
	router.HandleFunc("/api/lib/femalehistory", saveFemaleHistoryLib).Methods("POST")
	router.HandleFunc("/api/lib/femalehistory/{id}", deleteFemaleHistoryLib).Methods("DELETE")
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

// getFemaleHistory handles GET /api/patients/{patientId}/female-history
// Returns the stored female history row (structured columns) for the given patient.
func getFemaleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	// translate numeric ID -> case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Read row
	row := db.QueryRow(`SELECT menarche_age, last_menstrual, period_duration_days, cycle_length_days,
		pads_per_day, sexual_onset_age, birth_control_used, is_menopause, menopause_age, is_menstrual_applicable,
		gravidity, parity, delivery_type, full_term_pregnancy_count, premature_pregnancy_count, abortion_count,
		living_children, preg_induced_htn, has_family_planning, is_preg_history_applicable, notes, date_added, added_by
		FROM patient_femalehistory WHERE patno = ?`, patno)

	var fh = map[string]interface{}{}
	var (
		menarcheAge, periodDuration, cycleLength, padsPerDay, sexualOnset               sql.NullInt64
		isMenopause, menopauseAge, isMenstrualApplicable                                sql.NullInt64
		gravidity, parity, fullTermCount, prematureCount, abortionCount, livingChildren sql.NullInt64
		pregInducedHTN, hasFamilyPlanning, isPregApplicable                             sql.NullInt64
		lastMenstrual, birthControl, deliveryType, notes, dateAdded, addedBy            sql.NullString
	)

	err := row.Scan(&menarcheAge, &lastMenstrual, &periodDuration, &cycleLength,
		&padsPerDay, &sexualOnset, &birthControl, &isMenopause, &menopauseAge, &isMenstrualApplicable,
		&gravidity, &parity, &deliveryType, &fullTermCount, &prematureCount, &abortionCount,
		&livingChildren, &pregInducedHTN, &hasFamilyPlanning, &isPregApplicable, &notes, &dateAdded, &addedBy)
	if err != nil {
		// no row found -> return empty/default object
		log.Printf("getFemaleHistory: no row found for patno=%s (patientID=%s)", patno, patientID)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// helper: strip time suffix from MySQL DATE strings returned with parseTime=true
	// e.g. "2026-03-03T00:00:00Z" -> "2026-03-03" so <input type="date"> works in browser
	trimDate := func(s string) string {
		if len(s) > 10 {
			return s[:10]
		}
		return s
	}

	// Always return ALL keys with defaults so the frontend state is fully replaced.
	// This prevents stale local state when fields haven't been filled yet (NULL in DB).
	fh["ageOfFirstMenstruation"] = ""
	if menarcheAge.Valid {
		fh["ageOfFirstMenstruation"] = fmt.Sprintf("%d", menarcheAge.Int64)
	}
	fh["dateOfLastMenstrualPeriod"] = ""
	if lastMenstrual.Valid {
		fh["dateOfLastMenstrualPeriod"] = trimDate(lastMenstrual.String)
	}
	fh["durationOfMenstrualPeriod"] = ""
	if periodDuration.Valid {
		fh["durationOfMenstrualPeriod"] = fmt.Sprintf("%d", periodDuration.Int64)
	}
	fh["intervalCycleOfMenstruation"] = ""
	if cycleLength.Valid {
		fh["intervalCycleOfMenstruation"] = fmt.Sprintf("%d", cycleLength.Int64)
	}
	fh["numberOfPadsPerDay"] = ""
	if padsPerDay.Valid {
		fh["numberOfPadsPerDay"] = fmt.Sprintf("%d", padsPerDay.Int64)
	}
	fh["onsetOfSexualIntercourse"] = ""
	if sexualOnset.Valid {
		fh["onsetOfSexualIntercourse"] = fmt.Sprintf("%d", sexualOnset.Int64)
	}
	fh["birthControlMethod"] = ""
	if birthControl.Valid {
		fh["birthControlMethod"] = birthControl.String
	}
	fh["isMenopause"] = false
	if isMenopause.Valid {
		fh["isMenopause"] = isMenopause.Int64 == 1
	}
	fh["ageOfMenopause"] = ""
	if menopauseAge.Valid {
		fh["ageOfMenopause"] = fmt.Sprintf("%d", menopauseAge.Int64)
	}
	fh["isMenstrualHistoryApplicable"] = false
	if isMenstrualApplicable.Valid {
		fh["isMenstrualHistoryApplicable"] = isMenstrualApplicable.Int64 == 1
	}
	fh["numberOfPregnancyToDate"] = ""
	if gravidity.Valid {
		fh["numberOfPregnancyToDate"] = fmt.Sprintf("%d", gravidity.Int64)
	}
	fh["numberOfDeliveryToDate"] = ""
	if parity.Valid {
		fh["numberOfDeliveryToDate"] = fmt.Sprintf("%d", parity.Int64)
	}
	fh["typeOfDelivery"] = ""
	if deliveryType.Valid {
		fh["typeOfDelivery"] = deliveryType.String
	}
	fh["numberOfFullTermPregnancy"] = ""
	if fullTermCount.Valid {
		fh["numberOfFullTermPregnancy"] = fmt.Sprintf("%d", fullTermCount.Int64)
	}
	fh["numberOfPrematurePregnancy"] = ""
	if prematureCount.Valid {
		fh["numberOfPrematurePregnancy"] = fmt.Sprintf("%d", prematureCount.Int64)
	}
	fh["numberOfAbortion"] = ""
	if abortionCount.Valid {
		fh["numberOfAbortion"] = fmt.Sprintf("%d", abortionCount.Int64)
	}
	fh["numberOfLivingChildren"] = ""
	if livingChildren.Valid {
		fh["numberOfLivingChildren"] = fmt.Sprintf("%d", livingChildren.Int64)
	}
	fh["pregnancyInducedHypertension"] = false
	if pregInducedHTN.Valid {
		fh["pregnancyInducedHypertension"] = pregInducedHTN.Int64 == 1
	}
	fh["accessToFamilyPlanningCounselling"] = false
	if hasFamilyPlanning.Valid {
		fh["accessToFamilyPlanningCounselling"] = hasFamilyPlanning.Int64 == 1
	}
	fh["isPregnancyHistoryApplicable"] = false
	if isPregApplicable.Valid {
		fh["isPregnancyHistoryApplicable"] = isPregApplicable.Int64 == 1
	}
	fh["notes"] = ""
	if notes.Valid {
		fh["notes"] = notes.String
	}
	log.Printf("getFemaleHistory: returning data for patno=%s patientID=%s: %+v", patno, patientID, fh)
	json.NewEncoder(w).Encode(fh)
}

// saveFemaleHistory handles POST /api/patients/{patientId}/female-history
// Accepts a JSON payload matching the frontend state and upserts into patient_femalehistory
func saveFemaleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// translate id -> patno
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// helper to read int-like values from payload
	toInt := func(k string) interface{} {
		if v, ok := payload[k]; ok && v != nil && v != "" {
			switch t := v.(type) {
			case float64:
				return int(t)
			case string:
				if t == "" {
					return nil
				}
				if i, err := strconv.Atoi(t); err == nil {
					return i
				}
			}
		}
		return nil
	}
	toStr := func(k string) interface{} {
		if v, ok := payload[k]; ok && v != nil {
			return fmt.Sprintf("%v", v)
		}
		return nil
	}
	toBoolInt := func(k string) int {
		if v, ok := payload[k]; ok && v != nil {
			if b, ok2 := v.(bool); ok2 {
				if b {
					return 1
				}
				return 0
			}
			if s, ok3 := v.(string); ok3 {
				if s == "true" {
					return 1
				}
				return 0
			}
		}
		return 0
	}

	log.Printf("saveFemaleHistory: received payload for patientID=%s patno=%s: %+v", patientID, patno, payload)

	_, execErr := db.Exec(`INSERT INTO patient_femalehistory (
		patno, menarche_age, last_menstrual, period_duration_days, cycle_length_days,
		pads_per_day, sexual_onset_age, birth_control_used, is_menopause, menopause_age, is_menstrual_applicable,
		gravidity, parity, delivery_type, full_term_pregnancy_count, premature_pregnancy_count, abortion_count,
		living_children, preg_induced_htn, has_family_planning, is_preg_history_applicable, notes, date_added, added_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), 'system')
	ON DUPLICATE KEY UPDATE
		menarche_age = VALUES(menarche_age), last_menstrual = VALUES(last_menstrual), period_duration_days = VALUES(period_duration_days),
		cycle_length_days = VALUES(cycle_length_days), pads_per_day = VALUES(pads_per_day), sexual_onset_age = VALUES(sexual_onset_age),
		birth_control_used = VALUES(birth_control_used), is_menopause = VALUES(is_menopause), menopause_age = VALUES(menopause_age),
		is_menstrual_applicable = VALUES(is_menstrual_applicable), gravidity = VALUES(gravidity), parity = VALUES(parity),
		delivery_type = VALUES(delivery_type), full_term_pregnancy_count = VALUES(full_term_pregnancy_count),
		premature_pregnancy_count = VALUES(premature_pregnancy_count), abortion_count = VALUES(abortion_count),
		living_children = VALUES(living_children), preg_induced_htn = VALUES(preg_induced_htn), has_family_planning = VALUES(has_family_planning),
		is_preg_history_applicable = VALUES(is_preg_history_applicable), notes = VALUES(notes), date_added = NOW()`,
		patno,
		toInt("ageOfFirstMenstruation"), toStr("dateOfLastMenstrualPeriod"), toInt("durationOfMenstrualPeriod"), toInt("intervalCycleOfMenstruation"),
		toInt("numberOfPadsPerDay"), toInt("onsetOfSexualIntercourse"), toStr("birthControlMethod"), toBoolInt("isMenopause"), toInt("ageOfMenopause"), toBoolInt("isMenstrualHistoryApplicable"),
		toInt("numberOfPregnancyToDate"), toInt("numberOfDeliveryToDate"), toStr("typeOfDelivery"), toInt("numberOfFullTermPregnancy"), toInt("numberOfPrematurePregnancy"), toInt("numberOfAbortion"),
		toInt("numberOfLivingChildren"), toBoolInt("pregnancyInducedHypertension"), toBoolInt("accessToFamilyPlanningCounselling"), toBoolInt("isPregnancyHistoryApplicable"), toStr("notes"))

	if execErr != nil {
		log.Println("saveFemaleHistory error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}

	// After saving, re-query the row and return the saved object (same shape as GET)
	row2 := db.QueryRow(`SELECT menarche_age, last_menstrual, period_duration_days, cycle_length_days,
		pads_per_day, sexual_onset_age, birth_control_used, is_menopause, menopause_age, is_menstrual_applicable,
		gravidity, parity, delivery_type, full_term_pregnancy_count, premature_pregnancy_count, abortion_count,
		living_children, preg_induced_htn, has_family_planning, is_preg_history_applicable, notes, date_added, added_by
		FROM patient_femalehistory WHERE patno = ?`, patno)

	var (
		m2, pd2, cl2, pp2, so2                      sql.NullInt64
		im2, ma2, ima2                              sql.NullInt64
		g2, p2, ft2, prem2, ab2, lc2                sql.NullInt64
		pih2, hfp2, ipa2                            sql.NullInt64
		lm2, bc2, dt2, notes2, dateAdded2, addedBy2 sql.NullString
	)
	err2 := row2.Scan(&m2, &lm2, &pd2, &cl2, &pp2, &so2, &bc2, &im2, &ma2, &ima2, &g2, &p2, &dt2, &ft2, &prem2, &ab2, &lc2, &pih2, &hfp2, &ipa2, &notes2, &dateAdded2, &addedBy2)
	if err2 != nil {
		// saved but cannot re-read; return minimal confirmation
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved", "patno": patno})
		return
	}

	// helper for saveFemaleHistory: strip time suffix from date strings
	trimDate2 := func(s string) string {
		if len(s) > 10 {
			return s[:10]
		}
		return s
	}

	saved := map[string]interface{}{}
	saved["ageOfFirstMenstruation"] = ""
	if m2.Valid {
		saved["ageOfFirstMenstruation"] = fmt.Sprintf("%d", m2.Int64)
	}
	saved["dateOfLastMenstrualPeriod"] = ""
	if lm2.Valid {
		saved["dateOfLastMenstrualPeriod"] = trimDate2(lm2.String)
	}
	saved["durationOfMenstrualPeriod"] = ""
	if pd2.Valid {
		saved["durationOfMenstrualPeriod"] = fmt.Sprintf("%d", pd2.Int64)
	}
	saved["intervalCycleOfMenstruation"] = ""
	if cl2.Valid {
		saved["intervalCycleOfMenstruation"] = fmt.Sprintf("%d", cl2.Int64)
	}
	saved["numberOfPadsPerDay"] = ""
	if pp2.Valid {
		saved["numberOfPadsPerDay"] = fmt.Sprintf("%d", pp2.Int64)
	}
	saved["onsetOfSexualIntercourse"] = ""
	if so2.Valid {
		saved["onsetOfSexualIntercourse"] = fmt.Sprintf("%d", so2.Int64)
	}
	saved["birthControlMethod"] = ""
	if bc2.Valid {
		saved["birthControlMethod"] = bc2.String
	}
	saved["isMenopause"] = false
	if im2.Valid {
		saved["isMenopause"] = im2.Int64 == 1
	}
	saved["ageOfMenopause"] = ""
	if ma2.Valid {
		saved["ageOfMenopause"] = fmt.Sprintf("%d", ma2.Int64)
	}
	saved["isMenstrualHistoryApplicable"] = false
	if ima2.Valid {
		saved["isMenstrualHistoryApplicable"] = ima2.Int64 == 1
	}
	saved["numberOfPregnancyToDate"] = ""
	if g2.Valid {
		saved["numberOfPregnancyToDate"] = fmt.Sprintf("%d", g2.Int64)
	}
	saved["numberOfDeliveryToDate"] = ""
	if p2.Valid {
		saved["numberOfDeliveryToDate"] = fmt.Sprintf("%d", p2.Int64)
	}
	saved["typeOfDelivery"] = ""
	if dt2.Valid {
		saved["typeOfDelivery"] = dt2.String
	}
	saved["numberOfFullTermPregnancy"] = ""
	if ft2.Valid {
		saved["numberOfFullTermPregnancy"] = fmt.Sprintf("%d", ft2.Int64)
	}
	saved["numberOfPrematurePregnancy"] = ""
	if prem2.Valid {
		saved["numberOfPrematurePregnancy"] = fmt.Sprintf("%d", prem2.Int64)
	}
	saved["numberOfAbortion"] = ""
	if ab2.Valid {
		saved["numberOfAbortion"] = fmt.Sprintf("%d", ab2.Int64)
	}
	saved["numberOfLivingChildren"] = ""
	if lc2.Valid {
		saved["numberOfLivingChildren"] = fmt.Sprintf("%d", lc2.Int64)
	}
	saved["pregnancyInducedHypertension"] = false
	if pih2.Valid {
		saved["pregnancyInducedHypertension"] = pih2.Int64 == 1
	}
	saved["accessToFamilyPlanningCounselling"] = false
	if hfp2.Valid {
		saved["accessToFamilyPlanningCounselling"] = hfp2.Int64 == 1
	}
	saved["isPregnancyHistoryApplicable"] = false
	if ipa2.Valid {
		saved["isPregnancyHistoryApplicable"] = ipa2.Int64 == 1
	}
	saved["notes"] = ""
	if notes2.Valid {
		saved["notes"] = notes2.String
	}
	saved["date_added"] = ""
	if dateAdded2.Valid {
		saved["date_added"] = dateAdded2.String
	}
	saved["added_by"] = ""
	if addedBy2.Valid {
		saved["added_by"] = addedBy2.String
	}

	log.Printf("saveFemaleHistory: saved data for patno=%s patientID=%s: %+v", patno, patientID, saved)
	json.NewEncoder(w).Encode(saved)
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
	rows, err := db.Query("SELECT SURGERY_CODE, SURGERY_DESC FROM tsekap_lib_surgical ORDER BY SORT_NO, SURGERY_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SurgItem struct {
		SURGERY_CODE string `json:"SURGERY_CODE"`
		SURGERY_DESC string `json:"SURGERY_DESC"`
	}

	list := []SurgItem{}
	for rows.Next() {
		var code, desc string
		rows.Scan(&code, &desc)
		list = append(list, SurgItem{SURGERY_CODE: code, SURGERY_DESC: desc})
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

// createPatientFemaleHistoryTable creates a table to store the Female tab entries per patient.
// Structured columns are used so individual fields can be queried easily.
func createPatientFemaleHistoryTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_femalehistory (
		patno VARCHAR(20) NOT NULL PRIMARY KEY,
		menarche_age INT,
		last_menstrual DATE,
		period_duration_days INT,
		cycle_length_days INT,
		pads_per_day INT,
		sexual_onset_age INT,
		birth_control_used VARCHAR(100),
		is_menopause TINYINT(1),
		menopause_age INT,
		is_menstrual_applicable TINYINT(1),

		gravidity INT,
		parity INT,
		delivery_type VARCHAR(32),
		full_term_pregnancy_count INT,
		premature_pregnancy_count INT,
		abortion_count INT,
		living_children INT,
		preg_induced_htn TINYINT(1),
		has_family_planning TINYINT(1),
		is_preg_history_applicable TINYINT(1),

		notes TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_femalehistory table ready")
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

// ==================== FEMALE HISTORY LIBRARY ====================

func createFemaleHistoryLibTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tsekap_lib_femalehistory (
		FH_ID INT AUTO_INCREMENT PRIMARY KEY,
		FH_FIELD_KEY VARCHAR(100) NOT NULL,
		FH_LABEL VARCHAR(255) NOT NULL,
		FH_SECTION VARCHAR(100),
		FH_NUMBER INT DEFAULT 0,
		SORT_NO INT DEFAULT 0,
		LIB_STAT TINYINT(1) DEFAULT 1
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	if err != nil {
		log.Printf("Warning createFemaleHistoryLibTable: %v", err)
	} else {
		log.Println("✓ tsekap_lib_femalehistory table ready")
	}
	seedFemaleHistoryLib()
}

func seedFemaleHistoryLib() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM tsekap_lib_femalehistory").Scan(&count)
	if count > 0 {
		return
	}
	seeds := []struct {
		key     string
		label   string
		section string
		number  int
		sortNo  int
	}{
		{"ageOfFirstMenstruation", "Age of First Menstruation (Menarche)", "Menstrual History", 1, 1},
		{"dateOfLastMenstrualPeriod", "Date of Last Menstrual Period", "Menstrual History", 2, 2},
		{"durationOfMenstrualPeriod", "Duration of Menstrual Period in Number of Days", "Menstrual History", 3, 3},
		{"intervalCycleOfMenstruation", "Interval/Cycle of Menstruation in Number of Days", "Menstrual History", 4, 4},
		{"numberOfPadsPerDay", "Number of Pads/Napkins Used per Day during Menstruation", "Menstrual History", 5, 5},
		{"onsetOfSexualIntercourse", "Onset of Sexual Intercourse (Age of First Sexual Intercourse)", "Menstrual History", 6, 6},
		{"birthControlMethod", "Birth Control Method Used", "Menstrual History", 7, 7},
		{"isMenopause", "Is Menopause?", "Menstrual History", 8, 8},
		{"ageOfMenopause", "If Menopause, Age of Menopause", "Menstrual History", 9, 9},
		{"isMenstrualHistoryApplicable", "Is menstrual history applicable?", "Menstrual History", 10, 10},
		{"numberOfPregnancyToDate", "Number of Pregnancy to Date - Gravity Chief", "Pregnancy History", 1, 11},
		{"numberOfDeliveryToDate", "Number of Delivery to Date - Parity", "Pregnancy History", 2, 12},
		{"typeOfDelivery", "Type of Delivery", "Pregnancy History", 3, 13},
		{"numberOfFullTermPregnancy", "Number of Full Term Pregnancy", "Pregnancy History", 4, 14},
		{"numberOfPrematurePregnancy", "Number of Premature Pregnancy", "Pregnancy History", 5, 15},
		{"numberOfAbortion", "Number of Abortion", "Pregnancy History", 6, 16},
		{"numberOfLivingChildren", "Number of Living Children", "Pregnancy History", 7, 17},
		{"pregnancyInducedHypertension", "If Pregnancy - Induced Hypertension (Pre - Eclampsia)", "Pregnancy History", 8, 18},
		{"accessToFamilyPlanningCounselling", "If with access to Family Planning Counselling", "Pregnancy History", 9, 19},
		{"isPregnancyHistoryApplicable", "Is pregnancy history applicable?", "Pregnancy History", 10, 20},
	}
	for _, s := range seeds {
		db.Exec(`INSERT INTO tsekap_lib_femalehistory (FH_FIELD_KEY, FH_LABEL, FH_SECTION, FH_NUMBER, SORT_NO, LIB_STAT) VALUES (?, ?, ?, ?, ?, 1)`,
			s.key, s.label, s.section, s.number, s.sortNo)
	}
	log.Println("✓ tsekap_lib_femalehistory seeded with 20 fields")
}

func getFemaleHistoryLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT FH_ID, FH_FIELD_KEY, FH_LABEL, FH_SECTION, FH_NUMBER, SORT_NO, LIB_STAT FROM tsekap_lib_femalehistory ORDER BY SORT_NO, FH_ID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type FHItem struct {
		ID       int    `json:"id"`
		FieldKey string `json:"field_key"`
		Label    string `json:"label"`
		Section  string `json:"section"`
		Number   int    `json:"number"`
		SortNo   int    `json:"sort_no"`
		LibStat  int    `json:"lib_stat"`
	}
	list := []FHItem{}
	for rows.Next() {
		var it FHItem
		rows.Scan(&it.ID, &it.FieldKey, &it.Label, &it.Section, &it.Number, &it.SortNo, &it.LibStat)
		list = append(list, it)
	}
	json.NewEncoder(w).Encode(list)
}

func saveFemaleHistoryLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item struct {
		ID      int    `json:"id"`
		Label   string `json:"label"`
		SortNo  int    `json:"sort_no"`
		LibStat int    `json:"lib_stat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.ID > 0 {
		_, err := db.Exec(`UPDATE tsekap_lib_femalehistory SET FH_LABEL=?, SORT_NO=?, LIB_STAT=? WHERE FH_ID=?`,
			item.Label, item.SortNo, item.LibStat, item.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Adding new rows not supported for this library", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func deleteFemaleHistoryLib(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	_, err := db.Exec(`DELETE FROM tsekap_lib_femalehistory WHERE FH_ID = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted"})
}
