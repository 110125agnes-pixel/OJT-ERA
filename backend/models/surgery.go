package models

// Surgery represents a surgery record in the database
type Surgery struct {
	ID          int    `json:"id"`
	PatientName string `json:"patient_name"`
	SurgeryType string `json:"surgery_type"`
	SurgeonName string `json:"surgeon_name"`
	SurgeryDate string `json:"surgery_date"`
	SurgeryTime string `json:"surgery_time"`
	Duration    string `json:"duration"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

// CreateSurgeryTable creates the surgeries table if it doesn't exist
func CreateSurgeryTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS surgeries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			patient_name TEXT NOT NULL,
			surgery_type TEXT NOT NULL,
			surgeon_name TEXT NOT NULL,
			surgery_date TEXT,
			surgery_time TEXT,
			duration TEXT,
			status TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := DB.Exec(query)
	return err
}

// GetAllSurgeries retrieves all surgeries from the database
func GetAllSurgeries() ([]Surgery, error) {
	rows, err := DB.Query(`
		SELECT id, patient_name, surgery_type, surgeon_name, surgery_date, 
		       surgery_time, duration, status, notes 
		FROM surgeries 
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surgeries []Surgery
	for rows.Next() {
		var s Surgery
		err := rows.Scan(
			&s.ID, &s.PatientName, &s.SurgeryType, &s.SurgeonName,
			&s.SurgeryDate, &s.SurgeryTime, &s.Duration, &s.Status, &s.Notes,
		)
		if err != nil {
			return nil, err
		}
		surgeries = append(surgeries, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return surgeries, nil
}

// GetSurgeryByID retrieves a single surgery by ID
func GetSurgeryByID(id int) (*Surgery, error) {
	var s Surgery
	err := DB.QueryRow(`
		SELECT id, patient_name, surgery_type, surgeon_name, surgery_date, 
		       surgery_time, duration, status, notes 
		FROM surgeries 
		WHERE id = ?
	`, id).Scan(
		&s.ID, &s.PatientName, &s.SurgeryType, &s.SurgeonName,
		&s.SurgeryDate, &s.SurgeryTime, &s.Duration, &s.Status, &s.Notes,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSurgery inserts a new surgery into the database
func CreateSurgery(s *Surgery) error {
	result, err := DB.Exec(`
		INSERT INTO surgeries (patient_name, surgery_type, surgeon_name, surgery_date, 
		                      surgery_time, duration, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, s.PatientName, s.SurgeryType, s.SurgeonName, s.SurgeryDate,
		s.SurgeryTime, s.Duration, s.Status, s.Notes)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = int(id)
	return nil
}

// UpdateSurgery updates an existing surgery in the database
func UpdateSurgery(id int, s *Surgery) error {
	_, err := DB.Exec(`
		UPDATE surgeries 
		SET patient_name = ?, surgery_type = ?, surgeon_name = ?, surgery_date = ?,
		    surgery_time = ?, duration = ?, status = ?, notes = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, s.PatientName, s.SurgeryType, s.SurgeonName, s.SurgeryDate,
		s.SurgeryTime, s.Duration, s.Status, s.Notes, id)
	return err
}

// DeleteSurgery deletes a surgery from the database
func DeleteSurgery(id int) error {
	_, err := DB.Exec("DELETE FROM surgeries WHERE id = ?", id)
	return err
}
