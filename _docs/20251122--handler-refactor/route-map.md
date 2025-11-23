# Route Map & Refactoring Target

Source: `backend/main.go` (Nov 22, 2025)

This document maps existing Chi routes to their current handlers and defines the target structure for the refactor.

## 1. Public / Home
*   **Route**: `GET /`
*   **Current**: `handler.HandleHomePage` (`handle-home.go`)
*   **Target**: `pages/public/home/page.go`

## 2. Auth & Onboarding

### Login
*   **Route**: `GET /login/`
*   **Current**: `handler.HandleLoginPage` (`handle-login-page.go`)
*   **Target**: `pages/auth/login/page.go`
*   **Route**: `POST /login/`
*   **Current**: `handler.HandleLogin` (`handle-login.go`)
*   **Target**: `pages/auth/login/actions.go`

### Register
*   **Route**: `GET /register/`
*   **Current**: `handler.HandleRegisterPage` (`handle-register-page.go`)
*   **Target**: `pages/auth/register/page.go`
*   **Route**: `POST /register/`
*   **Current**: `handler.HandleRegister` (`handle-register.go`)
*   **Target**: `pages/auth/register/actions.go`

### Email Verify
*   **Route**: `GET /verify-email/`
*   **Current**: `handler.HandleEmailVerifyPage` (`handle-email-verify-page.go`)
*   **Target**: `pages/auth/verifyEmail/page.go`
*   **Route**: `POST /verify-email/`
*   **Current**: `handler.HandleEmailVerifyResend` (`handle-email-verify-resend.go`)
*   **Target**: `pages/auth/verifyEmail/actions.go`

### Invite Accept
*   **Route**: `GET /invite-accept/{token}/`
*   **Current**: `handler.HandleInviteAcceptPage` (`handle-invite-accept-page.go`)
*   **Target**: `pages/auth/inviteAccept/page.go`
*   **Route**: `POST /invite-accept/{token}/`
*   **Current**: `handler.HandleInviteAccept` (`handle-invite-accept.go`)
*   **Target**: `pages/auth/inviteAccept/actions.go`

### Password Forgot
*   **Route**: `GET /forgot-pw/`
*   **Current**: `handler.HandlePwForgotPage` (`handle-pw-forgot-page.go`)
*   **Target**: `pages/auth/passwordForgot/page.go`
*   **Route**: `POST /forgot-pw/`
*   **Current**: `handler.HandlePwForgot` (`handle-pw-forgot.go`)
*   **Target**: `pages/auth/passwordForgot/actions.go`

### Password Reset
*   **Route**: `GET /reset-pw/{token}/`
*   **Current**: `handler.HandlePwResetPage` (`handle-pw-reset-page.go`)
*   **Target**: `pages/auth/passwordReset/page.go`
*   **Route**: `POST /reset-pw/{token}/`
*   **Current**: `handler.HandlePwReset` (`handle-pw-reset.go`)
*   **Target**: `pages/auth/passwordReset/actions.go`

### Logout
*   **Route**: `GET /logout/`
*   **Current**: `handler.HandleLogout` (`handle-logout.go`)
*   **Target**: `pages/auth/logout/actions.go`

## 3. User Management & Profile

### Users
*   **Route**: `GET /users/`
*   **Current**: Migrated
*   **Target**: `pages/users/page.go` → `ServeUsersPage` ✅
*   **Route**: `DELETE /users/{id}`
*   **Current**: Migrated
*   **Target**: `pages/users/user_delete_action.go` → `HandleUserDelete` ✅

### Profile (MIGRATED TO SETTINGS)
*   **Route**: `PATCH /settings/password`
*   **Current**: Migrated
*   **Target**: `pages/settings/account/password_update_action.go` → `HandlePasswordUpdate` ✅

### Invites
*   **Route**: `POST /invites/`
*   **Current**: Migrated
*   **Target**: `pages/users/invite_create_action.go` → `HandleInviteCreate` ✅
*   **Route**: `DELETE /invites/{id}`
*   **Current**: Migrated
*   **Target**: `pages/users/invite_delete_action.go` → `HandleInviteDelete` ✅

**Note**: Each action has its own file for maximum clarity. Invite actions live in the users package since they're logically part of user management.

## 4. Release Notes (Dashboard)

