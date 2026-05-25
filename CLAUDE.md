# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test

```bash
# Build
go build ./...

# Run the server (MySQL + Jaeger agent expected to be running)
go run cmd/api/main.go

# Run all tests
go test ./...

# Run a single test
go test -v -run TestTracing ./internal/middleware/

# Generate Swagger docs
swag init -g cmd/api/main.go
```

## Architecture

**Layered flow:** Handler → Service → DAO → Model (GORM) → MySQL

- `internal/routers/api/` — HTTP handlers (bind params, call service, write response via `app.Response`)
- `internal/service/` — business logic and transaction orchestration
- `internal/dao/` — data access wrappers around GORM model methods
- `internal/model/` — GORM models, DB init, and global callbacks

All routes are defined in `internal/routers/router.go`. Public routes (`/auth`, `/upload/file`, `/static`) go on the engine directly; `/api/v1/*` routes go through a group that applies the `JWT()` middleware.

**Startup order** (`cmd/api/main.go` init): Settings → DBEngine → Logger → Tracer. If any step fails, the process exits via `log.Fatalf`.

**Handler response pattern** — always use one of three methods:
```go
response := app.NewResponse(c)
response.ToResponse(data)              // single object or empty {}
response.ToResponseList(list, total)   // paginated list with pager
response.ToErrorResponse(errcode.Xxx)  // error with code + msg
```

## Middleware chain (applied to every request)

1. `gin.Logger()` + `gin.Recovery()` (debug) / `AccessLog()` + `Recovery()` (production)
2. `Tracing()` — Jaeger span creation, injects `X-Trace-ID` / `X-Span-ID` into Gin context
3. `RateLimiter()` — token-bucket rate limiting
4. `ContextTimeout()` — request timeout context
5. `Translations()` — validation error i18n (zh/en)
6. `JWT()` — on `/api/v1/*` routes only

## Key Patterns

**Global singletons** (`global/` package): App-wide state is stored as package-level vars initialized by `init()` in `cmd/api/main.go`. This includes `DBEngine`, all settings, `Logger`, and `Tracer`. There is no DI container — all packages read from `global` directly.

**`Service` and `Dao` are per-request**: `service.New(ctx)` and `dao.New(engine)` create fresh instances. `Service` holds both a `context.Context` and a `*dao.Dao`. DO NOT reuse them across goroutines.

**Model embedding — `*Model` pointer trapping**: All model types embed `*Model` (pointer, not value). When you create a zero-value struct (`BlogArticle{}`), `Model` is nil. Accessing `article.ID` (which is `article.Model.ID`) on that zero value causes a nil-pointer dereference. Always initialize:
```go
&BlogArticle{Model: &Model{}, Title: "Go"}
```

- Custom model methods that return single records (e.g. `Get`, `GetByAID`) MUST NOT swallow `gorm.ErrRecordNotFound` — returning a zero-value struct with nil Model will panic when callers access embedded fields. Return the error instead.
- `First()` returns `ErrRecordNotFound`; `Find()` returns an empty slice without error. The `List`/`ListByIDs` methods that use `Find()` don't have this problem.

**Soft-delete via GORM callbacks** (`internal/model/model.go`): Three registered callbacks hijack all CRUD operations. DELETE is rewritten to UPDATE setting `is_del=1` and `deleted_on`. All queries are scoped with `is_del=0`. Use `db.Unscoped()` to bypass.

**Log-to-trace linking**: `logger.WithTrace()` reads `X-Trace-ID` and `X-Span-ID` from the Gin context and appends them to every log line as JSON fields. Call `global.Logger.Infof(c, ...)` where `c` is `*gin.Context` — not a bare `context.Context`.

**Error codes** (`pkg/errcode/`): Centralized retcode system with module-specific ranges (tags: 2001xxxx, articles: 2002xxxx). Each error is a pre-instantiated singleton used throughout handlers and services.

**Configuration** (`configs/config.yaml`): Loaded via Viper at startup into package-level vars in `global/`. Five sections: Server, App, JWT, Database, Email. Structs defined in `pkg/setting/section.go`.

**Logger clone pattern**: `logger.Logger` methods chain via `clone()`, which deep-copies the `fields` map to avoid shared mutable state. Every log method calls `WithContext(ctx).WithTrace().Output(...)`.
