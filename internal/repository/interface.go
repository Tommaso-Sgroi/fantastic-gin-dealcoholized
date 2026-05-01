package repository

import (
	"context"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/record"
)

// Repository is the full set of persistence operations exposed by this package.
type Repository interface {
	// ListFountains returns all fountains from the database.
	ListFountains(ctx context.Context) ([]record.FountainRecord, error)

	// ListFountainsInBoundingBox returns all fountains within the given bounding box.
	// The caller is responsible for post-filtering to a circle if needed.
	ListFountainsInBoundingBox(ctx context.Context, minLat, maxLat, minLng, maxLng float64) ([]record.FountainRecord, error)

	// CreateFountain inserts a new fountain and returns its assigned ID.
	CreateFountain(ctx context.Context, f record.FountainRecord) (uint64, error)

	// UpdateFountain replaces all mutable fields of the fountain identified by f.ID.
	// Returns errs.ErrNotFound when no row with that ID exists.
	UpdateFountain(ctx context.Context, f record.FountainRecord) error

	// DeleteFountain removes the fountain with the given ID.
	// Returns errs.ErrNotFound when no row with that ID exists.
	DeleteFountain(ctx context.Context, id uint64) error
}

// Transactor exposes the ability to run multiple Repository operations atomically.
// Use RunInTransaction instead of manual Commit/Rollback to prevent transaction leaks.
type Transactor interface {
	RunInTransaction(ctx context.Context, fn func(context.Context, Repository) error) error
}
