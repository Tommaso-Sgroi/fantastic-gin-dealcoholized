package dbquery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/errs"
	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// GetNameByID returns the NameRecord for the given id.
// Returns errs.ErrNotFound (wrapped with context) when no row matches.
func (r *RepositoryQueryExecutor) GetNameByID(ctx context.Context, id int64) (record.NameRecord, error) {
	var n string
	err := r.QueryRowContext(ctx, `SELECT name FROM example_table WHERE id=?`, id).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return record.NameRecord{}, fmt.Errorf("GetNameByID id=%d: %w", id, errs.ErrNotFound)
	} else if err != nil {
		return record.NameRecord{}, fmt.Errorf("GetNameByID id=%d: %w", id, err)
	}
	return record.NameRecord{ID: id, Value: n}, nil
}
