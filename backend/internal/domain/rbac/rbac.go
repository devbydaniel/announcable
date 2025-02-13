package rbac

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
)

type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	PermissionManageAccess      Permission = "manage_access"
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
