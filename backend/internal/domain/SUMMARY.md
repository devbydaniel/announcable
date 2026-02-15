# Domain Agent Brief

## Conventions

- Every bounded context exposes a `service` struct backed by a `repository`. Services contain orchestration/validation while repositories wrap `gorm` access via `database.DB`. Most packages also define `model.go` for persistence structs and `common.go` for helper builders.
- `database.BaseModel` embeds shared columns (UUID `ID`, timestamps). Cross-package associations use typed UUIDs plus imported structs (e.g. release notes load `organisation.Organisation`).
- Services are instantiated with `NewService(repo)`; repositories require a shared DB connection (`NewRepository(db)`).
- Long-running operations often start transactions by calling `repo.db.StartTransaction()` and passing `tx.Tx` down to repository methods.
- All packages pull a structured logger via `logger.Get()` and log at `Trace`/`Debug` granularity.

### Package Naming Conventions

Domain packages follow idiomatic Go naming principles:

**✅ Use singular, lowercase names**:

- `user` (not `users`, `userDomain`, or `user_domain`)
- `organisation` (not `organisations`)
- `session` (not `sessions`)

**Why singular?**

1. **Standard Go practice**: Represents a single domain entity
2. **Natural collision avoidance**: Handler packages use plural (`users` handler vs `user` domain)
3. **Clean imports**: `user.NewService()` reads naturally
4. **No stuttering**: Avoid `user.User` when `user.Model` or domain types suffice

**Example of natural collision avoidance**:

```go
// Handler package (plural)
package users

import (
    "github.com/.../domain/user"          // Domain: singular
    "github.com/.../domain/organisation"  // Domain: singular
    "github.com/.../domain/session"       // Domain: singular
)

func (h *Handlers) HandleUserDelete(...) {
    userService := user.NewService(...)       // Clear: domain service
    orgService := organisation.NewService(...)
    sessionService := session.NewService(...)
}
```

**When collisions DO occur** (rare cases like `admin` package used in both domains and handlers):

```go
// Use import aliases - the idiomatic Go solution
import (
    adminDomain "github.com/.../domain/admin"
    adminHandler "github.com/.../handler/pages/admin"
)

// Or shorter alias for one:
import (
    "github.com/.../domain/admin"
    adminpages "github.com/.../handler/pages/admin"
)
```

**Avoid**:

- ❌ Suffixes: `userDomain`, `userService`, `userRepo` as package names
- ❌ Plural domains: `users`, `organisations` (reserve plural for collections/handlers)
- ❌ Underscores: `user_domain`, `user_service` (use camelCase for multi-word like `releasePage`)
- ❌ Abbreviations: `usr`, `org`, `sess` (spell out full words)

## Usage Notes

- When adding a new domain, mirror the established trio (`model`, `repository`, `service`) and import shared utilities instead of reimplementing (e.g., use `password` for hashing, `random` for tokens, `imgUtil` for uploads).
- Avoid calling `config.New()` repeatedly inside hot code paths—prefer caching values at package level like existing services do.
- Prefer repository helpers that accept transactions if your feature spans multiple aggregates; transactions originate from `database.DB.StartTransaction()`.
- Keep cross-domain references flowing through IDs and clearly defined structs to prevent circular dependencies; existing packages re-export only the types needed by others.
