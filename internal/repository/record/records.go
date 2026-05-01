package record

// FountainRecord is the raw data bag produced by fountain SQL queries.
// db tags drive automatic column scanning.
// No domain logic lives here — the domain layer maps it into its own entities.
type FountainRecord struct {
	ID        uint64  `db:"id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
	Status    string  `db:"status"`
}
