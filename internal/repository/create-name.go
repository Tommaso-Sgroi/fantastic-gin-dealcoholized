package repository

import "fmt"

// CreateName inserts a new row and returns the generated id.
func (db *appdbimpl) CreateName(name string) (int64, error) {
	res, err := db.c.Exec(`INSERT INTO example_table (name) VALUES (?)`, name)
	if err != nil {
		return 0, fmt.Errorf("CreateName: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateName LastInsertId: %w", err)
	}
	return id, nil
}
