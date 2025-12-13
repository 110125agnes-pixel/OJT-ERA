package models

import "database/sql"

// Item represents an item in the database
type Item struct {
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

var DB *sql.DB

// InitDB sets the database connection
func InitDB(db *sql.DB) {
	DB = db
}

// CreateTable creates the items table if it doesn't exist
func CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			case_no VARCHAR(50),
			hospital_no VARCHAR(50),
			lastname VARCHAR(100) NOT NULL,
			firstname VARCHAR(100) NOT NULL,
			middlename VARCHAR(100),
			suffix VARCHAR(20),
			birthdate DATE,
			age VARCHAR(10),
			room VARCHAR(50),
			admission_date DATETIME,
			discharge_date DATETIME,
			sex VARCHAR(20),
			height VARCHAR(20),
			weight VARCHAR(20),
			complaint TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := DB.Exec(query)
	return err
}

// GetAllItems retrieves all items from the database
func GetAllItems() ([]Item, error) {
	rows, err := DB.Query("SELECT id, case_no, hospital_no, lastname, firstname, middlename, suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint FROM items ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.CaseNo, &item.HospitalNo, &item.Lastname, &item.Firstname, &item.Middlename, &item.Suffix, &item.Birthdate, &item.Age, &item.Room, &item.AdmissionDate, &item.DischargeDate, &item.Sex, &item.Height, &item.Weight, &item.Complaint)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetItemByID retrieves a single item by ID
func GetItemByID(id int) (*Item, error) {
	var item Item
	err := DB.QueryRow("SELECT id, case_no, hospital_no, lastname, firstname, middlename, suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint FROM items WHERE id = ?", id).
		Scan(&item.ID, &item.CaseNo, &item.HospitalNo, &item.Lastname, &item.Firstname, &item.Middlename, &item.Suffix, &item.Birthdate, &item.Age, &item.Room, &item.AdmissionDate, &item.DischargeDate, &item.Sex, &item.Height, &item.Weight, &item.Complaint)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem inserts a new item into the database
func CreateItem(item *Item) error {
	result, err := DB.Exec(
		"INSERT INTO items (case_no, hospital_no, lastname, firstname, middlename, suffix, birthdate, age, room, admission_date, discharge_date, sex, height, weight, complaint) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		item.CaseNo, item.HospitalNo, item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Age, item.Room, item.AdmissionDate, item.DischargeDate, item.Sex, item.Height, item.Weight, item.Complaint,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = int(id)
	return nil
}

// UpdateItem updates an existing item in the database
func UpdateItem(id int, item *Item) (int64, error) {
	result, err := DB.Exec(
		"UPDATE items SET case_no = ?, hospital_no = ?, lastname = ?, firstname = ?, middlename = ?, suffix = ?, birthdate = ?, age = ?, room = ?, admission_date = ?, discharge_date = ?, sex = ?, height = ?, weight = ?, complaint = ? WHERE id = ?",
		item.CaseNo, item.HospitalNo, item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Age, item.Room, item.AdmissionDate, item.DischargeDate, item.Sex, item.Height, item.Weight, item.Complaint, id,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// DeleteItem removes an item from the database
func DeleteItem(id int) (int64, error) {
	result, err := DB.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}
