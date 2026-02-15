# Handler Architecture

## Overview

Handlers are organized by feature/page under two main directories:

- **`pages/`**: HTML-serving handlers for user-facing pages
- **`api/`**: JSON API endpoints

Each handler package shares dependencies via `shared.Dependencies` (DB, ObjStore, Log, Decoder).

## Structure Pattern

### Feature Package Layout

```
pages/
└── {feature}/
    ├── page.go                      # Handlers struct, New(), page data types, GET handler
    └── {domain}_{action}_action.go  # One file per POST/PATCH/DELETE action
```

### Example: Users Feature

```
pages/users/
├── page.go                    # ServeUsersPage, Handlers struct, New()
├── user_delete_action.go      # HandleUserDelete
├── invite_create_action.go    # HandleInviteCreate
└── invite_delete_action.go    # HandleInviteDelete
```

## File Conventions

### `page.go` - Shared Setup

Contains:

1. **Handlers struct** with dependencies
2. **New()** constructor
3. **Page data structures** for templates
4. **Page serving handler** (typically GET)
5. **Template constructor** (if applicable)

```go
package users

type Handlers struct {
    deps *shared.Dependencies
}

func New(deps *shared.Dependencies) *Handlers {
    return &Handlers{deps: deps}
}

type pageData struct {
    shared.BaseTemplateData
    Users []UserData
}

var pageTmpl = templates.Construct(...)

func (h *Handlers) ServeUsersPage(w http.ResponseWriter, r *http.Request) {
    // Render page
}
```

### `{domain}_{action}_action.go` - Single Action

One action per file. All handlers are methods on `Handlers` struct from `page.go`.

```go
package users

func (h *Handlers) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
    h.deps.Log.Trace().Msg("HandleUserDelete")
    // Action logic
}
```

## Naming Conventions

**Packages**: Use plural or descriptive names (`users`, `auth`, `releaseNotes`)  
**Files**: Use snake_case (`user_delete_action.go`, `password_update_action.go`)  
**Actions**: Common verbs: `create`, `update`, `delete`, `publish`, `cancel`

## Organization Rules

1. **One file per action** - No grouping decisions needed
2. **Use subfolders** only when multiple distinct page views exist
3. **Co-locate related actions** - Actions stay with their logical page even if serving different routes

## Usage in main.go

```go
import "github.com/.../handler/pages/users"

deps := shared.New(db, objStore)
usersHandler := users.New(deps)

r.Route("/users", func(r chi.Router) {
    r.Get("/", usersHandler.ServeUsersPage)
    r.Delete("/{id}", usersHandler.HandleUserDelete)
})

r.Route("/invites", func(r chi.Router) {
    r.Post("/", usersHandler.HandleInviteCreate)
    r.Delete("/{id}", usersHandler.HandleInviteDelete)
})
```
