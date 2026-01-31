package rbac

// Role represents a user's role within an organisation.
type Role string

func (r Role) String() string {
	return string(r)
}

const (
	// RoleAdmin is the administrator role with full permissions.
	RoleAdmin Role = "admin"
	// RoleManager is the manager role with limited permissions.
	RoleManager Role = "manager"
)

// Permission represents an action that can be performed within an organisation.
type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	// PermissionManageAccess allows managing user access within an organisation.
	PermissionManageAccess Permission = "manage_access"
	// PermissionManageReleaseNote allows creating, editing, and deleting release notes.
	PermissionManageReleaseNote Permission = "manage_release_note"
)

var rolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionManageAccess,
		PermissionManageReleaseNote,
	},
	RoleManager: {
		PermissionManageReleaseNote,
	},
}

// HasPermission reports whether the given role has the specified permission.
func HasPermission(role Role, permission Permission) bool {
	permissions, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}
