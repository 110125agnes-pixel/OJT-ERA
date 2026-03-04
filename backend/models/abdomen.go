package models

import (
	"database/sql"
)

type Abdomen struct {
	ID   int    `json:"id"`
	Desc string `json:"ABDOMEN_DESC"`
}

func GetAllAbdomens(db *sql.DB) ([]Abdomen, error) {
	rows, err := db.Query("SELECT ABDOMEN_ID, ABDOMEN_DESC FROM tsekap_lib_abdomen ORDER BY SORT_NO ASC, ABDOMEN_ID ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abdomens []Abdomen
	for rows.Next() {
		var a Abdomen
		if err := rows.Scan(&a.ID, &a.Desc); err != nil {
			return nil, err
		}
		abdomens = append(abdomens, a)
	}
	return abdomens, nil
}
