package record

// NameRecord is the raw data bag produced by name dbquery.
// db tags drive automatic struct scanning (e.g. via sqlx).
// It has no domain logic — the domain layer maps it into its own entities.
type NameRecord struct {
	ID    int64  `db:"id"`
	Value string `db:"name"`
}
