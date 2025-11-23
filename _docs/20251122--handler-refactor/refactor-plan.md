# Handler Refactoring Plan

**Goal**: Group every backend HTTP handler by page/feature rather than having >50 flat `handle-*.go` files in `backend/internal/handler/`, so that the page-level view handler and its related action handlers live together.

## 1. Catalogue the current routes & handlers

- [x] Open `backend/main.go` and copy the existing Chi route groups (login, release notes, widget config, admin, API, etc.) into a new doc (e.g. `_docs/20251122--handler-refactor/route-map.md`). This captures which handlers serve each page/action today so nothing is lost when files move. Cite `backend/main.go` to keep the doc authoritative.
- [x] For each route group, list the concrete handler functions it uses (e.g. `/release-notes` uses `HandleReleaseNotesPage`, `HandleReleaseNoteCreate`, etc. from `handle-release-notes-*.go`). This list will drive the new folder layout.

## 2. Extract shared handler dependencies into a reusable container

- [x] Create `backend/internal/handler/shared/dependencies.go` with a `Dependencies` struct holding `DB`, `ObjStore`, `Log`, and `Decoder`.
- [x] Move `BaseTemplateData` (and the `PageData` interface) into a new `backend/internal/handler/shared/viewdata.go` file so page packages can import it without reaching back into the old flat package.
- [x] Refactor `backend/internal/handler/handler.go` to use `shared.Dependencies` internally while maintaining backward compatibility with existing handlers.
- [x] Keep existing `Handler` struct as a temporary wrapper (with lowercase `log` and `decoder` fields) to avoid breaking 52+ existing handler files.
- [x] Add TODO comment in `main.go` indicating future migration to `shared.New()`.
- **Note**: The `Handler` wrapper creates minimal duplication (two pointer copies) but ensures zero-risk incremental migration. This will be removed in Step 8 after all handlers are migrated.

## 3. Introduce a feature/page directory structure

- [x] Create `backend/internal/handler/pages/` and, under it, one folder per top-level route group from Step 1 (e.g. `auth`, `releaseNotes`, `widget`, `releasePageConfig`, `settings`, `admin`, `api`, `stripe`, `public`).
- [x] Within each feature folder, create subfolders for distinct pages when there are multiple templates/actions (e.g. `pages/releaseNotes/list`, `pages/releaseNotes/detail`, `pages/auth/login`, `pages/auth/register`, etc.).
- [x] In every subfolder add files by convention:
  - `page.go` – contains the page-serving handler plus template-specific data structs and shared Handlers struct
  - `{domain}_actions.go` – contains POST/PATCH/DELETE handlers for a specific domain (e.g., `user_actions.go`, `invite_actions.go`)
  - When multiple domains coexist, split actions into separate files prefixed by domain name
  - All action handlers are methods on the Handlers struct defined in `page.go`
  - Document this convention in `backend/internal/handler/pages/AGENTS.md`

## 4. Move code into the new packages feature-by-feature

- [ ] **Release Notes**: Start with the `/release-notes` group because it has the most handlers.
  - Create `backend/internal/handler/pages/releaseNotes/list/page.go` and move `HandleReleaseNotesPage` plus its helper structs.
  - Relocate actions: `HandleReleaseNoteCreatePage`, `HandleReleaseNoteCreate`, `HandleReleaseNoteUpdate`, `HandleReleaseNoteDelete`, `HandleReleaseNotePublish`.
  - Adjust imports to pull from `shared`.
- [ ] **Auth & Onboarding**: (`/login`, `/register`, `/verify-email`, `/invite-accept`, `/forgot-pw`, `/reset-pw`, `/logout`) -> `pages/auth/login`, `pages/auth/register`, etc.
- [x] **Users & Invites**:
  - `pages/users/page.go` - ServeUsersPage with UserData and InviteData
  - `pages/users/user_delete_action.go` - HandleUserDelete
  - `pages/users/invite_create_action.go` - HandleInviteCreate
  - `pages/users/invite_delete_action.go` - HandleInviteDelete
  - One file per action for maximum clarity
