package repository

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// repo implements Repository on top of a GORM handle (either the pool or a transaction).
type repo struct {
	gormdb *gorm.DB
}

// New returns a Repository backed by the provided GORM instance.
// The concrete type is unexported: callers only ever see the interface.
func New(gormdb *gorm.DB) (Repository, error) {
	if gormdb == nil {
		return nil, errors.New("repository: gorm.DB is required")
	}
	return repo{gormdb}, nil
}

// RunInTransaction executes fn within a single database transaction.
// If fn returns an error, the transaction is rolled back; otherwise committed.
func (r repo) RunInTransaction(ctx context.Context, fn func(context.Context, Repository) error) error {
	return r.gormdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, repo{tx})
	})
}

// Migrate applies all pending goose migrations found in migrationDir.
func Migrate(gormdb *gorm.DB, migrationDir fs.FS) error {
	if gormdb == nil {
		return errors.New("repository: gorm.DB is required")
	}
	dbconn, err := gormdb.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	goose.SetBaseFS(migrationDir)
	if err := goose.SetDialect(gormdb.Dialector.Name()); err != nil {
		return err
	}
	return goose.Up(dbconn, ".")
}
