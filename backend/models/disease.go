package models

// Disease represents a medical condition/disease
type Disease struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Barcode  string `json:"barcode"`
	Category string `json:"category"`
}

// EmployeeDisease represents the relationship between an employee and a disease
type EmployeeDisease struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	DiseaseID  int    `json:"disease_id"`
	DiseaseCode string `json:"disease_code"`
	DiseaseName string `json:"disease_name"`
	DateDiagnosed string `json:"date_diagnosed"`
}

// CreateDiseaseTable creates the disease table if it doesn't exist
func CreateDiseaseTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS diseases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE,
			barcode TEXT NOT NULL UNIQUE,
			category TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	// Create employee_disease junction table
	query2 := `
		CREATE TABLE IF NOT EXISTS employee_diseases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			employee_id INTEGER NOT NULL,
			disease_id INTEGER NOT NULL,
			date_diagnosed TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (employee_id) REFERENCES items(id),
			FOREIGN KEY (disease_id) REFERENCES diseases(id),
			UNIQUE(employee_id, disease_id)
		)
	`
	_, err = DB.Exec(query2)
	return err
}

// GetAllDiseases retrieves all diseases from the database
func GetAllDiseases() ([]Disease, error) {
	rows, err := DB.Query("SELECT id, name, code, barcode, category FROM diseases ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	diseases := []Disease{}
	for rows.Next() {
		var disease Disease
		err := rows.Scan(&disease.ID, &disease.Name, &disease.Code, &disease.Barcode, &disease.Category)
		if err != nil {
			return nil, err
		}
		diseases = append(diseases, disease)
	}
	return diseases, nil
}

// CreateDisease inserts a new disease into the database
func CreateDisease(disease *Disease) error {
	result, err := DB.Exec(
		"INSERT INTO diseases (name, code, barcode, category) VALUES (?, ?, ?, ?)",
		disease.Name, disease.Code, disease.Barcode, disease.Category,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	disease.ID = int(id)
	return err
}

// UpdateDisease updates an existing disease in the database
func UpdateDisease(id int, disease *Disease) (int64, error) {
	result, err := DB.Exec(
		"UPDATE diseases SET name = ?, code = ?, barcode = ?, category = ? WHERE id = ?",
		disease.Name, disease.Code, disease.Barcode, disease.Category, id,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// DeleteDisease removes a disease from the database
func DeleteDisease(id int) (int64, error) {
	result, err := DB.Exec("DELETE FROM diseases WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// GetEmployeeDiseases gets all diseases for a specific employee
func GetEmployeeDiseases(employeeID int) ([]EmployeeDisease, error) {
	rows, err := DB.Query(`
		SELECT ed.id, ed.employee_id, ed.disease_id, d.code, d.name, ed.date_diagnosed
		FROM employee_diseases ed
		JOIN diseases d ON ed.disease_id = d.id
		WHERE ed.employee_id = ?
		ORDER BY ed.id DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employeeDiseases := []EmployeeDisease{}
	for rows.Next() {
		var ed EmployeeDisease
		err := rows.Scan(&ed.ID, &ed.EmployeeID, &ed.DiseaseID, &ed.DiseaseCode, &ed.DiseaseName, &ed.DateDiagnosed)
		if err != nil {
			return nil, err
		}
		employeeDiseases = append(employeeDiseases, ed)
	}
	return employeeDiseases, nil
}

// AddEmployeeDisease adds a disease to an employee
func AddEmployeeDisease(employeeID int, diseaseID int, dateDiagnosed string) (int64, error) {
	result, err := DB.Exec(
		"INSERT INTO employee_diseases (employee_id, disease_id, date_diagnosed) VALUES (?, ?, ?)",
		employeeID, diseaseID, dateDiagnosed,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

// RemoveEmployeeDisease removes a disease from an employee
func RemoveEmployeeDisease(employeeID int, diseaseID int) (int64, error) {
	result, err := DB.Exec(
		"DELETE FROM employee_diseases WHERE employee_id = ? AND disease_id = ?",
		employeeID, diseaseID,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}
