package database

// DeleteNameByID deletes the row with the given id.
// Returns ErrNotFound if no row matched.
func (db *appdbimpl) DeleteNameByID(id int64) error {
	res, err := db.c.Exec(`DELETE FROM example_table WHERE id=?`, id)
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
