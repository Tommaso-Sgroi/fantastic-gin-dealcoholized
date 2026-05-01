/*
Package repository implements the Repository and Transactor interfaces
using SQLite as the backing store. All SQL queries are confined here; no SQL escapes
this package.

Usage:

	sqlDB, err := sql.Open("sqlite3", "./app.db")
	if err != nil { ... }
	defer sqlDB.Close()

	repo, txtor, err := repository.New(sqlDB)
	if err != nil { ... }

	// Inject repo into services; inject txtor into services that need atomicity:
	nameService := name.NewNameService(repo, txtor)
	healthService := health.NewHealthService(dbconn)
*/
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/repository/queries"
)

// appdbimpl is the top-level Repository + Transactor. It embeds RepositoryQueryExecutor
// for queries and keeps conn solely for connection-level ops (BeginTx).
// No query method is defined on appdbimpl — they all live on RepositoryQueryExecutor.
type appdbimpl struct {
	queries.RepositoryQueryExecutor
	conn *sql.DB
}

type trxdbimpl struct {
	*appdbimpl
}

// New returns a Repository and a Transactor both backed by the provided *sql.DB.
//
// Returns an error if sqlDB is nil.
func New(sqlDB *sql.DB) (Repository, Transactor, error) {
	if sqlDB == nil {
		return nil, nil, errors.New("sql.DB is required")
	}
	impl := appdbimpl{
		RepositoryQueryExecutor: *queries.NewRepositoryQueryExecutor(sqlDB),
		conn:                    sqlDB,
	}
	trx := trxdbimpl{
		&impl,
	}
	return &impl, &trx, nil
}

func (trx *trxdbimpl) RunInTransaction(ctx context.Context, fn func(context.Context, Repository) error) error {

	tx, err := trx.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// txRepo swaps the query executor to *sql.Tx so all repo methods run inside
	// the same transaction. inTx is set to prevent accidental nesting.
	txRepo := &appdbimpl{
		RepositoryQueryExecutor: *queries.NewRepositoryQueryExecutor(tx),
		conn:                    nil, // conn is not needed inside a transaction; all queries go through the txRepo's RepositoryQueryExecutor
	}

	if err = fn(ctx, txRepo); err != nil {
		_ = tx.Rollback() // best-effort; original error takes priority
		return err
	}

	return tx.Commit()
}
