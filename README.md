# Fantastic Gin Dealcoholized

This repository originally contained the basic structure for [Web and Software Architecture](http://gamificationlab.uniroma1.it/en/wasa/) homework project.
It has been rewritten using the [Gin](https://github.com/gin-gonic/gin) framework and updated dependencies. The backend architecture has been redesigned from scratch following a **layered architecture** and SOLID principles.

This project is a fork of [Fantastic Coffee (decaffeinated)](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated)
by [Enrico Bassetti](https://github.com/Enrico204) - originally created for the
[Web and Software Architecture](http://gamificationlab.uniroma1.it/en/wasa/) course at Sapienza University of Rome.


## What it does

The project itself shows an "actual" example implementation of the backend.  
It exposes a small REST API (documented in [`doc/api.yaml`](doc/api.yaml)) to create, read, update and delete "names" stored in a SQLite database:

| Method   | Path                 | Description                     |
|----------|----------------------|---------------------------------|
| `GET`    | `/healthz`           | Liveness probe                  |
| `POST`   | `/names`             | Create a name                   |
| `GET`    | `/names/:id`         | Get a name by ID                |
| `PUT`    | `/names/:id`         | Update a name by ID             |
| `DELETE` | `/names/:id`         | Delete a name by ID             |
| `POST`   | `/names/:id/replace` | Atomically replace a name by ID |

All error responses are JSON: `{"error": "<message>"}`.

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is a small daemon for checking server health; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` is the main web API server daemon
* `demo/` contains a demo YAML configuration file
* `db/` contains the SQL migration files
* `doc/` contains the OpenAPI 3.1 specification (`api.yaml`)
* `internal/` has all packages for implementing project-specific functionalities
	* `domain/` contains business entities, use-cases, and domain-owned port interfaces (for example the `name` and `health` domain packages)
	* `handler/` contains the HTTP transport layer (Gin router, middleware, request context wrapper, DTOs, and route handlers)
		* `handler/route` contains route handlers and HTTP-level error mapping
		* `handler/middleware` contains Gin middleware for request-scoped metadata and logging
		* `handler/reqcontext` contains the wrapped request context passed to handlers
		* `handler/dto` contains request/response data-transfer objects used at the HTTP boundary
	* `repository/` contains the persistence layer (SQLite adapter implementing the domain ports)
		* `repository/dbquery` contains the raw SQL query implementations
		* `repository/record` contains plain data-bag structs produced by queries, driven by `db` struct tags for automatic scanning
	* `errs/` contains shared application-level sentinel errors used to disambiguate failures across layers
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is a Vue.js web frontend; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

### Architecture

The backend follows a clean **layered architecture** organized in three layers:

* **HTTP layer (`internal/handler`)** receives requests, adapts Gin handlers to route handlers through a wrapped request context, maps DTOs to/from domain models, and maps domain errors to HTTP status codes.
* **Domain layer (`internal/domain`)** contains business rules and use-case orchestration, independent from any framework or infrastructure concern.
* **Repository layer (`internal/repository`)** is the persistence adapter. It implements the `Repository` and `Transactor` interfaces using SQLite, assembles `record.*` structs from raw query results, and manages transaction lifecycle.

Typical request flow:

```
HTTP request → Gin router/middleware → handler/route → domain use-case → repository → SQLite
```

Dependency direction:

```
handler → domain → repository
               ↑
             errs (leaf — no internal imports)
```

Dependencies flow strictly top-to-bottom. Shared error sentinels live in `internal/errs` so `errors.Is(...)` works consistently through wrapped errors. The HTTP contract is kept explicit via `internal/handler/dto`, domain entities stay in `internal/domain`, and raw DB projections stay in `internal/repository/record`.

### Notable files

| File | Role                                                                                                                                                                                                                                       |
|------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `cmd/webapi/main.go` | Entry point: loads config, wires all dependencies (DB → repository → domain services → HTTP router), starts the server and handles graceful shutdown on `SIGTERM`                                                                          |
| `internal/errs/errors.go` | Leaf package of sentinel errors (`ErrNotFound`, `ErrInvalidInput`, …) shared across all layers without creating import cycles                                                                                                              |
| `internal/repository/database.go` | Core repository implementation: `appdbimpl` (implements `Repository`) and `trxdbimpl` (implements `Transactor`). `RunInTransaction` swaps the SQL executor to `*sql.Tx` so all query methods run inside the same transaction transparently |
| `internal/repository/interface.go` | `Repository` and `Transactor` interface declarations — defined here and consumed by the domain layer                                                                                                                                       |
| `internal/repository/record/records.go` | Plain data-bag structs produced by SQL queries. `db` struct tags enable automatic column scanning (e.g. via `sqlx`). No domain logic lives here                                                                                            |
| `internal/repository/dbquery/query.go` | Defines `QueryExecutor` — the common interface between `*sql.DB` and `*sql.Tx` — and `RepositoryQueryExecutor`, the only type allowed to hold and execute SQL. This prevents any query method from accidentally bypassing a transaction    |
| `internal/domain/<entity>/<entity>-interface.go` | Declares the public use-case interface for an entity (e.g. `ServiceName`). Every domain service must satisfy this interface; it is the contract consumed by the handler layer. Adding a new use-case means adding a method here first      |
| `internal/domain/<entity>/<entity>-service.go` | Concrete instantiation of the service interface. Owns business rules, validation, and use-case orchestration                                                                                                                               |
| `internal/domain/<entity>/<entity>.go` | Domain entities definition and the private `xxxFromRow` mapping function that converts a `record.*` struct into the domain entity. This is the single place to update when an entity grows a field                                         |
| `internal/handler/api.go` | Defines the `Router` interface, `Config`, and `New()`. Configures the Gin engine, registers global middleware (CORS, recovery, request-context logger), and validates required dependencies at startup                                     |
| `internal/handler/api-handler.go` | Registers all routes on the Gin engine via `Handler()` and provides `wrap()` — the adapter that converts a in custom context `routes.HandlerFunc` into a `gin.HandlerFunc`                                                                 |
| `internal/handler/route/handler.go` | Defines `HandlerFunc` (the handler signature) and `WebAPIHandler` (the struct that carries domain service dependencies shared across all route handlers)                                                                |

Other project files include:
* `open-node.sh` starts a temporary container using the `node:20` image for safe web frontend development

## Run the demo server:

```shell
go run ./cmd/webapi/ --config-path demo/config.yml
```

## Configuration

Configuration is loaded in the following priority order (lowest → highest):

```
defaults CLI flags < environment variables < demo/config.yml
```

The YAML config file path defaults to `demo/config.yml` and can be changed with `-c <path>` or `API_CONFIG_PATH=<path>`.
In particular, the `demo/config.yml` **can be deactivated** with an empty value of `-c ""` or `API_CONFIG_PATH=`.

Key settings (with their defaults):

| Setting | Default | Description |
|---|---|---|
| `web.apihost` | `0.0.0.0:3000` | API server bind address |
| `web.debughost` | `0.0.0.0:4000` | pprof debug server address |
| `db.filename` | `db/fantastic-gin.db` | SQLite database file path |
| `logger.level` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `logger.file` | `/tmp/dealch.log` | Log output (`stderr`, `stdout`, or a file path) |

All keys can be overridden via environment variables prefixed with `API_`, e.g.:

```shell
API_WEB_APIHOST=0.0.0.0:3000 API_DB_FILENAME=db/fantastic-gin.db go run ./cmd/webapi/
```

## Database

The app uses **SQLite** via [`go-sqlite3`](https://github.com/mattn/go-sqlite3) (CGo).

## How to build

**Backend only:**

```shell
go build ./cmd/webapi/
```

**With embedded WebUI:**

```shell
./open-node.sh
# (inside the container)
yarn run build-embed
exit
# (outside the container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

**Development mode (backend only):**

```shell
go run ./cmd/webapi/
```

**With a custom config file:**

```shell
go run ./cmd/webapi/ -c demo/demo.yml
# or
API_CONFIG_PATH=demo/demo.yml go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-node.sh
# (here you're inside the container)
yarn run dev
```

## Go vendoring

This project uses [Go Vendoring](https://go.dev/ref/mod#vendoring). After adding or updating a dependency, run:

```shell
go mod tidy && go mod vendor
```

then commit everything under `vendor/`.

## How to build for production

```shell
./open-node.sh
# (inside the container)
yarn run build-prod
```

## License

MIT — see [LICENSE](LICENSE).

Originally authored by [Enrico Bassetti](https://github.com/Enrico204). Forked and extended by [Tommaso Sgroi](https://github.com/tommaso-sgroi).
