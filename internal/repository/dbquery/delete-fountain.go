package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
)

// DeleteFountain removes the fountain with the given ID.
// Returns errs.ErrNotFound when no matching row exists.
func (r *RepositoryQueryExecutor) DeleteFountain(ctx context.Context, id uint64) error {
	res, err := r.ExecContext(ctx, `DELETE FROM fountains WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("DeleteFountain id=%d: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteFountain RowsAffected id=%d: %w", id, err)
	}
	if affected == 0 {
		return fmt.Errorf("DeleteFountain id=%d: %w", id, errs.ErrNotFound)
	}
	return nil
}
