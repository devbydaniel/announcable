# Platform Administration

Super admin functionality for platform-wide management.

The `admin` package provides super admin capabilities. Unlike RBAC roles, the super admin is identified by a config-defined user ID — not a database role.

**Key components:**
- `AdminUser` struct — Holds the admin user ID
- `IsAdmin(userID, adminUserID)` — Checks if a user matches the configured admin
- `Service` — Platform-wide queries (list organisations, org details, stats)
- `Repository` — Cross-organisation GORM queries

**Integrations:**
- `middleware.AuthorizeSuperAdmin` gates admin routes
- Admin dashboard handler (`pages/admin/dashboard`) shows platform stats
- Admin org handler (`pages/admin/organisation`) manages individual orgs
- Admin user ID is set via environment configuration

**Notes:**
- Super admin is NOT part of the RBAC system — it's a separate mechanism
- The admin user ID is configured in `.env`, not stored in a database role
- Admin can view and manage all organisations regardless of membership
