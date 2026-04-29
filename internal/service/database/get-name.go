package database

// GetNameByID returns the name for the given id.
func (db *appdbimpl) GetNameByID(id int64) (string, error) {
	var name string
	err := db.c.QueryRow(`SELECT name FROM example_table WHERE id=?`, id).Scan(&name)
	return name, err
}
