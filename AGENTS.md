# AI Coding Agent Guidelines

> **Philosophy**: Verify through observable behavior. Lint, build, and run — don't trust assumptions.

For architecture overview and module navigation, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Repository Overview

**Announcable** is a release notes and announcement platform consisting of a Go web backend and an embeddable Lit/TypeScript widget.

```
announcable/
├── backend/          # Go web app (Chi, GORM, PostgreSQL, Minio)
├── widget/           # Lit/TypeScript embeddable widget (Vite, UMD)
├── docker-compose.yml        # Production stack
├── docker-compose.dev.yml    # Development services
├── ARCHITECTURE.md           # Full architecture docs with module index
└── AGENTS.md                 # This file
```

---

## Module Boundaries

The backend enforces clean separation between domain modules:

- **`internal/domain/*`** — Business logic (services, repositories, models)
- **`internal/handler/*`** — HTTP handlers (pages for HTML, api for JSON)
- **`internal/middleware/`** — Auth, RBAC, rate limiting
- **`internal/database/`** — GORM setup, migrations, base model
- **`config/`** — Environment configuration

Before modifying any module, read its documentation:

```bash
# Domain modules
cat backend/internal/domain/[module]/SUMMARY.md

# Handler architecture
cat backend/internal/handler/SUMMARY.md

# CSS/JS architecture
cat backend/assets/css/SUMMARY.md
cat backend/assets/js/SUMMARY.md
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the complete module index with links.

---

## Key Files

| Purpose | Location |
|---------|----------|
| Architecture & module index | [ARCHITECTURE.md](ARCHITECTURE.md) |
| Domain module summaries | `backend/internal/domain/[module]/SUMMARY.md` |
| Handler architecture | `backend/internal/handler/SUMMARY.md` |
| Domain conventions | `backend/internal/domain/SUMMARY.md` |
| CSS architecture | `backend/assets/css/SUMMARY.md` |
| JS architecture | `backend/assets/js/SUMMARY.md` |
| Database migrations | `backend/internal/database/migrations/` |
| Template system | `backend/templates/` |
| Widget source | `widget/src/` |

---

## Core Principles

### 1. Validation-First

Do NOT trust your own assessment of code correctness. Verify through observable behavior — build, lint, and run. See the development skills below for specific validation sequences.

### 2. Incremental Progress

- Make one change at a time
- Validate after each change
- Never batch multiple logical changes

### 3. Respect Boundaries

- Read the target module's SUMMARY.md before making changes
- Domain packages are singular (`user`, `session`) — handler packages are plural/descriptive
- Cross-domain references flow through IDs and clearly defined structs
- Never edit generated/built files in `static/dist/`

### 4. Follow Existing Patterns

- Handlers follow the structure in `internal/handler/SUMMARY.md`
- Domain modules follow the trio pattern: `model.go`, `repository.go`, `service.go`
- CSS/JS follow the Vite pipeline documented in their respective SUMMARY.md files

---

## Development Skills

For detailed development workflows, patterns, and validation checklists, use the appropriate skill:

- **Backend work** → `announcable-backend-dev` skill (Go, Chi, GORM, templates, CSS/JS)
- **Widget work** → `announcable-widget-dev` skill (Lit, TypeScript, Vite, UMD bundle)
- **Database migrations** → `announcable-migrations` skill (golang-migrate, raw SQL)

---

## Development Commands (Quick Reference)

### Backend

```bash
cd backend
make dev-services        # Start postgres, mail, minio
make dev-air             # Go hot-reload
npm run dev              # Vite CSS/JS watch
npm run build            # Production CSS/JS build
go build -o ./tmp/main . # Build Go binary
```

### Widget

```bash
cd widget
npm run dev              # Dev server with hot-reload
npm run build            # Production UMD bundle
npm run lint             # Lint check
```

### Migrations

```bash
cd backend
make migrations-new name=<name>   # Create new migration
make migrations-up                # Run one migration
make migrations-up-all            # Run all migrations
make migrations-down              # Rollback one
```