- [ ] **Widget Configuration**: (`/widget-config` + PATCH routes) -> `pages/widget/config`.
- [x] **Release Page Config**: -> `pages/releasePage/config`.
- [x] **Settings**: -> `pages/settings/account/`.
  - `pages/settings/account/page.go` - ServeSettingsPage (account settings overview)
  - `pages/settings/account/password_update_action.go` - HandlePasswordUpdate (moved from profile/)
  - `pages/settings/account/widget_id_regenerate_action.go` - HandleWidgetIdRegenerate (moved from /widget-config/external-id)
  - `pages/settings/account/release_page_url_update_action.go` - HandleReleasePageUrlUpdate (moved from /widget-config/base-url)
- [ ] **Admin Dashboard**: -> `pages/admin/dashboard` and `pages/admin/organisation`.
- [ ] **Public Pages**: (`/`, `/subscription/confirm`, `/subscription/cancel`, `/s/{orgSlug}`) -> `pages/public/...`.
- [ ] **API Endpoints**: (`/api/...`, `/widget`, `/img`, `/stripe/webhook`) -> `backend/internal/handler/api/<feature>/handlers.go`.
  - While moving:
    - Replace `func (h *Handler) HandleX` with `func (h *SpecificPage) Serve`.
    - Ensure helper structs move with their page.
    - Update template constructors.

### Naming Conventions (Idiomatic Go)

**Directory/Package Names:**
- Use lowercase, ideally single word (`login`, `logout`, `register`)
- For multi-word: use `snake_case` (`invite_accept`, `password_reset`, `verify_email`)
- Never use camelCase or kebab-case for directories

**File Names:**
- Always use `snake_case` (`password_update_action.go`, `user_delete_action.go`)

**Page Files:**
- `page.go` - Contains page-serving handler (GET), template data structures, and Handlers struct with New() constructor

**Action Files:**
- `{domain}_{action}_action.go` - One file per action (singular "action")
- Examples:
  - `user_delete_action.go` - HandleUserDelete
  - `invite_create_action.go` - HandleInviteCreate
  - `invite_delete_action.go` - HandleInviteDelete
  - `password_update_action.go` - HandlePasswordUpdate
  - `subscription_cancel_action.go` - HandleSubscriptionCancel

**Benefits:**

1. **Idiomatic Go**: Follows standard Go conventions seen in stdlib and popular projects
2. **Ultimate Clarity**: Each file contains exactly one action
3. **Easy to Find**: Action name is in the filename
4. **Perfect for Teams**: Zero merge conflicts - each developer works on separate files
5. **Self-Documenting**: File structure mirrors the API surface
6. **Simple**: No decisions about when to split or group - one action, one file

## 5. Provide constructors and registration helpers

- [ ] In each feature package expose a constructor like `func New(deps *shared.Dependencies) *Handlers`.
- [ ] Optionally add a `RegisterRoutes(r chi.Router)` helper per feature.

## 6. Update `backend/main.go`

- [ ] Replace `handler.NewHandler` with shared deps and feature-specific initializations.
- [ ] Swap every `handler.HandleX` reference with the appropriate feature struct method.

## 7. Delete the obsolete flat files and run gofmt

- [ ] Remove original `handle-*.go` files.
- [ ] Run `gofmt` / `goimports`.

## 8. Remove backward compatibility wrapper

- [ ] Delete `backend/internal/handler/handler.go` (the old Handler struct wrapper).
- [ ] Update `backend/main.go` to use `shared.New(db, objStore)` directly.
- [ ] Verify no remaining references to old `handler.Handler` or `handler.NewHandler`.
- [ ] This removes the temporary duplication of `log` and `decoder` fields introduced in Step 2.

## 9. Regression testing & verification

- [ ] Run `go test ./...`.
- [ ] Start dev server (`make dev-start`) and manually verify key pages.
- [ ] Update documentation (`CLAUDE.md` or new docs).
