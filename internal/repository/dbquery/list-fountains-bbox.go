package dbquery

import (
	"context"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// ListFountainsInBoundingBox returns all fountains whose coordinates fall within the
// given bounding box [minLat, maxLat] × [minLng, maxLng].
// The caller is responsible for further circle-filtering if needed.
func (r *RepositoryQueryExecutor) ListFountainsInBoundingBox(
	ctx context.Context,
	minLat, maxLat, minLng, maxLng float64,
) ([]record.FountainRecord, error) {
	const query = `
SELECT id, latitude, longitude, status
FROM fountains
WHERE ? <= latitude AND latitude <= ? AND ? <= longitude AND longitude <= ?`

	rows, err := r.QueryContext(ctx, query, minLat, maxLat, minLng, maxLng)
	if err != nil {
		return nil, fmt.Errorf("ListFountainsInBoundingBox: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var ret []record.FountainRecord
	for rows.Next() {
		var f record.FountainRecord
		if err = rows.Scan(&f.ID, &f.Latitude, &f.Longitude, &f.Status); err != nil {
			return nil, fmt.Errorf("ListFountainsInBoundingBox scan: %w", err)
		}
		ret = append(ret, f)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ListFountainsInBoundingBox rows: %w", err)
	}
	return ret, nil
}
