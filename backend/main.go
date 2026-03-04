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

	// Ensure patient_famhist summary table exists (family history stored as pipe-separated 1|0)
	createFamilyHistSummaryTable()
	// Ensure patient_immunization summary table exists (patno keyed, 1|0 bits)
	createImmunizationSummaryTable()

	// Ensure patient_pe_skin summary table exists (skin findings stored as pipe-separated 1|0)
	createPatientPeSkinTable()

	// Ensure patient_pe_genitourinary summary table exists (genitourinary findings stored as pipe-separated 1|0)
	createPatientPeGenitourinaryTable()

	// Ensure patient_pe_digital_rectal summary table exists (digital rectal findings stored as pipe-separated 1|0)
	createPatientPeDigitalRectalTable()

	// Ensure patient_pe_neuro summary table exists (neurological findings stored as pipe-separated 1|0)
	createPatientPeNeuroTable()

	// Ensure patient_pe_heent summary table exists (HEENT findings stored as pipe-separated 1|0)
	createPatientPeHeentTable()

	// Ensure patient_pe_chest summary table exists (Chest findings stored as pipe-separated 1|0)
	createPatientPeChestTable()

	// Ensure patient_pe_heart summary table exists (Heart findings stored as pipe-separated 1|0)
	createPatientPeHeartTable()

	// Ensure patient_pe_abdomen summary table exists (Abdomen findings stored as pipe-separated 1|0)
	createPatientPeAbdomenTable()

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

// ==================== SOCIAL HISTORY HANDLERS ====================

// getSocialHistory reads social history from patient_socialhistory using patient's case_no (patno).
func getSocialHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	// translate numeric ID -> case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// read row from patient_socialhistory (store strings like "Yes"/"No"/"Quit")
	var (
		isSmoker, cigPacks, isAlcohol, bottlesPerDay, isIllicit, isSexually sql.NullString
	)
	err := db.QueryRow(`SELECT is_patient_smoker, cigarette_packs_per_year, is_alcohol_drinker, bottles_per_day, is_illicit_drug_user, is_sexually_active
		FROM patient_socialhistory WHERE patno = ?`, patno).Scan(&isSmoker, &cigPacks, &isAlcohol, &bottlesPerDay, &isIllicit, &isSexually)
	if err != nil {
		// return empty/default object when not found
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	resp := map[string]interface{}{
		"is_patient_smoker": func() string {
			if isSmoker.Valid {
				return isSmoker.String
			}
			return "No"
		}(),
		"cigarette_packs_per_year": func() int {
			if cigPacks.Valid {
				if i, err := strconv.Atoi(cigPacks.String); err == nil {
					return i
				}
			}
			return 0
		}(),
		"is_alcohol_drinker": func() string {
			if isAlcohol.Valid {
				return isAlcohol.String
			}
			return "No"
		}(),
		"bottles_per_day": func() int {
			if bottlesPerDay.Valid {
				if i, err := strconv.Atoi(bottlesPerDay.String); err == nil {
					return i
				}
			}
			return 0
		}(),
		"is_illicit_drug_user": func() string {
			if isIllicit.Valid {
				return isIllicit.String
			}
			return "No"
		}(),
		"is_sexually_active": func() string {
			if isSexually.Valid {
				return isSexually.String
			}
			return "No"
		}(),
	}
	json.NewEncoder(w).Encode(resp)
}

