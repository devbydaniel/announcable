# Migration Checklist

- [x] **Refactor Shared Dependencies**
  - [x] Create `backend/internal/handler/shared/dependencies.go` (Dependency container)
  - [x] Create `backend/internal/handler/shared/viewdata.go` (BaseTemplateData)
  - [x] Refactor `handler.go` to use `shared.Dependencies` internally with backward compatibility wrapper
  - [x] Add TODO comment in `backend/main.go` (will update to use `shared.New()` in final step)

- [x] **Release Notes Feature**
  - [x] Create `backend/internal/handler/pages/release_notes/` structure
  - [x] Migrate List (`HandleReleaseNotesPage` → `pages/release_notes/list/page.go`)
  - [x] Migrate Create (`HandleReleaseNoteCreate` → `pages/release_notes/create/release_note_create_action.go`, `HandleReleaseNoteCreatePage` → `pages/release_notes/create/page.go`)
  - [x] Migrate Details/Update (`HandleReleaseNotePage` → `pages/release_notes/detail/page.go`, `HandleReleaseNoteUpdate` → `pages/release_notes/detail/release_note_update_action.go`, `HandleReleaseNoteDelete` → `pages/release_notes/detail/release_note_delete_action.go`, `HandleReleaseNotePublish` → `pages/release_notes/detail/release_note_publish_action.go`)
  - [x] Update `backend/main.go` routes for release notes
  - [x] Delete old handler files

- [x] **Auth Feature**
  - [x] Create `backend/internal/handler/pages/auth/` structure
  - [x] Migrate Login (`pages/auth/login/page.go`, `pages/auth/login/login_action.go`)
  - [x] Migrate Register (`pages/auth/register/page.go`, `pages/auth/register/register_action.go`)
  - [x] Migrate Email Verify (`pages/auth/verifyEmail/page.go`, `pages/auth/verifyEmail/resend_action.go`)
  - [x] Migrate Invite Accept (`pages/auth/inviteAccept/page.go`, `pages/auth/inviteAccept/accept_action.go`)
  - [x] Migrate Password Forgot/Reset (`pages/auth/passwordForgot/`, `pages/auth/passwordReset/`)
  - [x] Migrate Logout (`pages/auth/logout/logout_action.go`)
  - [x] Update `backend/main.go` routes for auth
  - [x] Delete old auth handler files

- [x] **User Feature**
  - [x] Create `backend/internal/handler/pages/users/` directory
  - [x] Create `pages/users/page.go` (Handlers struct, ServeUsersPage)
  - [x] Create `pages/users/user_delete_action.go` (HandleUserDelete)
  - [x] Create `pages/users/invite_create_action.go` (HandleInviteCreate)
  - [x] Create `pages/users/invite_delete_action.go` (HandleInviteDelete)
  - [x] Update `backend/main.go` routes
  - [x] Delete old handler files from flat structure
  - [x] Apply one-file-per-action naming convention

- [x] **Widget Config Feature**
  - [x] Create `backend/internal/handler/pages/widget/config/`
  - [x] Create `page.go` (Handlers struct, ServeWidgetConfigPage)
  - [x] Create `config_update_action.go` (HandleConfigUpdate)
  - [x] Update `backend/main.go` routes
  - [x] Delete old handler files

- [x] **Release Page Config Feature**
  - [x] Create `backend/internal/handler/pages/releasePage/config/`
  - [x] Create `page.go` (Handlers struct, ServeReleasePageConfigPage)
  - [x] Create `config_update_action.go` (HandleConfigUpdate)
  - [x] Update `backend/main.go` routes
  - [x] Delete old handler files

