package queries

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
)

// UpdateNameByID updates the name for the given id.
// Returns errs.ErrNotFound if no row matched.
func (r *RepositoryQueryExecutor) UpdateNameByID(ctx context.Context, id int64, n string) error {
	res, err := r.ExecContext(ctx, `UPDATE example_table SET name=? WHERE id=?`, n, id)
	if err != nil {
		return fmt.Errorf("UpdateNameByID id=%d: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateNameByID RowsAffected id=%d: %w", id, err)
	} else if affected == 0 {
		return fmt.Errorf("UpdateNameByID id=%d: %w", id, errs.ErrNotFound)
	}
	return nil
}