*   **Route**: `GET /release-notes/`
*   **Current**: `handler.HandleReleaseNotesPage` (`handle-release-notes-page.go`)
*   **Target**: `pages/releaseNotes/list/page.go`
*   **Route**: `POST /release-notes/`
*   **Current**: `handler.HandleReleaseNoteCreate` (`handle-release-note-create.go`)
*   **Target**: `pages/releaseNotes/create/actions.go` (or `pages/releaseNotes/list/actions.go`)
*   **Route**: `GET /release-notes/new`
*   **Current**: `handler.HandleReleaseNoteCreatePage` (`handle-release-note-create-page.go`)
*   **Target**: `pages/releaseNotes/create/page.go`
*   **Route**: `GET /release-notes/{id}`
*   **Current**: `handler.HandleReleaseNotePage` (`handle-release-note-page.go`)
*   **Target**: `pages/releaseNotes/detail/page.go`
*   **Route**: `PATCH /release-notes/{id}`
*   **Current**: `handler.HandleReleaseNoteUpdate` (`handle-release-note-update.go`)
*   **Target**: `pages/releaseNotes/detail/actions.go`
*   **Route**: `DELETE /release-notes/{id}`
*   **Current**: `handler.HandleReleaseNoteDelete` (`handle-release-note-delete.go`)
*   **Target**: `pages/releaseNotes/detail/actions.go`
*   **Route**: `PATCH /release-notes/{id}/publish`
*   **Current**: `handler.HandleReleaseNotePublish` (`handle-release-note-publish.go`)
*   **Target**: `pages/releaseNotes/detail/actions.go`

## 5. Widget Config

*   **Route**: `GET /widget-config/`
*   **Current**: `handler.HandleWidgetPage` (`handle-widget-page.go`)
*   **Target**: `pages/widget/config/page.go`
*   **Route**: `PATCH /widget-config/`
*   **Current**: `handler.HandleWidgetUpdate` (`handle-widget-update.go`)
*   **Target**: `pages/widget/config/actions.go`
*   **Route**: `PATCH /settings/widget-id` (moved from `/widget-config/external-id`)
*   **Current**: Migrated
*   **Target**: `pages/settings/account/widget_id_regenerate_action.go` → `HandleWidgetIdRegenerate` ✅
*   **Route**: `PATCH /settings/release-page-url` (moved from `/widget-config/base-url`)
*   **Current**: Migrated
*   **Target**: `pages/settings/account/release_page_url_update_action.go` → `HandleReleasePageUrlUpdate` ✅
*   **Note**: These actions were moved from `/widget-config` to `/settings` because they're displayed on the settings page and are user-facing configuration options

## 6. Release Page Config

*   **Route**: `GET /release-page-config/`
*   **Current**: `handler.HandleReleasePageConfigPage` (`handle-landing-page-config-page.go`)
*   **Target**: `pages/releasePage/config/page.go`
*   **Route**: `PATCH /release-page-config/`
*   **Current**: `handler.HandleReleasePageConfigUpdate` (`handle-lp-config-update.go`)
*   **Target**: `pages/releasePage/config/actions.go`

## 7. Settings

*   **Route**: `GET /settings/`
*   **Current**: Migrated
*   **Target**: `pages/settings/account/page.go` → `ServeSettingsPage` ✅
*   **Route**: `PATCH /settings/password`
*   **Current**: Migrated
*   **Target**: `pages/settings/account/password_update_action.go` → `HandlePasswordUpdate` ✅
*   **Route**: `PATCH /settings/widget-id`
*   **Current**: Migrated (was `/widget-config/external-id`)
*   **Target**: `pages/settings/account/widget_id_regenerate_action.go` → `HandleWidgetIdRegenerate` ✅
*   **Route**: `PATCH /settings/release-page-url`
*   **Current**: Migrated (was `/widget-config/base-url`)
*   **Target**: `pages/settings/account/release_page_url_update_action.go` → `HandleReleasePageUrlUpdate` ✅

## 8. Admin Dashboard

