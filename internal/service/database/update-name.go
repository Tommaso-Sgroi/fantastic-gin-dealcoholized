package database

// UpdateNameByID updates the name for the given id.
func (db *appdbimpl) UpdateNameByID(id int64, name string) error {
	res, err := db.c.Exec(`UPDATE example_table SET name=? WHERE id=?`, name, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
