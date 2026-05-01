package queries

import (
	"context"
	"database/sql"
)

// QueryExecutor is the common interface between *sql.DB and *sql.Tx.
// All repository methods use this so they work transparently
// with or without an active transaction.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// RepositoryQueryExecutor is the ONLY type that holds db and therefore the only receiver for query methods.
// It intentionally has no conn field so a developer writing a new query method
// cannot accidentally bypass a transaction by reaching *sql.DB directly.
type RepositoryQueryExecutor struct {
	QueryExecutor // *sql.DB normally, *sql.Tx inside a transaction
}

func NewRepositoryQueryExecutor(qe QueryExecutor) *RepositoryQueryExecutor {
	return new(RepositoryQueryExecutor).WithQueryExecutor(qe)
}

func (r *RepositoryQueryExecutor) WithQueryExecutor(executor QueryExecutor) *RepositoryQueryExecutor {
	r.QueryExecutor = executor
	return r
}