*   **Route**: `GET /admin/`
*   **Current**: `handler.HandleAdminDashboard` (`handle-admin-dashboard.go`)
*   **Target**: `pages/admin/dashboard/page.go`
*   **Route**: `GET /admin/organisations/{orgId}`
*   **Current**: `handler.HandleAdminOrgDetails` (`handle-admin-org-details.go`)
*   **Target**: `pages/admin/organisation/page.go`
*   **Route**: `PATCH /admin/organisations/{orgId}`
*   **Current**: `handler.HandleAdminOrgUpdate` (`handle-admin-org-update.go`)
*   **Target**: `pages/admin/organisation/actions.go`
*   **Route**: `PATCH /admin/organisations/{orgId}/release-page`
*   **Current**: `handler.HandleAdminOrgReleasePageUpdate` (`handle-admin-org-release-page-update.go`)
*   **Target**: `pages/admin/organisation/actions.go`
*   **Route**: `POST /admin/organisations/{orgId}/subscriptions`
*   **Current**: `handler.HandleAdminOrgSubscriptionCreate` (`handle-admin-org-subscription-create.go`)
*   **Target**: `pages/admin/organisation/actions.go`
*   **Route**: `DELETE /admin/organisations/{orgId}/subscriptions/{id}`
*   **Current**: `handler.HandleAdminOrgSubscriptionDelete` (`handle-admin-org-subscription-delete.go`)
*   **Target**: `pages/admin/organisation/actions.go`

## 9. Public Release Page & Widget Script

*   **Route**: `GET /s/{orgSlug}`
*   **Current**: `handler.HandleReleasePage` (`handle-release-note-page.go` - WAIT, verify: `handle-release-notes-website-serve.go` probably?)
    *   *Correction*: `handle-release-note-page.go` is usually single note. `HandleReleasePage` serves the public release notes list. The file is likely `handle-release-notes-website-serve.go` or similar. Checked file list: `handle-release-notes-website-serve.go` exists. `handle-release-note-page.go` is for the internal dashboard.
    *   **Target**: `pages/public/releasePage/page.go`
*   **Route**: `GET /widget/`
*   **Current**: `handler.HandleWidgetjsServe` (`handle-widgetjs-serve.go`)
*   **Target**: `pages/public/widgetScript/handler.go`

## 10. API

*   **Route**: `/api/release-notes/{orgId}`
*   **Current**: `handler.HandleReleaseNotesServe` (`handle-release-notes-widget-serve.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `/api/release-notes/{orgId}/status`
*   **Current**: `handler.HandleReleaseNotesStatusServe` (`handle-release-notes-status-serve.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `/api/release-notes/{orgId}/metrics`
*   **Current**: `handler.HandleReleaseNoteMetricCreate` (`handle-release-note-metrics.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `/api/release-notes/{orgId}/{releaseNoteId}/like`
*   **Current**: `handler.HandleGetReleaseNoteLikeState` (`handle-get-release-note-like-state.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `POST .../like`
*   **Current**: `handler.HandleReleaseNoteToggleLike` (`handle-release-note-likes.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `/api/widget-config/{orgId}`
*   **Current**: `handler.HandleWidgetConfigServe` (`handle-widget-config-serve.go`)
*   **Target**: `api/widget/handlers.go`
*   **Route**: `/api/img/*`
*   **Current**: `handler.HandleObjStore` (`handle-objstore.go`)
*   **Target**: `api/shared/handlers.go`

## 11. Payment / Stripe

*   **Route**: `/payment/create-checkout-session`
*   **Current**: `handler.HandleCheckoutSession` (`handle-subscription.go`)
*   **Target**: `pages/payment/actions.go`
*   **Route**: `/payment/create-portal-session`
*   **Current**: `handler.HandlePortalSession` (`handle-subscription.go`)
*   **Target**: `pages/payment/actions.go`
*   **Route**: `/subscription/confirm`
*   **Current**: `handler.HandleSubscriptionConfirm` (`handle-subscription-confirm.go`)
*   **Target**: `pages/public/subscription/page.go`
*   **Route**: `/subscription/cancel`
*   **Current**: `handler.HandleSubscriptionCancel` (`handle-subscription-cancel.go`)
*   **Target**: `pages/public/subscription/page.go`
*   **Route**: `/stripe/webhook`
*   **Current**: `handler.HandleWebhook` (`handle-subscription.go` likely, or separate) - Checked list: `handle-subscription.go` seems to cover multiple? No, `handle-subscription.go` is likely generic. There isn't a specific `handle-webhook.go`. Need to check content of `handle-subscription.go`.
*   **Target**: `api/stripe/handlers.go`

## 12. Other
*   **Route**: `NotFound`
*   **Current**: `handler.HandleNotFound` (`handle-not-found.go`)
*   **Target**: `pages/shared/errors/handlers.go` or similar.