- [x] **Settings Feature**
  - [x] Create `backend/internal/handler/pages/settings/account/`
  - [x] Migrate Settings Page (`ServeSettingsPage`)
  - [x] Move Password Update action from `profile/` to `settings/account/`
  - [x] Migrate Widget ID Regenerate action from `/widget-config/external-id` to `/settings/widget-id`
  - [x] Migrate Release Page URL Update action from `/widget-config/base-url` to `/settings/release-page-url`
  - [x] Update `backend/main.go` routes
  - [x] Update template routes
  - [x] Delete old files (`handle-settings-page.go`, `handle-widget-extid-regenerate.go`, `handle-lp-baseurl-update.go`, and `pages/profile/` directory)

- [x] **Admin Feature**
  - [x] Create `backend/internal/handler/pages/admin/`
  - [x] Migrate Dashboard (`pages/admin/dashboard/page.go`)
  - [x] Migrate Organisation Details & Actions (`pages/admin/organisation/`)
    - [x] `page.go` - ServeOrganisationDetailsPage
    - [x] `org_update_action.go` - HandleOrgUpdate
    - [x] `release_page_update_action.go` - HandleReleasePageUpdate
    - [x] `subscription_create_action.go` - HandleSubscriptionCreate
    - [x] `subscription_delete_action.go` - HandleSubscriptionDelete
  - [x] Update `backend/main.go` routes

- [x] **Public Pages**
  - [x] Create `backend/internal/handler/pages/public/`
  - [x] Migrate Home (`pages/public/home/page.go`)
  - [x] Migrate Public Release Page (`pages/public/release_page/page.go`)
  - [x] Migrate Widget Script (`pages/public/widget_script/page.go`)
  - [x] Migrate Subscription Confirm/Cancel (`pages/public/subscription/`)
    - [x] `page.go` - Handlers struct
    - [x] `confirm_action.go` - HandleSubscriptionConfirm
    - [x] `cancel_action.go` - HandleSubscriptionCancel
  - [x] Update `backend/main.go` routes

- [x] **API & Webhooks**
  - [x] Create `backend/internal/handler/pages/api/`
  - [x] Migrate Widget API endpoints (`pages/api/widget/`)
    - [x] `handlers.go` - Handlers struct
    - [x] `release_notes.go` - HandleReleaseNotesServe
    - [x] `status.go` - HandleReleaseNotesStatusServe
    - [x] `metrics.go` - HandleReleaseNoteMetricCreate
    - [x] `like_state.go` - HandleGetReleaseNoteLikeState
    - [x] `toggle_like.go` - HandleReleaseNoteToggleLike
    - [x] `config.go` - HandleWidgetConfigServe
  - [x] Migrate Shared API (`pages/api/shared/`)
    - [x] `handlers.go` - Handlers struct
    - [x] `objstore.go` - HandleObjStore
    - [x] `not_found.go` - HandleNotFound
  - [x] Migrate Stripe Webhook (`pages/api/stripe/`)
    - [x] `handlers.go` - Handlers struct
    - [x] `webhook.go` - HandleWebhook
  - [x] Migrate Payment handlers (`pages/payment/`)
    - [x] `handlers.go` - Handlers struct
    - [x] `checkout_session_action.go` - HandleCheckoutSession
    - [x] `portal_session_action.go` - HandlePortalSession
  - [x] Update `backend/main.go` routes

- [x] **Cleanup Phase 1: Remove Old Handlers**
  - [x] Remove old `backend/internal/handler/handle-*.go` files
  - [x] Run `go fmt ./...`
  - [x] Run `go mod tidy`

- [x] **Cleanup Phase 2: Remove Backward Compatibility Wrapper**
  - [x] Delete `backend/internal/handler/handler.go` (the old Handler wrapper)
  - [x] Update `backend/main.go` to use `shared.New(db, objStore)` directly (was already done)
  - [x] Remove TODO comment from `backend/main.go` (was already done)
  - [x] Verify no remaining references to `handler.Handler` or `handler.NewHandler`
  - [x] This removes the temporary field duplication (`log`, `decoder`) from Step 2

- [x] **Verification**
  - [x] Verify no compilation errors (build successful)
