package repository

import (
	"context"
)

// Repository is the full set of persistence operations exposed by this package.
type Repository interface {
	GetNameByID(ctx context.Context, id int64) (string, error)
	UpdateNameByID(ctx context.Context, id int64, name string) error
	CreateName(ctx context.Context, name string) (int64, error)
	DeleteNameByID(ctx context.Context, id int64) error
}

// Transactor exposes the ability to run multiple Repository operations atomically.
// Use RunInTransaction instead of manual Commit/Rollback to prevent transaction leaks.
type Transactor interface {
	RunInTransaction(ctx context.Context, fn func(context.Context, Repository) error) error
}
