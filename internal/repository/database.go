/*
Package repository implements the Repository and Transactor interfaces
using SQLite as the backing store. All SQL dbquery are confined here; no SQL escapes
this package.

Usage:

	sqlDB, err := sql.Open("sqlite3", "./app.db")
	if err != nil { ... }
	defer sqlDB.Close()

	repo, txtor, err := repository.New(sqlDB)
	if err != nil { ... }

	// Inject repo into services; inject txtor into services that need atomicity:
	fountainService := fountain.NewFountainService(repo)
	healthService := health.NewHealthService(dbconn)
*/
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/dbquery"
)

// appdbimpl is the top-level Repository + Transactor. It embeds RepositoryQueryExecutor
// for dbquery and keeps conn solely for connection-level ops (BeginTx).
// No dbquery method is defined on appdbimpl — they all live on RepositoryQueryExecutor.
type appdbimpl struct {
	dbquery.RepositoryQueryExecutor
	conn *sql.DB
}

type trxdbimpl struct {
	*appdbimpl
}

// New returns a Repository and a Transactor both backed by the provided *sql.DB.
// It also ensures the database schema is up to date.
//
// Returns an error if sqlDB is nil.
func New(sqlDB *sql.DB) (Repository, Transactor, error) {
	if sqlDB == nil {
		return nil, nil, errors.New("sql.DB is required")
	}
	if err := initSchema(sqlDB); err != nil {
		return nil, nil, fmt.Errorf("repository.New initSchema: %w", err)
	}
	impl := appdbimpl{
		RepositoryQueryExecutor: *dbquery.NewRepositoryQueryExecutor(sqlDB),
		conn:                    sqlDB,
	}
	trx := trxdbimpl{
		&impl,
	}
	return &impl, &trx, nil
}

// initSchema creates the fountains table if it does not exist yet.
func initSchema(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='fountains'`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS fountains (
			id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			latitude  FLOAT   NOT NULL,
			longitude FLOAT   NOT NULL,
			status    TEXT    NOT NULL
		)`)
		if err != nil {
			return fmt.Errorf("creating fountains table: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("checking fountains table: %w", err)
	}
	return nil
}

func (trx *trxdbimpl) RunInTransaction(ctx context.Context, fn func(context.Context, Repository) error) error {

	tx, err := trx.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// txRepo swaps the dbquery executor to *sql.Tx so all repo methods run inside
	// the same transaction. conn is nil to prevent accidental nesting.
	txRepo := &appdbimpl{
		RepositoryQueryExecutor: *dbquery.NewRepositoryQueryExecutor(tx),
		conn:                    nil,
	}

	if err = fn(ctx, txRepo); err != nil {
		_ = tx.Rollback() // best-effort; original error takes priority
		return err
	}

	return tx.Commit()
}
