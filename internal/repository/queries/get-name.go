package queries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
)

// GetNameByID returns the name for the given id.
// Returns errs.ErrNotFound (wrapped with context) when no row matches.
func (r *RepositoryQueryExecutor) GetNameByID(ctx context.Context, id int64) (string, error) {
	var n string
	err := r.QueryRowContext(ctx, `SELECT name FROM example_table WHERE id=?`, id).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("GetNameByID id=%d: %w", id, errs.ErrNotFound)
	} else if err != nil {
		return "", fmt.Errorf("GetNameByID id=%d: %w", id, err)
	}
	return n, nil
}
