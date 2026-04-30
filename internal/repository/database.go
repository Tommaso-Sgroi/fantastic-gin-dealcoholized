/*
Package repository implements the domain.NameRepository and domain.Pinger interfaces
using SQLite as the backing store. All SQL queries are confined here; no SQL escapes
this package.

Usage:

	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil { ... }
	defer db.Close()

	repo, err := repository.New(db)
	if err != nil { ... }

	// Pass repo to domain services:
	nameService := domain.NewNameService(repo)
	healthService := domain.NewHealthService(repo)
*/
package repository

import (
	"context"
	"database/sql"
	"errors"
)

type appdbimpl struct {
	c *sql.DB
}

// New returns a new Repository backed by the provided *sql.DB.
// Returns an error if db is nil.
func New(db *sql.DB) (Repository, error) {
	if db == nil {
		return nil, errors.New("sql.DB is required")
	}
	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping(ctx context.Context) error {
	return db.c.PingContext(ctx)
}
