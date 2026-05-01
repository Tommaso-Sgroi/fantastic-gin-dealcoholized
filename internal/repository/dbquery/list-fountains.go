package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// ListFountains returns all rows from the fountains table.
func (r *RepositoryQueryExecutor) ListFountains(ctx context.Context) ([]record.FountainRecord, error) {
	rows, err := r.QueryContext(ctx, `SELECT id, latitude, longitude, status FROM fountains`)
	if err != nil {
		return nil, fmt.Errorf("ListFountains: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var ret []record.FountainRecord
	for rows.Next() {
		var f record.FountainRecord
		if err = rows.Scan(&f.ID, &f.Latitude, &f.Longitude, &f.Status); err != nil {
			return nil, fmt.Errorf("ListFountains scan: %w", err)
		}
		ret = append(ret, f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ListFountains rows: %w", err)
	}
	return ret, nil
}
