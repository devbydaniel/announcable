# Role-Based Access Control

RBAC system for organization-level permissions.

The `rbac` package defines roles and permissions for controlling access to organization resources. It is a pure logic package with no database models or services.

**Roles:**
- `RoleAdmin` — Full access (manage access + manage release notes)
- `RoleManager` — Release note management only

**Permissions:**
- `PermissionManageAccess` — User management, invites, settings
- `PermissionManageReleaseNote` — Release note CRUD, widget/release page config

**Key function:**
- `HasPermission(role, permission) bool` — Checks if a role has a specific permission

**Integrations:**
- `organisation.OrganisationUser.Role` stores the user's role per org
- `middleware.Authorize(permission)` enforces permissions on routes
- `middleware.AuthorizeSuperAdmin` is separate from RBAC (checks config-defined admin user ID)

**Notes:**
- This is a standalone package with no dependencies on other domain modules
- Role-permission mapping is hardcoded in `rolePermissions` map
- Super admin is NOT part of RBAC — it's a separate config-based check in the `admin` domain
