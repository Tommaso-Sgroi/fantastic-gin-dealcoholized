# Fantastic Gin Dealcoholized

This repository contains the basic structure for [Web and Software Architecture](http://gamificationlab.uniroma1.it/en/wasa/) homework project.
It has been rewritten using [Gin](https://github.com/gin-gonic/gin) framework and other updated dependencies.

This project is a fork of [Fantastic Coffee (decaffeinated)](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated)
by [Enrico Bassetti](https://github.com/Enrico204) - originally created for the
[Web and Software Architecture](http://gamificationlab.uniroma1.it/en/wasa/) course at Sapienza University of Rome.

> ⚠️ **This is a _for-fun_ project and is not suitable for a real production environment.**

## What it does

The project itself shows an "actual" example implementation of the backend.  
It exposes a small REST API (documented in [`doc/api.yaml`](doc/api.yaml)) to create, read, update and delete "names" stored in a SQLite database:

| Method   | Path          | Description           |
|----------|---------------|-----------------------|
| `GET`    | `/healthz`    | Liveness probe        |
| `POST`   | `/name`       | Create a name         |
| `GET`    | `/name/:id`   | Get a name by ID      |
| `PUT`    | `/name/:id`   | Update a name by ID   |
| `DELETE` | `/name/:id`   | Delete a name by ID   |

All error responses are JSON: `{"error": "<message>"}`.

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is a small daemon for checking server health; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` is the main web API server daemon
* `demo/` contains the default YAML configuration file
* `db/` contains the SQL migration files
* `doc/` contains the OpenAPI 3.1 specification (`api.yaml`)
* `internal/` has all packages for implementing project-specific functionalities
	* `domain/` contains business entities and use-cases (for example the `name` and `health` domain packages)
	* `handler/` contains the HTTP transport layer (Gin router, middleware, request context, and route handlers)
	* `repository/` contains the persistence layer (SQLite queries and data access)
	* `errs/` contains shared application-level sentinel errors used to disambiguate failures across layers
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is a Vue.js web frontend; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

### Layered approach (brief)

The backend follows a simple layered approach:

* **HTTP layer (`internal/handler`)** receives requests, decodes/encodes JSON and maps domain errors to HTTP status codes.
* **Domain layer (`internal/domain`)** contains business rules and use-case orchestration, independent from underlying technology stack: Gin and SQLite.
* **Repository layer (`internal/repository`)** handles storage details and SQL queries.

Typical flow:

`HTTP request -> handler/routes -> domain use-case -> repository -> SQLite`

Dependencies are one-way, and shared error sentinels live in `internal/errs` so `errors.Is(...)` works consistently through wrapped errors.

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
