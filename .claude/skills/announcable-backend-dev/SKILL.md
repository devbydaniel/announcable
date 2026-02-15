---
name: announcable-backend-dev
description: Backend development in Announcable. Use when creating, modifying, or debugging backend code (Go, Chi router, GORM, templates, CSS/JS).
---

# Backend Development — Announcable

## Working Directory

**All commands run from `backend/`:**

```bash
cd backend
```

## Before You Start

Read the relevant documentation for what you're changing:

```bash
# Domain module
cat internal/domain/[module]/SUMMARY.md

# Handler architecture
cat internal/handler/SUMMARY.md

# Domain conventions
cat internal/domain/SUMMARY.md

# CSS architecture
cat assets/css/SUMMARY.md

# JS architecture
cat assets/js/SUMMARY.md
```

## Validation Sequence

Run after every change. Do NOT trust your own assessment — verify through observable behavior.

```bash
# Go code changes
go build -o ./tmp/main .          # Must compile
go vet ./...                      # No issues

# CSS/JS changes
npm run build                     # Vite must succeed

# Full check (if dev environment is running)
# Backend should respond at configured PORT
```

## Project Structure

```
backend/
├── main.go                    # Entry point, routes, server setup
├── config/                    # Environment configuration
├── internal/
│   ├── domain/                # Business logic (model, repository, service per module)
│   ├── handler/
│   │   ├── pages/             # HTML-serving handlers
│   │   ├── api/               # JSON API handlers
│   │   └── shared/            # Shared dependencies struct
│   ├── database/              # GORM setup, migrations, base model
│   ├── middleware/             # Auth, RBAC, rate limiting
│   ├── objstore/              # Minio wrapper
│   ├── email/                 # Email sending
│   └── logger/                # Structured logging
├── templates/                 # Go html/template (layouts, pages, partials)
├── assets/                    # Source CSS/JS (Vite input)
│   ├── css/                   # CSS source files
│   └── js/                    # JS source files
└── static/                    # Embedded assets (dist/ is Vite output)
```

## Domain Module Pattern

Every domain module follows the same trio structure:

```
internal/domain/{module}/
├── SUMMARY.md       # ← Read this first
├── model.go         # GORM model (embeds database.BaseModel)
├── repository.go    # Database access via GORM
├── service.go       # Business logic, orchestration
└── common.go        # Helper builders, shared types (optional)
```

**Key conventions:**
- Services are created with `NewService(repo)`
- Repositories require `NewRepository(db *database.DB)`
- `database.BaseModel` provides UUID `ID` and timestamps
- Use `logger.Get()` for structured logging
- Transactions via `repo.db.StartTransaction()` for multi-aggregate operations

**Package naming:** Singular, lowercase (`user`, `session`, `organisation` — NOT `users`, `sessions`)

## Handler Pattern

Handlers are organized by feature under `pages/` (HTML) and `api/` (JSON):

```
internal/handler/pages/{feature}/
├── page.go                         # Handlers struct, New(), GET handler, template
└── {domain}_{action}_action.go     # One POST/PATCH/DELETE per file
```

**Key conventions:**
- All handlers are methods on `*Handlers` struct receiving `*shared.Dependencies`
- One action per file, named `{domain}_{action}_action.go`
- Templates constructed with `templates.Construct()` helper
- Use `shared.BaseTemplateData` for template data

**Adding a new page:**
1. Create handler package under `pages/`
2. Define `Handlers` struct with `New(deps)` constructor
3. Add template in `templates/pages/`
4. Add CSS entry in `assets/css/pages/`
5. Register routes in `main.go`

## Template System

Templates use Go's `html/template` with composition:

```
root.html → layout (appframe/onboard/fullscreen) → page → partials
```

**Blocks defined by root.html:** `layout-css`, `page-css`, `body`, `layout-js`, `page-js`

```go
var myTmpl = templates.Construct(
    "my-page",
    "layouts/root.html",
    "layouts/appframe.html",
    "pages/my-page.html",
)

// In handler:
myTmpl.ExecuteTemplate(w, "root", data)
```

## CSS/JS Changes

Source files are in `assets/`, built by Vite to `static/dist/`:

```bash
# Development (watch mode)
npm run dev

# Production build
npm run build
```

**CSS:** Entry files in `assets/css/layouts/` and `assets/css/pages/` import from `assets/css/components/` and `assets/css/base/`. All `@import` must come first.

**JS:** Vanilla JavaScript (no modules). Alpine.js components, HTMX handlers, global functions. External deps available via `window.*` from vendor bundle.

**Never edit files in `static/dist/`** — they are auto-generated.

## Adding a New Feature (Checklist)

1. [ ] Read relevant SUMMARY.md files
2. [ ] Create/modify domain module (`model.go`, `repository.go`, `service.go`)
3. [ ] Create migration if schema changes needed (see `announcable-migrations` skill)
4. [ ] Create/modify handler (`page.go` + action files)
5. [ ] Create/modify templates
6. [ ] Add CSS/JS if needed
7. [ ] Register routes in `main.go`
8. [ ] `go build -o ./tmp/main .` compiles
9. [ ] `go vet ./...` passes
10. [ ] `npm run build` succeeds (if CSS/JS changed)
11. [ ] Test manually in browser

## External Dependencies

| Library | Purpose | Global |
|---------|---------|--------|
| Alpine.js | Reactive UI framework | `window.Alpine` |
| HTMX | AJAX and partial page updates | `window.htmx` |
| SweetAlert | Modal dialogs | `window.swal` |
| Toastify | Toast notifications | `window.Toastify` |
| Feather Icons | Icon library | `window.feather` |

All bundled in `/static/dist/vendor.js` and `/static/dist/vendor.css`.

## Anti-Patterns

| Don't | Why | Instead |
|-------|-----|---------|
| Edit `static/dist/` files | Auto-generated by Vite | Edit source in `assets/` |
| Call `config.New()` in hot paths | Creates new config each time | Cache at package level |
| Import across domain modules directly | Circular dependencies | Use IDs and defined structs |
| Use plural package names for domains | Convention conflict with handlers | Singular: `user`, `session` |
| Put `@import` after CSS rules | CSS spec violation | All `@import` at top of file |
| Skip reading SUMMARY.md | Miss module conventions | Always read before modifying |
| Batch multiple logical changes | Hard to identify breakage | One change → validate |
