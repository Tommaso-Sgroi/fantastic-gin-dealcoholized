package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// CreateFountain inserts a new fountain row and returns the auto-assigned ID.
func (r *RepositoryQueryExecutor) CreateFountain(ctx context.Context, f record.FountainRecord) (uint64, error) {
	res, err := r.ExecContext(ctx,
		`INSERT INTO fountains (id, latitude, longitude, status) VALUES (NULL, ?, ?, ?)`,
		f.Latitude, f.Longitude, f.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateFountain: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateFountain LastInsertId: %w", err)
	}
	return uint64(id), nil
}
