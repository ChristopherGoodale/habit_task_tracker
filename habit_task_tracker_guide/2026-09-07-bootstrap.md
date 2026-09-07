# 2026-09-07 — Project Bootstrap

## The "Why" (Design Decisions)

**Stdlib `net/http` instead of Chi/Gin/etc.** Go 1.22 added method- and
path-parameter-aware routing directly to `http.ServeMux`
(`mux.HandleFunc("PUT /api/tasks/{id}", ...)`, `r.PathValue("id")`), which
used to be the main reason people reached for a router library. Since this
project's whole point (per the project ranking doc) is seeing every piece of
the request lifecycle with nothing hidden, adding a router dependency here
would work against the goal rather than for it. A router library becomes
worth it once you need things stdlib genuinely lacks — grouped
sub-routers, built-in middleware chaining helpers — which is a fair
trade-off to revisit in project #3 (the real-time project) if it comes up.

**Plain `schema.sql` instead of a migration tool.** With one table and no
collaborators, a versioned migration tool (goose, golang-migrate, atlas)
buys nothing yet — it's solving a team/history problem this project doesn't
have. `db.Connect` embeds `schema.sql` via `go:embed` and runs it with
`CREATE TABLE IF NOT EXISTS` on every startup, which is idempotent and keeps
"what's the schema" answerable by reading one file. This is a deliberate
choice to revisit once the schema needs to *change* over time rather than
just exist.

**`database/sql` + `pgx/v5/stdlib`, no ORM.** `pgx` is used only as a
driver (registered under the `database/sql` interface), not as its own
query API. Every query in `internal/task/store.go` is hand-written SQL —
the point of doing this in Go before doing the ASP.NET Core + EF Core
project (#2) is to have a concrete "what does an ORM actually save me"
comparison once EF Core's migrations and change-tracking show up.

**Postgres via Docker Compose, not SQLite.** The ranking doc lists both as
options; Postgres was chosen deliberately over the simpler SQLite path to
practice a connection-pooled, networked database from the start, since
that's the shape every later project in the sequence uses. Docker Compose
keeps "start the database" to one command without installing Postgres
system-wide.

**Hand-rolled logging/CORS middleware instead of a middleware library.**
Both are `func(http.Handler) http.Handler` — the same pattern any Go
middleware library builds on top of — so this project shows that pattern
directly rather than through a library's abstraction over it.

**Postgres on host port 5433, not 5432.** During verification, the Go
server's connection to `localhost:5432` was silently hijacked by a
pre-existing native Postgres Windows service running on this machine
(different credentials, so it failed with a SASL auth error that looked
like a Docker Compose misconfiguration but wasn't). Docker's own port
mapping is not exclusive — if something else already owns the host port,
connections can land on the wrong server. `docker compose ps` showing
`healthy` only confirms the *container* is fine, not that nothing else on
the host is fighting it for the port.

## Core Concepts

- **`http.ServeMux` pattern routing** (Go 1.22+): patterns can include an
  HTTP method prefix (`"GET /api/tasks"`) and `{param}` path segments,
  retrieved in a handler via `r.PathValue("param")`.
- **`go:embed`**: compiles `schema.sql`'s contents into the Go binary as a
  string constant, so the binary has no runtime dependency on a file being
  present alongside it.
- **Middleware as handler-wrapping functions**: `func(http.Handler) http.Handler`
  takes a handler and returns a new one that runs code before/after calling
  the original — composed in `main.go` by wrapping `mux` twice
  (`middleware.CORS(...)` then `middleware.Logging(...)`).
- **Controlled components in React**: `TaskForm`'s inputs are driven by
  `useState`, with `value` and `onChange` wired together — the standard
  React form pattern used throughout this sequence.

## Implementation Breakdown

- [`server/internal/db/db.go`](../server/internal/db/db.go) — opens the
  connection, pings it (fails fast if Postgres isn't reachable), then execs
  the embedded schema.
- [`server/internal/task/store.go`](../server/internal/task/store.go) — one
  method per operation (`List`, `Create`, `Update`, `Delete`), each a single
  parameterized query. `ErrNotFound` is a sentinel error the handlers check
  with `errors.Is` to decide between a 404 and a 500.
- [`server/internal/task/handlers.go`](../server/internal/task/handlers.go) —
  translates HTTP (JSON body, path params, status codes) to/from calls on
  the store. No business logic lives here beyond input validation
  (non-empty title).
- [`server/main.go`](../server/main.go) — the only place all the pieces
  (config from env vars, db connection, store, handlers, middleware, mux)
  get wired together, so the whole app's dependency graph is readable in one
  short file.
- [`client/src/api/tasks.js`](../client/src/api/tasks.js) — a thin `fetch`
  wrapper; no axios, no React Query yet (those show up in later projects in
  the sequence once caching/refetching actually matter).
- [`client/src/App.jsx`](../client/src/App.jsx) — owns the `tasks` array in
  state and passes down callbacks (`onCreate`/`onToggle`/`onDelete`) to
  `TaskForm`/`TaskList`/`TaskItem`, which stay presentational.
