# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo for Announcable, a release notes and announcement platform. It consists of two main components:

- **Backend**: Go web application using Chi router, GORM, PostgreSQL, and Minio for object storage
- **Widget**: Lit/TypeScript embeddable widget built with Vite that displays release notes

## Development Commands

### Starting Development Environment

```bash
cd backend

# Start dev services (postgres, mail, minio) - runs in background
make dev-services

# In a new terminal - start Go backend with Air hot-reload
make dev-air

# In a new terminal - start Vite for CSS/JS hot-reload
npm run dev

# Stop all services
make dev-stop

# Follow logs from all services
make dev-logs
```

The backend runs with Air for hot-reloading. Changes to Go files, templates, or static assets trigger automatic rebuilds. The app runs directly on your host (not in Docker) for faster rebuilds and easier debugging.

### Production Deployment

```bash
# Start complete production stack (external nginx handles reverse proxy)
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

### Backend Development

```bash
cd backend

# Build (includes templ code generation)
go generate ./... && go build -o ./tmp/main .

# Database migrations
make migrations-new name=<migration_name>    # Create new migration
make migrations-up                           # Run one migration
make migrations-up-all                       # Run all migrations
make migrations-down                         # Rollback one migration
make migration-force version=<version>       # Force to specific version
make migrations-unfuck                       # Fix migration state issues

```

### Widget Development

```bash
cd widget

# Development mode with hot-reload
npm run dev

# Build for production
npm run build

# Build for development/testing
npm run build:test

# Lint
npm run lint
```

## Architecture

### Backend Structure

The backend follows a handler-based architecture with clean separation of concerns:

- **`main.go`**: Application entry point, initializes database, object storage, middleware, and routes
- **`internal/handler/`**: HTTP handlers (one file per route/handler)
- **`internal/domain/`**: Domain models organized by entity (user, organisation, release-notes, session, rbac, etc.)
- **`internal/database/`**: Database layer with GORM setup and model definitions
- **`internal/middleware/`**: Custom middleware for auth, RBAC, etc.
- **`internal/objstore/`**: Minio object storage wrapper
- **`templates/`**: Go HTML templates organized into layouts, pages, and partials
- **`static/`**: Static assets (CSS, JS, media, widget)
- **`config/`**: Configuration management

**Key Patterns**:

- Handlers are named `handle-<action>.go` (e.g., `handle-login.go`, `handle-release-note-create.go`)
- Each handler receives a `Handler` struct with dependencies (DB, ObjStore, logger, decoder)
- Templates use Go's html/template with a base template data structure
- RBAC system controls access to organization resources

### Template System

The backend uses Go's `html/template` with a structured, composable template system:

**Template Organization** (`backend/templates/`):

- **`layouts/`**: Base page structures that define the overall page skeleton
  - `root.html`: The foundational template with `<head>`, external dependencies (HTMX, Alpine.js, Feather icons), and block definitions for CSS/JS injection
  - `appframe.html`: Authenticated app layout with navigation and header
  - `onboard.html`: Minimal layout for auth flows (login, register)
  - `fullscreenmessage.html`: Centered message layout
- **`pages/`**: Page-specific content that fills layout blocks (login.html, release-notes-list.html, etc.)
- **`partials/`**: Reusable components (nav.html, header.html, hx-prefixed HTMX partials)

**Template Composition Pattern**:

Templates are composed using Go's `define` and `block` directives:

1. `root.html` defines the base structure with blocks: `layout-css`, `page-css`, `body`, `layout-js`, `page-js`
2. Layout templates (e.g., `appframe.html`) implement `layout-css`, `layout-js`, and `body`, and define a `main` block
3. Page templates implement `page-css`, `page-js`, and `main` blocks with actual content

**Handler Template Usage**:

```go
// Construct template once (typically as package-level var)
var loginTmpl = templates.Construct(
    "login",
    "layouts/root.html",
    "layouts/onboard.html",
    "pages/login.html",
)

// Execute in handler
loginTmpl.ExecuteTemplate(w, "root", data)
```

The `templates.Construct()` helper automatically includes all partials and embeds templates via `//go:embed`.

**CSS Organization**:

CSS uses Vite bundling with a component-based architecture:

**Source CSS** (`backend/static/css/`):
- **`base/`**: CSS reset and CSS variables (design tokens)
- **`components/`**: Reusable UI components (button.css, card.css, form.css, nav.css, etc.)
- **`pages/`**: Page-specific styles (login.css, release-notes-list.css, etc.)

**Bundled CSS Entry Files** (`backend/assets/css/`):
- **`layouts/`**: Layout entry files with inline styles + component imports (appframe.css, onboard.css, fullscreen.css)
- **`pages/`**: Page entry files that import components and page-specific CSS (to be created)

**Build Output** (`backend/static/dist/`):
- Vite bundles and minifies CSS from `assets/css/` into `static/dist/`
- Each entry file becomes one bundled, minified CSS file

