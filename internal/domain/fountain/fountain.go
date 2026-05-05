/*
Package fountain contains the core business entities and rules for the Fountain domain.
It has no dependency on any framework or infrastructure concern (no Gin, no SQL, no HTTP).
*/
package fountain

import (
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

const (
	StatusGood   = "good"
	StatusFaulty = "faulty"
)

// Fountain represents a drinking-water fountain ("nasone") in the system.
// json tags are here for HTTP serialisation only — no db tags, no infra concerns.
type Fountain struct {
	ID        uint64  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

// Validate checks domain-level invariants on a Fountain value.
// It does NOT check the ID (callers supply the ID from the URL path, not the body).
func (f Fountain) Validate() error {
	if f.Latitude < -90 || f.Latitude > 90 {
		return fmt.Errorf("latitude %f out of range [-90, 90]: %w", f.Latitude, errs.ErrInvalidInput)
	}
	if f.Longitude < -180 || f.Longitude > 180 {
		return fmt.Errorf("longitude %f out of range [-180, 180]: %w", f.Longitude, errs.ErrInvalidInput)
	}
	if f.Status != StatusGood && f.Status != StatusFaulty {
		return fmt.Errorf("status %q must be %q or %q: %w", f.Status, StatusGood, StatusFaulty, errs.ErrInvalidInput)
	}
	return nil
}

// fountainFromRow maps a record.FountainRecord (DB-boundary DTO) to a domain Fountain.
// All record-to-entity conversions go through here; adding a field to Fountain means
// updating this function — the compiler will catch any missing assignments.
func fountainFromRow(row record.FountainRecord) Fountain {
	return Fountain{
		ID:        row.ID,
		Latitude:  row.Latitude,
		Longitude: row.Longitude,
		Status:    row.Status,
	}
}
