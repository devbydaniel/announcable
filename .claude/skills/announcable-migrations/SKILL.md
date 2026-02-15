---
name: announcable-migrations
description: Creating database migrations in Announcable. Use when schema changes are needed (new columns, tables, indexes, constraints).
---

# Database Migrations — Announcable Backend

## Working Directory

```bash
cd backend
```

## Migration Tool

Announcable uses **golang-migrate** with raw SQL migration files. Migrations are NOT auto-generated — you write them by hand.

Migration files live in `internal/database/migrations/` with sequential numbering:

```
internal/database/migrations/
├── 000001_initial.up.sql
├── 000001_initial.down.sql
├── 000002_add_widget_configs.up.sql
├── 000002_add_widget_configs.down.sql
└── ...
```

## Workflow

### 1. Create the Migration Files

```bash
make migrations-new name=descriptive_name
```

This creates two files:
- `{seq}_descriptive_name.up.sql` — applies the change
- `{seq}_descriptive_name.down.sql` — reverts the change

### 2. Write the SQL

**Up migration** — apply the schema change:

```sql
-- 000005_add_tags_to_release_notes.up.sql
ALTER TABLE release_notes ADD COLUMN tags TEXT[];
CREATE INDEX idx_release_notes_tags ON release_notes USING GIN(tags);
```

**Down migration** — revert symmetrically:

```sql
-- 000005_add_tags_to_release_notes.down.sql
DROP INDEX IF EXISTS idx_release_notes_tags;
ALTER TABLE release_notes DROP COLUMN IF EXISTS tags;
```

### 3. Run the Migration

```bash
make migrations-up        # Run one migration
make migrations-up-all    # Run all pending
```

### 4. Update the GORM Model

Update the corresponding `model.go` in the domain module to match the new schema:

```go
// internal/domain/release-notes/model.go
type ReleaseNote struct {
    database.BaseModel
    Title       string
    Tags        pq.StringArray `gorm:"type:text[]"`  // NEW
    // ...
}
```

### 5. Validate

```bash
go build -o ./tmp/main .    # Must compile
go vet ./...                # No issues
```

If the dev environment is running, verify the migration applied correctly.

## Naming Convention

Use snake_case describing the change:

| Change | Migration Name |
|--------|---------------|
| Add a column | `add_tags_to_release_notes` |
| Create a table | `create_categories_table` |
| Add an index | `add_index_on_release_notes_created_at` |
| Remove a column | `remove_legacy_flag_from_users` |
| Add a constraint | `add_unique_constraint_on_org_slug` |

## Commands Reference

```bash
make migrations-new name=<name>      # Create up/down SQL files
make migrations-up                    # Apply one pending migration
make migrations-up-all                # Apply all pending migrations
make migrations-down                  # Revert one migration
make migration-force version=<ver>    # Force to specific version
make migrations-unfuck                # Fix dirty migration state
```

## Recovery

If migrations get into a bad state (dirty flag):

```bash
make migrations-unfuck
```

This reads the current version, forces to that version, then rolls back one.

## Important Notes

- Migrations run automatically on production startup (`cfg.Env == "production"`)
- In development, run migrations manually with `make` commands
- The `up` and `down` must be perfectly symmetric
- Always update the GORM model after writing the migration
- Test both `up` and `down` locally before committing

## Anti-Patterns

| Don't | Why | Instead |
|-------|-----|---------|
| Skip the down migration | Can't rollback | Always write symmetric up/down |
| Update GORM model without migration | Schema drift | Migration first, then model |
| Write migration without testing down | May fail on rollback | Run `make migrations-down` to verify |
| Use GORM AutoMigrate | Uncontrolled schema changes | Always use explicit SQL migrations |
| Edit a migration already run in production | Breaks migration history | Create a new migration |
