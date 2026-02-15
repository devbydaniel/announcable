# Architecture

**2-word:** Release Notes

**8-word:** Release notes platform with embeddable widget and admin.

**32-word:** Announcable is a release notes and announcement platform. Users create, manage, and publish release notes through a Go web application. An embeddable Lit/TypeScript widget displays release notes on customer websites. Multi-tenant with RBAC.

---

## Repository Structure

```
announcable/
├── backend/                  # Go web application (Chi, GORM, PostgreSQL, Minio)
├── widget/                   # Lit/TypeScript embeddable widget (Vite, UMD)
├── docker-compose.yml        # Production stack
├── docker-compose.dev.yml    # Development services
├── ARCHITECTURE.md           # This file
├── AGENTS.md                 # AI coding agent guidelines
└── .env.example              # Environment configuration template
```

---

## Backend (Go + Chi + GORM)

📁 **[`backend/`](backend/)**

### Domain Modules — Business Logic

| Module | Purpose | Path |
|--------|---------|------|
| [user](backend/internal/domain/user/SUMMARY.md) | User accounts & credentials | `internal/domain/user/` |
| [organisation](backend/internal/domain/organisation/SUMMARY.md) | Multi-tenant organization management | `internal/domain/organisation/` |
| [session](backend/internal/domain/session/SUMMARY.md) | Session-based authentication | `internal/domain/session/` |
| [rbac](backend/internal/domain/rbac/SUMMARY.md) | Role-based access control | `internal/domain/rbac/` |
| [release-notes](backend/internal/domain/release-notes/SUMMARY.md) | Release note CRUD & publishing | `internal/domain/release-notes/` |
| [release-note-likes](backend/internal/domain/release-note-likes/SUMMARY.md) | User reactions to release notes | `internal/domain/release-note-likes/` |
| [release-note-metrics](backend/internal/domain/release-note-metrics/SUMMARY.md) | View/engagement tracking | `internal/domain/release-note-metrics/` |
| [widget-configs](backend/internal/domain/widget-configs/SUMMARY.md) | Embeddable widget configuration | `internal/domain/widget-configs/` |
| [release-page-configs](backend/internal/domain/release-page-configs/SUMMARY.md) | Public release page configuration | `internal/domain/release-page-configs/` |
| [admin](backend/internal/domain/admin/SUMMARY.md) | Super admin platform management | `internal/domain/admin/` |

### Handler Layer — HTTP Interface

| Area | Purpose | Path |
|------|---------|------|
| [pages/auth](backend/internal/handler/pages/auth/) | Login, register, password flows, email verification | `internal/handler/pages/auth/` |
| [pages/release_notes](backend/internal/handler/pages/release_notes/) | Release note list, create, detail pages | `internal/handler/pages/release_notes/` |
| [pages/settings](backend/internal/handler/pages/settings/) | Account settings | `internal/handler/pages/settings/` |
| [pages/users](backend/internal/handler/pages/users/) | User management & invites | `internal/handler/pages/users/` |
| [pages/widget](backend/internal/handler/pages/widget/) | Widget configuration page | `internal/handler/pages/widget/` |
| [pages/release_page](backend/internal/handler/pages/release_page/) | Release page configuration | `internal/handler/pages/release_page/` |
| [pages/admin](backend/internal/handler/pages/admin/) | Admin dashboard & org management | `internal/handler/pages/admin/` |
| [pages/public](backend/internal/handler/pages/public/) | Home, public release page, widget script | `internal/handler/pages/public/` |
| [api/widget](backend/internal/handler/api/widget/) | Widget JSON API (release notes, metrics, likes) | `internal/handler/api/widget/` |
| [api/shared](backend/internal/handler/api/shared/) | Shared API handlers (object storage proxy, 404) | `internal/handler/api/shared/` |

### Infrastructure

| Module | Purpose | Path |
|--------|---------|------|
| database | GORM setup, migrations, base model | `internal/database/` |
| middleware | Auth, RBAC, rate limiting middleware | `internal/middleware/` |
| objstore | Minio object storage wrapper | `internal/objstore/` |
| email | Email sending (Postmark/Mailcatcher) | `internal/email/` |
| logger | Structured logging (Zerolog + Axiom) | `internal/logger/` |
| config | Environment configuration | `config/` |
| templates | Go html/template system | `templates/` |
| static | Embedded static assets (CSS, JS, media) | `static/` |
| assets | Source CSS/JS (Vite-bundled) | `assets/` |

---

## Widget (Lit + TypeScript + Vite)

📁 **[`widget/`](widget/)**

| Layer | Purpose | Path |
|-------|---------|------|
| main.ts | Entry point, registers custom element | `src/main.ts` |
| app.ts | Main widget web component (LitElement) | `src/app.ts` |
| [components](widget/src/components/) | Lit UI components (icons, widget views) | `src/components/` |
| [controllers](widget/src/controllers/) | Lit reactive controllers | `src/controllers/` |
| [tasks](widget/src/tasks/) | Async data fetching tasks | `src/tasks/` |
| [lib](widget/src/lib/) | Utilities, config, types, contexts | `src/lib/` |

Built as UMD bundle (`widget.js`) embedded on customer websites via script tag.

---

## Key Architectural Patterns

### Handler-Based Architecture (Backend)

```
┌──────────────────────────────────────────────────────────────┐
│                    Handler Layer                              │
│         pages/ (HTML) and api/ (JSON)                        │
│         Each handler receives shared.Dependencies            │
├──────────────────────────────────────────────────────────────┤
│                    Domain Layer                               │
│         service → repository → database.DB                   │
│         Pure business logic, no HTTP concepts                │
├──────────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                         │
│         database, objstore, email, middleware                │
└──────────────────────────────────────────────────────────────┘
```

### Template Composition (Backend)

```
root.html (base skeleton)
  └── layout (appframe/onboard/fullscreen)
        └── page content
              └── partials (nav, header, HTMX fragments)
```

### CSS/JS Build Pipeline

```
assets/css/ ──→ Vite ──→ static/dist/ ──→ //go:embed ──→ binary
assets/js/  ──→ Vite ──→ static/dist/ ──→ //go:embed ──→ binary
```

---

## Quick Navigation

- **Adding a release note feature**: Start at [release-notes SUMMARY](backend/internal/domain/release-notes/SUMMARY.md)
- **Adding a page**: See [handler SUMMARY](backend/internal/handler/SUMMARY.md)
- **Understanding auth**: [session](backend/internal/domain/session/SUMMARY.md) + [rbac](backend/internal/domain/rbac/SUMMARY.md)
- **Widget API flow**: [api/widget](backend/internal/handler/api/widget/) serves data → widget fetches via [tasks](widget/src/tasks/)
- **CSS/JS changes**: See [CSS SUMMARY](backend/assets/css/SUMMARY.md) and [JS SUMMARY](backend/assets/js/SUMMARY.md)
