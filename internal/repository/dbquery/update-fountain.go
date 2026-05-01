package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// UpdateFountain replaces the latitude, longitude and status of the fountain
// identified by f.ID. Returns errs.ErrNotFound when no matching row exists.
func (r *RepositoryQueryExecutor) UpdateFountain(ctx context.Context, f record.FountainRecord) error {
	res, err := r.ExecContext(ctx,
		`UPDATE fountains SET latitude=?, longitude=?, status=? WHERE id=?`,
		f.Latitude, f.Longitude, f.Status, f.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateFountain id=%d: %w", f.ID, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateFountain RowsAffected id=%d: %w", f.ID, err)
	}
	if affected == 0 {
		return fmt.Errorf("UpdateFountain id=%d: %w", f.ID, errs.ErrNotFound)
	}
	return nil
}
