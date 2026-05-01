package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
)

// DeleteNameByID deletes the row with the given id.
// Returns errs.ErrNotFound if no row matched.
func (r *RepositoryQueryExecutor) DeleteNameByID(ctx context.Context, id int64) error {
	res, err := r.ExecContext(ctx, `DELETE FROM example_table WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("DeleteNameByID id=%d: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteNameByID RowsAffected id=%d: %w", id, err)
	} else if affected == 0 {
		return fmt.Errorf("DeleteNameByID id=%d: %w", id, errs.ErrNotFound)
	}
	return nil
}
