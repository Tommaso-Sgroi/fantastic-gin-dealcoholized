# Plan: Fix Repository Transaction Safety Issues

The current implementation has 4 concrete issues: `New` hides the `Transactor`, `WithTransaction` is a dangerous manual-commit stub, nested transactions silently misbehave, and `Ping` is dead code leaking infra concerns into the wrong layer. All fixes are self-contained within the `repository` package and `main.go`.

## Steps

1. **Remove `WithTransaction`/`TransactionRepository` from [`interface.go`](internal/repository/interface.go)** — the manual `Commit`/`Rollback` API is error-prone (leaks if caller forgets) and a stub returning `nil`. `Transactor.RunInTransaction` (callback pattern) already covers the use case better and safely.

2. **Change `New` in [`database.go`](internal/repository/database.go) to return `(Repository, Transactor, error)`** — `appdbimpl` already satisfies both; callers currently must type-assert to get `Transactor`, which is fragile. Update the doc comment to match.

3. **Add `inTx bool` to `appdbimpl` and guard `RunInTransaction` in [`database.go`](internal/repository/database.go)** — if `inTx` is `true`, return `errors.New("nested transactions are not supported")` immediately. Set `inTx: true` on the `txRepo` passed into `fn` so the guard propagates naturally.

4. **Remove the dead `Ping(ctx context.Context) error` method from [`database.go`](internal/repository/database.go)** — `Repository` no longer declares `Ping`, and `health.Pinger` requires `Ping() error` (no `ctx`). Since `main.go` already passes `dbconn` (`*sql.DB`) directly to `health.NewHealthService`, this method serves no interface and is dead code.

5. **Update [`main.go`](cmd/webapi/main.go)** — unpack the new 3-value return from `New`; wire `Transactor` into any service that needs atomicity (e.g. `name.NewNameService(db, txtor)`).

## Further Considerations

1. **`Pinger` via repo or raw `*sql.DB`?** `main.go` currently passes `dbconn` directly to `health.NewHealthService` — acceptable, but it bypasses the repo layer. Alternative: add `Ping() error` (no ctx) to `appdbimpl` so `db` satisfies `health.Pinger` and `dbconn` can stop leaking into the wiring layer. Low priority but cleaner.
2. **Savepoints for future nested tx support?** The `inTx` guard returns an error for now; SQLite supports `SAVEPOINT` for true nesting. Worth noting if the codebase grows to need it.

