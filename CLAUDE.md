# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo for Announcable, a release notes and announcement platform. It consists of two main components:

- **Backend**: Go web application using Chi router, GORM, PostgreSQL, Stripe, and Minio for object storage
- **Widget**: React/TypeScript embeddable widget built with Vite that displays release notes

## Development Commands

### Starting Development Environment

```bash
cd backend

# Start all services (backend with hot-reload, postgres, pgadmin, mail, minio)
make dev-start

# Stop all services
make dev-stop

# Follow logs from all services
make dev-logs
```

The backend runs with Air for hot-reloading. Changes to Go files, templates, or static assets trigger automatic rebuilds.

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

# Stripe webhook testing (requires Stripe CLI)
make stripe-webhook
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
- **`internal/stripeUtil/`**: Stripe payment integration utilities
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

**CSS Organization** (`backend/static/css/`):

CSS follows a component-based architecture mirroring template structure:

- **`base/`**: CSS reset and CSS variables (design tokens)
- **`components/`**: Reusable UI components (button.css, card.css, form.css, nav.css, etc.)
- **`layouts/`**: Layout-specific styles (appframe.css, onboard.css, fullscreenmessage.css)
- **`pages/`**: Page-specific styles (login.css, release-notes-list.css, etc.)

CSS is loaded progressively through template blocks:

1. Base CSS always loads in `root.html` (`reset.css`, `variables.css`)
2. Layout CSS loads via `layout-css` block (e.g., `appframe.css`, `nav.css`)
3. Page CSS loads via `page-css` block (e.g., `login.css` plus needed component CSS)

**JavaScript Organization** (`backend/static/js/`):

- **`app/`**: App-level utilities (confirmDialog.js, successMsg.js)
- **`components/`**: Reusable component scripts (toast.js)
- **`pages/`**: Page-specific scripts

JS loads similarly to CSS using `layout-js` and `page-js` blocks.

**External Dependencies**:

The application uses CDN-hosted libraries loaded in `root.html`:

- **HTMX**: Enables AJAX requests and partial page updates without JavaScript
- **Alpine.js**: Lightweight reactive framework for UI interactions
- **Feather Icons**: Icon library
- **SweetAlert**: Modal dialogs
- **Toastify**: Toast notifications

**Static Asset Serving**:

All static assets are embedded in the binary via `//go:embed` in `backend/static/static.go` and served at the `/static/*` route.

### Widget Structure

The widget is a self-contained React application built as a UMD bundle:

- **`src/main.tsx`**: Widget entry point that mounts to a container element
- **`src/App.tsx`**: Main widget component
- **`src/components/`**: React components
- **`src/hooks/`**: Custom React hooks
- **`src/lib/`**: Utilities and helpers
- **`vite.config.ts`**: Builds as UMD library with CSS namespaced to `.announcable-widget`

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

- **Stripe**: Payment processing and subscription management
- **Axiom**: Logging and observability
- **Postmark/Mailcatcher**: Email delivery (Postmark for production, Mailcatcher for development)
- **Minio**: S3-compatible object storage

## Environment Setup

Copy `.env.example` to `.env` and configure:

- Database credentials (Postgres)
- Stripe API keys and webhook secret
- Minio/object storage credentials
- Email service credentials (Postmark for production, Mailcatcher for dev)
- Axiom logging token (optional)
- Legal document versions (TOS_VERSION, PP_VERSION)

## Docker Services

Development uses Docker Compose with two config files:

- **`compose-base.yml`**: Base services (postgres, pgadmin)
- **`compose-dev.yml`**: Development overrides (app with Air, mailcatcher, minio, port mappings)

Services:

- `app`: Go backend with Air hot-reload
- `postgres`: PostgreSQL database
- `pgadmin`: Database admin interface
- `mail`: Mailcatcher for email testing (http://localhost:1080)
- `objstorage`: Minio (ports 9000/9001)

## Code Generation

The backend uses `go generate` to build templ templates. This runs automatically during Air rebuilds but must be run manually if building without Air.

## Important Notes

- Handler files follow strict naming: `handle-<action>.go`
- Templates are Go html/template, not templ (despite templ being in dependencies)
- Widget CSS is namespaced to `.announcable-widget` to avoid conflicts when embedded
- Migrations must be run manually - they don't auto-run on startup
- Use `make migrations-unfuck` to recover from migration version mismatches
