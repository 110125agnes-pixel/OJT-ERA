package models

import "database/sql"

// Item represents an item in the database
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

var DB *sql.DB

// InitDB sets the database connection
func InitDB(db *sql.DB) {
	DB = db
}

// CreateTable creates the items table if it doesn't exist
func CreateTable() error {
	query := `
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
	_, err := DB.Exec(query)
	return err
}

// GetAllItems retrieves all items from the database
func GetAllItems() ([]Item, error) {
	rows, err := DB.Query("SELECT id, lastname, firstname, middlename, suffix, birthdate, sex, civil_status FROM items ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Lastname, &item.Firstname, &item.Middlename, &item.Suffix, &item.Birthdate, &item.Sex, &item.CivilStatus)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// CreateItem inserts a new item into the database
func CreateItem(item *Item) error {
	result, err := DB.Exec(
		"INSERT INTO items (lastname, firstname, middlename, suffix, birthdate, sex, civil_status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Sex, item.CivilStatus,
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
		"UPDATE items SET lastname = ?, firstname = ?, middlename = ?, suffix = ?, birthdate = ?, sex = ?, civil_status = ? WHERE id = ?",
		item.Lastname, item.Firstname, item.Middlename, item.Suffix, item.Birthdate, item.Sex, item.CivilStatus, id,
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
