# Backend Agent Brief

## Runtime & Entry Point

- Module `github.com/devbydaniel/release-notes-go`, targeting Go 1.23 with toolchain 1.23.4 (`go.mod`).
- `main.go` loads `.env` via `godotenv`, builds global config (`config.New()`), connects to Postgres through `internal/database`, configures Stripe (`internal/stripeUtil`) and MinIO-backed object storage (`internal/objstore`), then wires a `chi` router.
- Global logging is handled by `internal/logger` which bootstraps Zerolog with both console and Axiom writers; the logger must be cleaned up on shutdown.
- HTTP stack layers `chi` middlewares (logger, recoverer) plus custom middleware from `internal/middleware` before delegating to handlers.

## Directory Layout (high-level)

- `config/`: environment-backed runtime configuration structs (product info, email, payment, etc.).
- `internal/`: application code split by concern:
  - `database/`: Gorm setup, connection helpers, raw SQL migrations in `internal/database/migrations`.
  - `domain/<bounded-context>/`: each domain (users, release notes, subscriptions, widgets, etc.) follows a `model.go` + `repository.go` + `service.go` pattern, occasionally with `common.go`.
  - `handler/`: HTTP handlers grouped per route/page; each file owns a single area (login, release notes, widget, admin, etc.).
  - Supporting subsystems (`middleware`, `email`, `objstore`, `imgUtil`, `stripeUtil`, `logger`, `memcache`, `ratelimit`, `password`, `random`, `util`).
- `static/`: embedded CSS/JS/media plus widget assets (`static/static.go` uses `go:embed`).
- `templates/`: Go HTML templates (layouts/pages/partials) embedded via `templates/templates.go`.
- `Makefile`: helper tasks for migrations (`golang-migrate` CLI), Docker Compose dev stack, and Stripe webhook forwarding.
- `package.json`: only Prettier plus Go-template plugin for formatting templates/embedded assets.

## Request Lifecycle & Presentation

1. Router setup (`main.go`) defines public pages, authenticated dashboards, `/api` JSON endpoints, widget hosting, Stripe webhooks, and static asset serving.
2. `handler.NewHandler` bundles shared dependencies (DB, object store, Zerolog logger, Gorilla schema decoder). Each handler populates page-specific structs embedding `handler.BaseTemplateData`.
3. Templates are parsed/executed via `templates.Construct`/`templates.ExecuteTemplate`, mixing layouts and partials (navigation, header, HTMX snippets). Static files are exposed at `/static/*` and the widget script at `/widget`.
4. Public API routes enable CORS (`github.com/go-chi/cors`) with permissive defaults because the widget hard-codes `/api` and `/s` paths.

## Domain & Data Access Patterns

- Each domain service wraps a repository interface, driving business logic while keeping persistence concerns in repositories (`internal/domain/**/repository.go`). Services may open explicit transactions via `database.DB.StartTransaction`.
- Gorm (`gorm.io/gorm`) is the ORM of choice; repositories abstract filtering, pagination, and partial updates (using `Updates`/`Select` to whitelist fields).
- Object storage paths are managed inside repositories when media is involved. For example, `internal/domain/release-notes/service.go` calls `imgUtil` to transcode/rescale uploads before storing them via `objstore`.
- Session management hashes tokens (`sha256`) and persists them with rolling expiration; validation refreshes expiry and evicts expired sessions.
- RBAC is codified in `internal/domain/rbac` with `Role` + `Permission` enums and a helper `HasPermission`.
- Subscriptions interact with Stripe metadata; `subscription.Service` exposes helpers for CRUD and free/paid checks which feed middleware/handlers.

## Middleware & Security

- `mw.Handler` is instantiated with the DB and offers:
  - `Authenticate`: reads the session cookie, validates against the session domain, loads organisation/user context, and injects rich context keys (user/org IDs, roles, verification state, ToS/PP versions).
  - `Authorize` + `AuthorizeSuperAdmin`: enforce RBAC and super-admin-only routes via context data and `config.AdminUserId`.
  - `WithSubscriptionStatus`: augments the context with `HasActiveSubscription` for gating UI/actions.
  - `RateLimit`: simple token-bucket guard (per-user) backed by `internal/ratelimit`.
- Many handlers assume context keys exist; when adding new middleware ensure keys cascade before reaching handlers.
- CORS policies are explicitly defined for `/api`, `/widget`, `/s`, and `/stripe` routes; keep them in sync with frontend/widget expectations.

## Infrastructure & Integrations

- **Config**: `config.New()` reads environment variables (panic if missing) for base URL, product/legal copy, Postgres, MinIO, email, Stripe, Axiom, etc. Populate `.env` for local work—`main.initEnv()` loads it automatically.
- **Logging**: `internal/logger` sets `zerolog.TraceLevel` globally and multiplexes logs to stderr + Axiom. Always acquire loggers via `logger.Get()` to keep fields consistent.
- **Email**: `internal/email` switches between Postmark templates (production) and Mailcatcher SMTP (non-production). Templates expect specific `TemplateAlias` names (password-reset, welcome, user-invitation).
- **Object Storage**: `internal/objstore` provisions MinIO buckets (`release-notes`, `landing-page`), generates presigned URLs (proxied in non-prod), and exposes helpers for upload/delete.
- **Stripe**: `internal/stripeUtil` wraps checkout session creation, billing portal sessions, webhook verification, and subscription parsing. Metadata links Stripe subscriptions back to organisation IDs.
- **Caching & Rate Limiting**: `internal/memcache` wraps `patrickmn/go-cache` for ephemeral caches; `internal/ratelimit` implements an in-memory token bucket consumed by middleware—no cross-process coordination.
- **Binary assets**: `static/static.go` and `templates/templates.go` rely on `go:embed`. When adding files ensure glob patterns (`css/**/*`, `pages/*`, etc.) include the new assets.

## Operations & Local Dev

- Use `Makefile` targets to manage migrations (`migrate` CLI), Dockerized dependencies (`dev-start`, `dev-stop`, `dev-logs`), and Stripe webhook forwarding.
- Database migrations live in `internal/database/migrations` as timestamped SQL files; `migrations-unfuck` helps recover from partially applied migrations during development.
- `tmp/` holds build artifacts/logs; safe to clean locally but not tracked.

## Working Conventions

- Follow the handler/service/repository split when introducing new functionality: HTTP handlers stay thin, services coordinate validation + transactions, repositories own persistence.
- Prefer `logger.Get()` and structured log fields; trace logs are prevalent and expected.
- Use `config.New()` sparingly inside hot paths—consider reusing the module-level `cfg` created in `main.go` or relevant packages to avoid repeated env reads.
- Preserve context propagation: add values via middleware and access them in handlers/domain logic instead of re-querying the database.
- Respect existing CORS paths and widget-contract URLs (`/api`, `/widget`, `/s`) before renaming routes—front-end snippets reference them directly.