**CSS Loading Pattern**:
1. Layout CSS bundle loads via `layout-css` block (e.g., `/static/dist/layouts/appframe.css`)
2. Page CSS bundle loads via `page-css` block (e.g., `/static/dist/pages/login.css`)
3. Each bundle includes all required base, component, and page CSS in one optimized file

**Important**: All `@import` statements must come at the very beginning of CSS files before any other CSS rules.

**JavaScript Organization**:

JavaScript uses Vite bundling with the same architecture as CSS:

**Source JS** (`backend/assets/js/`):
- **`app/`**: App-level utilities (confirmDialog.js, successMsg.js)
- **`components/`**: Reusable component scripts (toast.js, file-input.js, popover.js)
- **`pages/`**: Page-specific scripts

**Note**: JS files are built individually (not bundled together)

**Build Output** (`backend/static/dist/`):
- Vite minifies each JS file from `assets/js/` into `static/dist/`
- Directory structure is preserved (e.g., `assets/js/app/confirmDialog.js` → `static/dist/app/confirmDialog.js`)

**JS Loading Pattern**:
1. App-level JS files load in `root.html` header (e.g., `/static/dist/app/confirmDialog.js`, `/static/dist/components/toast.js`)
2. Page-specific JS files load via `page-js` block (e.g., `/static/dist/pages/settings.js`)
3. Component JS files are loaded explicitly where needed

**Development Workflow**:
- Source files in `assets/js/` for editing
- `npm run dev` for watch mode with hot-reload
- `npm run build` for production minification
- Minified output in `static/dist/` (embedded in Go binary)

For detailed JS architecture guidance, see `backend/assets/js/AGENTS.md`.

**External Dependencies**:

The application uses npm packages bundled into a vendor bundle via Vite:

- **Alpine.js**: Reactive UI framework for interactive components
- **HTMX**: AJAX requests and partial page updates
- **Toastify**: Toast notifications
- **SweetAlert**: Modal dialogs
- **Feather Icons**: Icon library

All dependencies are bundled into `/static/dist/vendor.js` and `/static/dist/vendor.css` and loaded in `root.html`. They are exposed as global objects (window.Alpine, window.htmx, etc.) for use in vanilla JavaScript files.

**Static Asset Serving**:

All static assets are embedded in the binary via `//go:embed` in `backend/static/static.go` and served at the `/static/*` route.

### Widget Structure

The widget is a self-contained Lit web component application built as a UMD bundle:

- **`src/main.ts`**: Widget entry point that registers the custom element
- **`src/app.ts`**: Main widget web component (LitElement)
- **`src/components/`**: Lit components
- **`src/controllers/`**: Lit reactive controllers
- **`src/tasks/`**: Async task handlers for data fetching
- **`src/lib/`**: Utilities and helpers
- **`vite.config.ts`**: Builds as UMD library

The widget is embedded in customer websites via a script tag and displays release notes fetched from the backend API.

### Database

- PostgreSQL for primary data storage
- GORM for ORM with migrations in `internal/database/migrations/`
- Minio for object storage (images, attachments)

### Authentication & Authorization

- Session-based authentication
- RBAC (Role-Based Access Control) for organization-level permissions
- Middleware enforces authentication and authorization rules
- Email verification required for new accounts

### External Services

- **Axiom**: Logging and observability
- **Postmark/Mailcatcher**: Email delivery (Postmark for production, Mailcatcher for development)
- **Minio**: S3-compatible object storage

## Environment Setup

Copy `.env.example` to `.env` and configure:

- Database credentials (Postgres)
- Minio/object storage credentials
- Email service credentials (Postmark for production, Mailcatcher for dev)
- Axiom logging token (optional)

## Docker Services

### Development (docker-compose.dev.yml)
Services run in Docker, app runs on host:
- `postgres`: PostgreSQL database (port 5432)
- `pgadmin`: Database admin interface (port from PGADMIN_PORT)
- `mail`: Mailcatcher for email testing (SMTP: 1025, Web: 1080)
- `minio`: Object storage (API: 9000, Console: 9001)

### Production (docker-compose.yml)
Complete stack (external nginx handles SSL/reverse proxy):
- `app`: Go backend (exposed on PORT)
- `postgres`: PostgreSQL database (internal only)
- `minio`: Object storage (internal only)
- `pgadmin`: Database admin (exposed on PGADMIN_PORT)

## Code Generation

The backend uses `go generate` to build templ templates. This runs automatically during Air rebuilds but must be run manually if building without Air.

## Important Notes

- Handler files follow strict naming: `handle-<action>.go`
- Templates are Go html/template, not templ (despite templ being in dependencies)
- Widget CSS is namespaced to `.announcable-widget` to avoid conflicts when embedded
- Migrations must be run manually - they don't auto-run on startup
- Use `make migrations-unfuck` to recover from migration version mismatches