// saveSocialHistory upserts social history into patient_socialhistory mapping Yes/No/Quit -> 1/2/0
func saveSocialHistory(w http.ResponseWriter, r *http.Request) {
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

	// helpers to extract string and int values from payload
	toStr := func(k string) string {
		if v, ok := payload[k]; ok && v != nil {
			return fmt.Sprintf("%v", v)
		}
		return "No"
	}
	toInt := func(k string) int {
		if v, ok := payload[k]; ok && v != nil {
			switch t := v.(type) {
			case float64:
				return int(t)
			case string:
				if t == "" {
					return 0
				}
				if i, err := strconv.Atoi(t); err == nil {
					return i
				}
			}
		}
		return 0
	}

	// Upsert into patient_socialhistory (patno primary key)
	_, execErr := db.Exec(`INSERT INTO patient_socialhistory (
		patno, is_patient_smoker, cigarette_packs_per_year, is_alcohol_drinker, bottles_per_day, is_illicit_drug_user, is_sexually_active, date_added, added_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), 'system')
	ON DUPLICATE KEY UPDATE
		is_patient_smoker = VALUES(is_patient_smoker),
		cigarette_packs_per_year = VALUES(cigarette_packs_per_year),
		is_alcohol_drinker = VALUES(is_alcohol_drinker),
		bottles_per_day = VALUES(bottles_per_day),
		is_illicit_drug_user = VALUES(is_illicit_drug_user),
		is_sexually_active = VALUES(is_sexually_active),
		date_added = NOW()`,
		patno,
		toStr("is_patient_smoker"), toInt("cigarette_packs_per_year"), toStr("is_alcohol_drinker"), toInt("bottles_per_day"), toStr("is_illicit_drug_user"), toStr("is_sexually_active"))

	if execErr != nil {
		log.Println("saveSocialHistory error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}

	// return the saved payload as frontend expects strings for radio fields
	saved := map[string]interface{}{
		"is_patient_smoker":        payload["is_patient_smoker"],
		"cigarette_packs_per_year": toInt("cigarette_packs_per_year"),
		"is_alcohol_drinker":       payload["is_alcohol_drinker"],
		"bottles_per_day":          toInt("bottles_per_day"),
		"is_illicit_drug_user":     payload["is_illicit_drug_user"],
		"is_sexually_active":       payload["is_sexually_active"],
	}
	json.NewEncoder(w).Encode(saved)
}

// ==================== PHYSICAL EXAM HANDLERS ====================

func getPertinentPhysicalExam(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["patientId"]
	var ppe PertinentPhysicalExam
	// Read from the new table patient_pertinent_physical_exam. Column names in DB use _cm/_kg suffixes for measurements.
	err := db.QueryRow(`SELECT id, patient_id, systolic_bp, diastolic_bp, heart_rate, respiratory_rate,
			temperature, height_cm, weight_kg, bmi, pzscore, right_eye_vision, left_eye_vision,
			length_pediatric_cm, head_circumference_cm, skinfold_thickness_cm, waist_cm, hip_cm, limbs_cm, arm_circumference_cm
			FROM patient_pertinent_physical_exam WHERE patient_id = ?`, patientID).Scan(
		&ppe.ID, &ppe.PatientID, &ppe.SystolicBP, &ppe.DiastolicBP, &ppe.HeartRate, &ppe.RespiratoryRate,
		&ppe.Temperature, &ppe.Height, &ppe.Weight, &ppe.BMI, &ppe.PZScore, &ppe.RightEyeVision,
		&ppe.LeftEyeVision, &ppe.LengthPediatric, &ppe.HeadCircumference, &ppe.SkinfoldThickness,
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
	if err := json.NewDecoder(r.Body).Decode(&ppe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists int
	db.QueryRow("SELECT COUNT(*) FROM patient_pertinent_physical_exam WHERE patient_id = ?", patientID).Scan(&exists)

	if exists > 0 {
		db.Exec(`UPDATE patient_pertinent_physical_exam SET systolic_bp=?, diastolic_bp=?, heart_rate=?,
			respiratory_rate=?, temperature=?, height_cm=?, weight_kg=?, bmi=?, pzscore=?, right_eye_vision=?,
			left_eye_vision=?, length_pediatric_cm=?, head_circumference_cm=?, skinfold_thickness_cm=?,
			waist_cm=?, hip_cm=?, limbs_cm=?, arm_circumference_cm=? WHERE patient_id=?`,
			ppe.SystolicBP, ppe.DiastolicBP, ppe.HeartRate, ppe.RespiratoryRate, ppe.Temperature,
			ppe.Height, ppe.Weight, ppe.BMI, ppe.PZScore, ppe.RightEyeVision, ppe.LeftEyeVision,
			ppe.LengthPediatric, ppe.HeadCircumference, ppe.SkinfoldThickness, ppe.Waist, ppe.Hip,
			ppe.Limbs, ppe.ArmCircumference, patientID)
	} else {
		db.Exec(`INSERT INTO patient_pertinent_physical_exam (patient_id, systolic_bp, diastolic_bp, heart_rate,
			respiratory_rate, temperature, height_cm, weight_kg, bmi, pzscore, right_eye_vision,
			left_eye_vision, length_pediatric_cm, head_circumference_cm, skinfold_thickness_cm,
			waist_cm, hip_cm, limbs_cm, arm_circumference_cm)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			patientID, ppe.SystolicBP, ppe.DiastolicBP, ppe.HeartRate, ppe.RespiratoryRate, ppe.Temperature,
			ppe.Height, ppe.Weight, ppe.BMI, ppe.PZScore, ppe.RightEyeVision, ppe.LeftEyeVision,
			ppe.LengthPediatric, ppe.HeadCircumference, ppe.SkinfoldThickness, ppe.Waist, ppe.Hip,
			ppe.Limbs, ppe.ArmCircumference)
	}
	json.NewEncoder(w).Encode(ppe)
}

// createPatientPertinentPhysicalExamTable creates a table to store Pertinent Physical Examination per patient.
func createPatientPertinentPhysicalExamTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pertinent_physical_exam (
		id INT AUTO_INCREMENT PRIMARY KEY,
		patient_id INT NOT NULL UNIQUE,
		systolic_bp INT,
		diastolic_bp INT,
		heart_rate INT,
		respiratory_rate INT,
		temperature DECIMAL(5,2),
		height_cm DECIMAL(6,2),
		weight_kg DECIMAL(6,2),
		bmi DECIMAL(5,2),
		pzscore INT,
		right_eye_vision VARCHAR(32),
		left_eye_vision VARCHAR(32),
		length_pediatric_cm DECIMAL(6,2),
		head_circumference_cm DECIMAL(6,2),
		skinfold_thickness_cm DECIMAL(6,2),
		waist_cm DECIMAL(6,2),
		hip_cm DECIMAL(6,2),
		limbs_cm DECIMAL(6,2),
		arm_circumference_cm DECIMAL(6,2),
		remarks TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY (patient_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pertinent_physical_exam table ready")
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
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	// Translate numeric patient ID → case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Load library order (defines bit positions)
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

	// Fetch saved pipe-separated string for this patient
	var saved string
	var notes sql.NullString
	db.QueryRow("SELECT fdisease_code, notes FROM patient_famhist WHERE patno = ?", patno).Scan(&saved, &notes)
	bits := []string{}
	if saved != "" {
		bits = strings.Split(saved, "|")
	}

	// Build response
	var list []FamilyHistoryItem
	for i, d := range lib {
		isChecked := false
		if i < len(bits) {
			isChecked = bits[i] == "1"
		}
		fh := FamilyHistoryItem{
			DiseaseCode: d.Code,
			DiseaseName: d.Desc,
			IsChecked:   isChecked,
		}
		// attach notes on first element for frontend compatibility
		if i == 0 && notes.Valid {
			fh.Notes = notes.String
		}
		list = append(list, fh)
	}
	json.NewEncoder(w).Encode(list)
}

func saveFamilyHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	// Decode incoming items array
	var items []FamilyHistoryItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Translate numeric patient ID → case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Re-fetch library order
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

	// Build lookup map from request body
	checkedMap := map[string]bool{}
	var notes string
	for _, it := range items {
		checkedMap[it.DiseaseCode] = it.IsChecked
		if it.Notes != "" {
			notes = it.Notes
		}
	}

	bits := make([]string, len(libOrder))
	for i, code := range libOrder {
		if checkedMap[code] {
			bits[i] = "1"
		} else {
			bits[i] = "0"
		}
	}
	fdiseaseCode := strings.Join(bits, "|")

	// UPSERT into patient_famhist
	_, execErr := db.Exec(
		`INSERT INTO patient_famhist (patno, fdisease_code, notes, date_added, added_by)
		 VALUES (?, ?, ?, NOW(), 'system')
		 ON DUPLICATE KEY UPDATE fdisease_code = VALUES(fdisease_code), notes = VALUES(notes), date_added = NOW()`,
		patno, fdiseaseCode, notes)
	if execErr != nil {
		log.Println("saveFamilyHistory error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved", "patno": patno})
}

func getSurgicalHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	// Load the ordered surgical library
	libRows, err := db.Query("SELECT SURG_CODE, SURG_DESC FROM tsekap_lib_surgical WHERE LIB_STAT=1 ORDER BY SORT_NO, SURG_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer libRows.Close()
	type SurgLib struct {
		Code string
		Desc string
	}
	var libList []SurgLib
	for libRows.Next() {
		var s SurgLib
		libRows.Scan(&s.Code, &s.Desc)
		libList = append(libList, s)
	}

	// Load patient's saved bit string
	var surgCode string
	db.QueryRow("SELECT surg_code FROM patient_surgery WHERE patno = ?", patno).Scan(&surgCode)

	bits := strings.Split(surgCode, "|")

	type SurgItem struct {
		SurgeryCode string `json:"SurgeryCode"`
		SurgeryName string `json:"SurgeryName"`
		IsChecked   bool   `json:"IsChecked"`
	}
	list := []SurgItem{}
	for i, s := range libList {
		checked := false
		if i < len(bits) && bits[i] == "1" {
			checked = true
		}
		list = append(list, SurgItem{SurgeryCode: s.Code, SurgeryName: s.Desc, IsChecked: checked})
	}
	json.NewEncoder(w).Encode(list)
}

func saveSurgicalHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		http.Error(w, "patient not found", http.StatusNotFound)
		return
	}

	// Load ordered library to build bit positions
	libRows, err := db.Query("SELECT SURG_CODE FROM tsekap_lib_surgical WHERE LIB_STAT=1 ORDER BY SORT_NO, SURG_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer libRows.Close()
	var libCodes []string
	for libRows.Next() {
		var code string
		libRows.Scan(&code)
		libCodes = append(libCodes, code)
	}

	// Decode incoming payload (frontend sends PascalCase keys)
	type SaveItem struct {
		SurgeryCode string `json:"SurgeryCode"`
		IsChecked   bool   `json:"IsChecked"`
	}
	var items []SaveItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build checked set
	checkedSet := map[string]bool{}
	for _, it := range items {
		if it.IsChecked {
			checkedSet[it.SurgeryCode] = true
		}
	}

	// Build pipe-separated bit string
	bits := make([]string, len(libCodes))
	for i, code := range libCodes {
		if checkedSet[code] {
			bits[i] = "1"
		} else {
			bits[i] = "0"
		}
	}
	surgCode := strings.Join(bits, "|")

	_, execErr := db.Exec(`INSERT INTO patient_surgery (patno, surg_code, date_added, added_by)
		VALUES (?, ?, NOW(), 'system')
		ON DUPLICATE KEY UPDATE surg_code=VALUES(surg_code), date_added=NOW()`,
		patno, surgCode)
	if execErr != nil {
		log.Println("saveSurgicalHistory error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
}

func getImmunization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	// Translate numeric patient ID → case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Build library order: child, young, preg, elderly (preserves positions)
	type libEntry struct{ Code, Name, Category string }
	var lib []libEntry

	rows, err := db.Query("SELECT IMM_CODE, IMM_DESC FROM tsekap_lib_immchild ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c, n string
			rows.Scan(&c, &n)
			lib = append(lib, libEntry{Code: c, Name: n, Category: "child"})
		}
		rows.Close()
	}

	rows, err = db.Query("SELECT IMM_CODE, IMM_DESC FROM tsekap_lib_immyoungw ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c, n string
			rows.Scan(&c, &n)
			lib = append(lib, libEntry{Code: c, Name: n, Category: "young"})
		}
		rows.Close()
	}

	rows, err = db.Query("SELECT IMM_CODE, IMM_DESC FROM tsekap_lib_immpregw ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c, n string
			rows.Scan(&c, &n)
			lib = append(lib, libEntry{Code: c, Name: n, Category: "pregnant"})
		}
		rows.Close()
	}

	rows, err = db.Query("SELECT IMM_CODE, IMM_DESC FROM tsekap_lib_immelderly ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c, n string
			rows.Scan(&c, &n)
			lib = append(lib, libEntry{Code: c, Name: n, Category: "elderly"})
		}
		rows.Close()
	}

	// Fetch saved group strings for this patient (one column per category)
	var savedChild, savedYoung, savedPreg, savedElderly, otherNotes sql.NullString
	db.QueryRow("SELECT imm_child, imm_young, imm_pregnant, imm_elderly, other_notes FROM patient_immunization WHERE patno = ?", patno).
		Scan(&savedChild, &savedYoung, &savedPreg, &savedElderly, &otherNotes)

	childBits := []string{}
	youngBits := []string{}
	pregBits := []string{}
	elderlyBits := []string{}
	if savedChild.Valid && savedChild.String != "" {
		childBits = strings.Split(savedChild.String, "|")
	}
	if savedYoung.Valid && savedYoung.String != "" {
		youngBits = strings.Split(savedYoung.String, "|")
	}
	if savedPreg.Valid && savedPreg.String != "" {
		pregBits = strings.Split(savedPreg.String, "|")
	}
	if savedElderly.Valid && savedElderly.String != "" {
		elderlyBits = strings.Split(savedElderly.String, "|")
	}

	// Build response using per-category bit arrays
	var list []ImmunizationItem
	// counters to track index within each category
	childIdx, youngIdx, pregIdx, elderlyIdx := 0, 0, 0, 0
	for _, e := range lib {
		isChecked := false
		switch e.Category {
		case "child":
			if childIdx < len(childBits) {
				isChecked = childBits[childIdx] == "1"
			}
			childIdx++
		case "young":
			if youngIdx < len(youngBits) {
				isChecked = youngBits[youngIdx] == "1"
			}
			youngIdx++
		case "pregnant":
			if pregIdx < len(pregBits) {
				isChecked = pregBits[pregIdx] == "1"
			}
			pregIdx++
		case "elderly":
			if elderlyIdx < len(elderlyBits) {
				isChecked = elderlyBits[elderlyIdx] == "1"
			}
			elderlyIdx++
		}
		list = append(list, ImmunizationItem{
			VaccineCode:      e.Code,
			VaccineName:      e.Name,
			Category:         e.Category,
			IsChecked:        isChecked,
			OtherDescription: otherNotes.String,
		})
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
	w.Header().Set("Content-Type", "application/json")
	patientID := mux.Vars(r)["patientId"]

	var items []ImmunizationItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Translate numeric patient ID → case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Rebuild per-category library orders (same order as getImmunization)
	childOrder := []string{}
	youngOrder := []string{}
	pregOrder := []string{}
	elderlyOrder := []string{}

	rows, err := db.Query("SELECT IMM_CODE FROM tsekap_lib_immchild ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c string
			rows.Scan(&c)
			childOrder = append(childOrder, c)
		}
		rows.Close()
	}
	rows, err = db.Query("SELECT IMM_CODE FROM tsekap_lib_immyoungw ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c string
			rows.Scan(&c)
			youngOrder = append(youngOrder, c)
		}
		rows.Close()
	}
	rows, err = db.Query("SELECT IMM_CODE FROM tsekap_lib_immpregw ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c string
			rows.Scan(&c)
			pregOrder = append(pregOrder, c)
		}
		rows.Close()
	}
	rows, err = db.Query("SELECT IMM_CODE FROM tsekap_lib_immelderly ORDER BY IMM_CODE")
	if err == nil {
		for rows.Next() {
			var c string
			rows.Scan(&c)
			elderlyOrder = append(elderlyOrder, c)
		}
		rows.Close()
	}

	// Map incoming items by code
	checkedMap := map[string]bool{}
	var otherNotes string
	for _, it := range items {
		if it.VaccineCode != "" {
			checkedMap[it.VaccineCode] = it.IsChecked
		}
		// capture other_notes if frontend provides it as OtherDescription
		if it.OtherDescription != "" {
			otherNotes = it.OtherDescription
		}
	}

	// Build bits per category
	buildBits := func(order []string) string {
		if len(order) == 0 {
			return ""
		}
		bits := make([]string, len(order))
		for i, code := range order {
			if checkedMap[code] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		return strings.Join(bits, "|")
	}

	immChild := buildBits(childOrder)
	immYoung := buildBits(youngOrder)
	immPreg := buildBits(pregOrder)
	immElderly := buildBits(elderlyOrder)

	_, execErr := db.Exec(
		`INSERT INTO patient_immunization (patno, imm_child, imm_young, imm_pregnant, imm_elderly, other_notes, date_added, added_by)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), 'system')
		 ON DUPLICATE KEY UPDATE imm_child = VALUES(imm_child), imm_young = VALUES(imm_young), imm_pregnant = VALUES(imm_pregnant), imm_elderly = VALUES(imm_elderly), other_notes = VALUES(other_notes), date_added = NOW()`,
		patno, immChild, immYoung, immPreg, immElderly, otherNotes)
	if execErr != nil {
		log.Println("saveImmunization error:", execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Saved", "patno": patno})
}

// createImmunizationSummaryTable creates patient_immunization table used to store
// a pipe-separated 1|0 string keyed by patno (case_no) similar to patient_medhist
func createImmunizationSummaryTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_immunization (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		imm_child VARCHAR(2000) NOT NULL DEFAULT '',
		imm_young VARCHAR(2000) NOT NULL DEFAULT '',
		imm_pregnant VARCHAR(2000) NOT NULL DEFAULT '',
		imm_elderly VARCHAR(2000) NOT NULL DEFAULT '',
		other_notes TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_immunization table ready")
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
	rows, err := db.Query("SELECT SURG_CODE, SURG_DESC FROM tsekap_lib_surgical WHERE LIB_STAT=1 ORDER BY SORT_NO, SURG_CODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SurgItem struct {
		Code string `json:"code"`
		Desc string `json:"desc"`
	}

	list := []SurgItem{}
	for rows.Next() {
		var code, desc string
		rows.Scan(&code, &desc)
		list = append(list, SurgItem{Code: code, Desc: desc})
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

// createFamilyHistSummaryTable auto-creates the patient_famhist table on backend startup.
// Schema design:
//
//	patno        - the patient's case_no (e.g. C2026-00001), serves as PRIMARY KEY
//	fdisease_code - pipe-separated 0/1 string, one bit per disease in tsekap_lib_mdiseases order
//	               e.g. "1|0|1|0" means disease 1 and 3 are checked
//	notes        - optional freeform notes saved from the frontend
//	date_added   - timestamp of last save
//	added_by     - who saved (currently always 'system')
func createFamilyHistSummaryTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_famhist (
		patno VARCHAR(20) NOT NULL PRIMARY KEY,
		fdisease_code VARCHAR(500) NOT NULL DEFAULT '',
		notes TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_famhist table ready")
}

// createPatientPeSkinTable auto-creates the patient_pe_skin table on backend startup.
// Schema design:
//
//	patno     - patient's case_no, PRIMARY KEY
//	skin_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_skin_extremities order
//	others_text - optional text saved for 'others' option
//	date_added, added_by
func createPatientPeSkinTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_skin (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		skin_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_skin table ready")
}

// createPatientPeGenitourinaryTable auto-creates the patient_pe_genitourinary table on backend startup.
// Schema design:
//
//	patno     - patient's case_no, PRIMARY KEY
//	gu_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_genitourinary order
//	others_text - optional text saved for 'others' option
//	date_added, added_by
func createPatientPeGenitourinaryTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_genitourinary (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		gu_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_genitourinary table ready")
}

// createPatientPeDigitalRectalTable auto-creates the patient_pe_digital_rectal table on backend startup.
// Schema design:
//
//	patno     - patient's case_no, PRIMARY KEY
//	dr_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_digital_rectal order
//	others_text - optional text saved for 'others' option
//	date_added, added_by
func createPatientPeDigitalRectalTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_digital_rectal (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		dr_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_digital_rectal table ready")
}

// createPatientPeNeuroTable auto-creates the patient_pe_neuro table on backend startup.
// Schema design:
//
//	patno     - patient's case_no, PRIMARY KEY
//	neuro_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_neuro order
//	others_text - optional text saved for 'others' option
//	date_added, added_by
func createPatientPeNeuroTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_neuro (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		neuro_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_neuro table ready")
}

// createPatientPeHeentTable auto-creates the patient_pe_heent table on backend startup.
// Schema design:
//
//    patno     - patient's case_no, PRIMARY KEY
//    heent_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_heent order
//    others_text - optional text saved for 'others' option
//    date_added, added_by
func createPatientPeHeentTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_heent (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		heent_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_heent table ready")
}

// createPatientPeChestTable auto-creates the patient_pe_chest table on backend startup.
// Schema design:
//
//    patno     - patient's case_no, PRIMARY KEY
//    chest_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_chest order
//    others_text - optional text saved for 'others' option
//    date_added, added_by
func createPatientPeChestTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_chest (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		chest_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_chest table ready")
}

// createPatientPeHeartTable auto-creates the patient_pe_heart table on backend startup.
// Schema design:
//
//    patno     - patient's case_no, PRIMARY KEY
//    heart_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_heart order
//    others_text - optional text saved for 'others' option
//    date_added, added_by
func createPatientPeHeartTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_heart (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		heart_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_heart table ready")
}

// createPatientPeAbdomenTable auto-creates the patient_pe_abdomen table on backend startup.
// Schema design:
//
//    patno     - patient's case_no, PRIMARY KEY
//    abdomen_code - pipe-separated 0/1 string, one bit per option in tsekap_lib_abdomen order
//    others_text - optional text saved for 'others' option
//    date_added, added_by
func createPatientPeAbdomenTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS patient_pe_abdomen (
		patno VARCHAR(50) NOT NULL PRIMARY KEY,
		abdomen_code VARCHAR(2000) NOT NULL DEFAULT '',
		others_text TEXT,
		date_added DATETIME,
		added_by VARCHAR(50)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	log.Println("✓ patient_pe_abdomen table ready")
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
	// First, load existing rows for non-skin categories
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

	// Now load skin findings from the summary table (if present), and expand into FindingRow entries
	// Translate numeric patient ID → case_no (patno)
	var patno string
	db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
	if patno == "" {
		patno = patientID
	}

	// Load skin library order
	libRows, libErr := db.Query("SELECT SKIN_ID, SKIN_DESC FROM tsekap_lib_skin_extremities ORDER BY SORT_NO, SKIN_ID")
	var skinOrder []struct{ ID, Desc string }
	if libErr == nil {
		defer libRows.Close()
		for libRows.Next() {
			var id, desc string
			libRows.Scan(&id, &desc)
			skinOrder = append(skinOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved skin bits
	var saved string
	var others sql.NullString
	db.QueryRow("SELECT skin_code, others_text FROM patient_pe_skin WHERE patno = ?", patno).Scan(&saved, &others)
	skinBits := []string{}
	if saved != "" {
		skinBits = strings.Split(saved, "|")
	}

	// Build skin finding rows using keys that match frontend ('skin_<ID>')
	for i, s := range skinOrder {
		isChecked := false
		if i < len(skinBits) {
			isChecked = skinBits[i] == "1"
		}
		code := "skin_" + s.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "skin",
			FindingCode: code,
			FindingDesc: s.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	// If there are saved 'others' text, append an 'others' row
	if others.Valid && others.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "skin",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  others.String,
		})
	}

	// Now expand Genitourinary summary into finding rows (if present)
	guLibRows, guErr := db.Query("SELECT GU_ID, GU_DESC FROM tsekap_lib_genitourinary ORDER BY SORT_NO, GU_ID")
	var guOrder []struct{ ID, Desc string }
	if guErr == nil {
		defer guLibRows.Close()
		for guLibRows.Next() {
			var id, desc string
			guLibRows.Scan(&id, &desc)
			guOrder = append(guOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved genitourinary bits
	var guSaved string
	var guOthers sql.NullString
	db.QueryRow("SELECT gu_code, others_text FROM patient_pe_genitourinary WHERE patno = ?", patno).Scan(&guSaved, &guOthers)
	guBits := []string{}
	if guSaved != "" {
		guBits = strings.Split(guSaved, "|")
	}

	for i, g := range guOrder {
		isChecked := false
		if i < len(guBits) {
			isChecked = guBits[i] == "1"
		}
		code := "genitourinary_" + g.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "genitourinary",
			FindingCode: code,
			FindingDesc: g.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if guOthers.Valid && guOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "genitourinary",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  guOthers.String,
		})
	}

	// Now expand Digital Rectal summary into finding rows (if present)
	drLibRows, drErr := db.Query("SELECT RECTAL_ID, RECTAL_DESC FROM tsekap_lib_digital_rectal ORDER BY SORT_NO, RECTAL_ID")
	var drOrder []struct{ ID, Desc string }
	if drErr == nil {
		defer drLibRows.Close()
		for drLibRows.Next() {
			var id, desc string
			drLibRows.Scan(&id, &desc)
			drOrder = append(drOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved digital rectal bits
	var drSaved string
	var drOthers sql.NullString
	db.QueryRow("SELECT dr_code, others_text FROM patient_pe_digital_rectal WHERE patno = ?", patno).Scan(&drSaved, &drOthers)
	drBits := []string{}
	if drSaved != "" {
		drBits = strings.Split(drSaved, "|")
	}

	for i, d := range drOrder {
		isChecked := false
		if i < len(drBits) {
			isChecked = drBits[i] == "1"
		}
		code := "digitalRectal_" + d.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "digitalRectal",
			FindingCode: code,
			FindingDesc: d.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if drOthers.Valid && drOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "digitalRectal",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  drOthers.String,
		})
	}

	// Now expand Neuro/Neurological summary into finding rows (if present)
	neuroLibRows, neuroErr := db.Query("SELECT NEURO_ID, NEURO_DESC FROM tsekap_lib_neuro ORDER BY SORT_NO, NEURO_ID")
	var neuroOrder []struct{ ID, Desc string }
	if neuroErr == nil {
		defer neuroLibRows.Close()
		for neuroLibRows.Next() {
			var id, desc string
			neuroLibRows.Scan(&id, &desc)
			neuroOrder = append(neuroOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved neuro bits
	var neuroSaved string
	var neuroOthers sql.NullString
	db.QueryRow("SELECT neuro_code, others_text FROM patient_pe_neuro WHERE patno = ?", patno).Scan(&neuroSaved, &neuroOthers)
	neuroBits := []string{}
	if neuroSaved != "" {
		neuroBits = strings.Split(neuroSaved, "|")
	}

	for i, n := range neuroOrder {
		isChecked := false
		if i < len(neuroBits) {
			isChecked = neuroBits[i] == "1"
		}
		code := "neuro_" + n.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "neurological",
			FindingCode: code,
			FindingDesc: n.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if neuroOthers.Valid && neuroOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "neurological",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  neuroOthers.String,
		})
	}

	// Now expand CHEST summary into finding rows (if present)
	chestLibRows, chestErr := db.Query("SELECT CHEST_ID, CHEST_DESC FROM tsekap_lib_chest ORDER BY SORT_NO, CHEST_ID")
	var chestOrder []struct{ ID, Desc string }
	if chestErr == nil {
		defer chestLibRows.Close()
		for chestLibRows.Next() {
			var id, desc string
			chestLibRows.Scan(&id, &desc)
			chestOrder = append(chestOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved chest bits
	var chestSaved string
	var chestOthers sql.NullString
	db.QueryRow("SELECT chest_code, others_text FROM patient_pe_chest WHERE patno = ?", patno).Scan(&chestSaved, &chestOthers)
	chestBits := []string{}
	if chestSaved != "" {
		chestBits = strings.Split(chestSaved, "|")
	}

	for i, c := range chestOrder {
		isChecked := false
		if i < len(chestBits) {
			isChecked = chestBits[i] == "1"
		}
		code := "chest_" + c.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "chest",
			FindingCode: code,
			FindingDesc: c.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if chestOthers.Valid && chestOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "chest",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  chestOthers.String,
		})
	}

	// Now expand HEART summary into finding rows (if present)
	heartLibRows, heartErr := db.Query("SELECT HEART_ID, HEART_DESC FROM tsekap_lib_heart ORDER BY SORT_NO, HEART_ID")
	var heartOrder []struct{ ID, Desc string }
	if heartErr == nil {
		defer heartLibRows.Close()
		for heartLibRows.Next() {
			var id, desc string
			heartLibRows.Scan(&id, &desc)
			heartOrder = append(heartOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved heart bits
	var heartSaved string
	var heartOthers sql.NullString
	db.QueryRow("SELECT heart_code, others_text FROM patient_pe_heart WHERE patno = ?", patno).Scan(&heartSaved, &heartOthers)
	heartBits := []string{}
	if heartSaved != "" {
		heartBits = strings.Split(heartSaved, "|")
	}

	for i, h := range heartOrder {
		isChecked := false
		if i < len(heartBits) {
			isChecked = heartBits[i] == "1"
		}
		code := "heart_" + h.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "heart",
			FindingCode: code,
			FindingDesc: h.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if heartOthers.Valid && heartOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "heart",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  heartOthers.String,
		})
	}

	// Now expand ABDOMEN summary into finding rows (if present)
	abdomenLibRows, abdomenErr := db.Query("SELECT ABDOMEN_ID, ABDOMEN_DESC FROM tsekap_lib_abdomen ORDER BY SORT_NO, ABDOMEN_ID")
	var abdomenOrder []struct{ ID, Desc string }
	if abdomenErr == nil {
		defer abdomenLibRows.Close()
		for abdomenLibRows.Next() {
			var id, desc string
			abdomenLibRows.Scan(&id, &desc)
			abdomenOrder = append(abdomenOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved abdomen bits
	var abdomenSaved string
	var abdomenOthers sql.NullString
	db.QueryRow("SELECT abdomen_code, others_text FROM patient_pe_abdomen WHERE patno = ?", patno).Scan(&abdomenSaved, &abdomenOthers)
	abdomenBits := []string{}
	if abdomenSaved != "" {
		abdomenBits = strings.Split(abdomenSaved, "|")
	}

	for i, a := range abdomenOrder {
		isChecked := false
		if i < len(abdomenBits) {
			isChecked = abdomenBits[i] == "1"
		}
		code := "abdomen_" + a.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "abdomen",
			FindingCode: code,
			FindingDesc: a.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if abdomenOthers.Valid && abdomenOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "abdomen",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  abdomenOthers.String,
		})
	}

	// Now expand HEENT summary into finding rows (if present)
	heentLibRows, heentErr := db.Query("SELECT HEENT_ID, HEENT_DESC FROM tsekap_lib_heent ORDER BY SORT_NO, HEENT_ID")
	var heentOrder []struct{ ID, Desc string }
	if heentErr == nil {
		defer heentLibRows.Close()
		for heentLibRows.Next() {
			var id, desc string
			heentLibRows.Scan(&id, &desc)
			heentOrder = append(heentOrder, struct{ ID, Desc string }{ID: id, Desc: desc})
		}
	}

	// Read saved heent bits
	var heentSaved string
	var heentOthers sql.NullString
	db.QueryRow("SELECT heent_code, others_text FROM patient_pe_heent WHERE patno = ?", patno).Scan(&heentSaved, &heentOthers)
	heentBits := []string{}
	if heentSaved != "" {
		heentBits = strings.Split(heentSaved, "|")
	}

	for i, h := range heentOrder {
		isChecked := false
		if i < len(heentBits) {
			isChecked = heentBits[i] == "1"
		}
		code := "heent_" + h.ID
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "heent",
			FindingCode: code,
			FindingDesc: h.Desc,
			IsChecked:   isChecked,
			OthersText:  "",
		})
	}

	if heentOthers.Valid && heentOthers.String != "" {
		list = append(list, FindingRow{
			ID:          0,
			PatientID:   0,
			Category:    "heent",
			FindingCode: "others",
			FindingDesc: "Others",
			IsChecked:   true,
			OthersText:  heentOthers.String,
		})
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

	// Special-case: if saving skin category, store as a summary bit-string
	if payload.Category == "skin" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild skin library order
		libRows, err := db.Query("SELECT SKIN_ID FROM tsekap_lib_skin_extremities ORDER BY SORT_NO, SKIN_ID")
		if err != nil {
			log.Println("skin lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'skin_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "skin_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		skinCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_skin
		_, execErr := db.Exec(`INSERT INTO patient_pe_skin (patno, skin_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE skin_code = VALUES(skin_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, skinCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (skin) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving HEART category, store as a summary bit-string
	if payload.Category == "heart" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild heart library order
		libRows, err := db.Query("SELECT HEART_ID FROM tsekap_lib_heart ORDER BY SORT_NO, HEART_ID")
		if err != nil {
			log.Println("heart lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'heart_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "heart_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		heartCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_heart
		_, execErr := db.Exec(`INSERT INTO patient_pe_heart (patno, heart_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE heart_code = VALUES(heart_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, heartCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (heart) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving ABDOMEN category, store as a summary bit-string
	if payload.Category == "abdomen" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild abdomen library order
		libRows, err := db.Query("SELECT ABDOMEN_ID FROM tsekap_lib_abdomen ORDER BY SORT_NO, ABDOMEN_ID")
		if err != nil {
			log.Println("abdomen lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'abdomen_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "abdomen_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		abdomenCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_abdomen
		_, execErr := db.Exec(`INSERT INTO patient_pe_abdomen (patno, abdomen_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE abdomen_code = VALUES(abdomen_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, abdomenCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (abdomen) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving CHEST category, store as a summary bit-string
	if payload.Category == "chest" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild chest library order
		libRows, err := db.Query("SELECT CHEST_ID FROM tsekap_lib_chest ORDER BY SORT_NO, CHEST_ID")
		if err != nil {
			log.Println("chest lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'chest_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "chest_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		chestCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_chest
		_, execErr := db.Exec(`INSERT INTO patient_pe_chest (patno, chest_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE chest_code = VALUES(chest_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, chestCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (chest) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving HEENT category, store as a summary bit-string
	if payload.Category == "heent" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild heent library order
		libRows, err := db.Query("SELECT HEENT_ID FROM tsekap_lib_heent ORDER BY SORT_NO, HEENT_ID")
		if err != nil {
			log.Println("heent lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'heent_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "heent_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		heentCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_heent
		_, execErr := db.Exec(`INSERT INTO patient_pe_heent (patno, heent_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE heent_code = VALUES(heent_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, heentCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (heent) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving genitourinary category, store as a summary bit-string
	if payload.Category == "genitourinary" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild genitourinary library order
		libRows, err := db.Query("SELECT GU_ID FROM tsekap_lib_genitourinary ORDER BY SORT_NO, GU_ID")
		if err != nil {
			log.Println("genitourinary lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'genitourinary_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "genitourinary_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		guCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_genitourinary
		_, execErr := db.Exec(`INSERT INTO patient_pe_genitourinary (patno, gu_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE gu_code = VALUES(gu_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, guCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (genitourinary) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving digitalRectal category, store as a summary bit-string
	if payload.Category == "digitalRectal" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild digital rectal library order
		libRows, err := db.Query("SELECT RECTAL_ID FROM tsekap_lib_digital_rectal ORDER BY SORT_NO, RECTAL_ID")
		if err != nil {
			log.Println("digital rectal lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'digitalRectal_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "digitalRectal_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		drCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_digital_rectal
		_, execErr := db.Exec(`INSERT INTO patient_pe_digital_rectal (patno, dr_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE dr_code = VALUES(dr_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, drCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (digitalRectal) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Special-case: if saving neurological category, store as a summary bit-string
	if payload.Category == "neurological" {
		// Translate numeric patient ID → case_no (patno)
		var patno string
		db.QueryRow("SELECT case_no FROM patients WHERE id = ?", patientID).Scan(&patno)
		if patno == "" {
			patno = patientID
		}

		// Rebuild neuro library order
		libRows, err := db.Query("SELECT NEURO_ID FROM tsekap_lib_neuro ORDER BY SORT_NO, NEURO_ID")
		if err != nil {
			log.Println("neuro lib query error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer libRows.Close()
		var order []string
		for libRows.Next() {
			var id string
			libRows.Scan(&id)
			order = append(order, id)
		}

		// Build lookup map from incoming findings
		checkedMap := map[string]bool{}
		var othersText string
		for _, f := range payload.Findings {
			if f.FindingCode == "others" {
				othersText = f.OthersText
				continue
			}
			// incoming finding codes for library items are expected in the form 'neuro_<ID>'
			checkedMap[f.FindingCode] = f.IsChecked
		}

		// Build bits in library order
		bits := make([]string, len(order))
		for i, id := range order {
			key := "neuro_" + id
			if checkedMap[key] {
				bits[i] = "1"
			} else {
				bits[i] = "0"
			}
		}
		neuroCode := strings.Join(bits, "|")

		// UPSERT into patient_pe_neuro
		_, execErr := db.Exec(`INSERT INTO patient_pe_neuro (patno, neuro_code, others_text, date_added, added_by)
			VALUES (?, ?, ?, NOW(), 'system')
			ON DUPLICATE KEY UPDATE neuro_code = VALUES(neuro_code), others_text = VALUES(others_text), date_added = NOW()`,
			patno, neuroCode, othersText)
		if execErr != nil {
			log.Println("savePhysicalExamFindings (neurological) error:", execErr)
			http.Error(w, execErr.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Saved"})
		return
	}

	// Default behaviour: Use DELETE + INSERT for simplicity and reliability
	_, err := db.Exec(`DELETE FROM tsekap_tbl_prof_pe_findings WHERE patient_id = ? AND category = ?`,
		patientID, payload.Category)
	if err != nil {
		log.Println("Delete error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert all findings with explicit is_checked values (store 1 for checked, 0 for unchecked)
	for _, f := range payload.Findings {
		_, err := db.Exec(`INSERT INTO tsekap_tbl_prof_pe_findings 
			(patient_id, category, finding_code, finding_desc, is_checked, others_text)
			VALUES (?, ?, ?, ?, ?, ?)`,
			patientID, payload.Category, f.FindingCode, f.FindingDesc, f.IsChecked, f.OthersText)
		if err != nil {
			log.Println("Insert finding error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
