# Multi-Tenant Organizations

Organization management with user membership and invites.

The `organisation` package manages the multi-tenant structure. Each organisation has a name, a public `ExternalID` (UUID, auto-generated on create), and serves as the scoping boundary for all release notes, configs, and user access.

**Key entities:**
- `Organisation` — Core tenant entity with name and external UUID
- `OrganisationUser` — Join table linking users to organisations with an RBAC `Role`
- `OrganisationInvite` — Pending invitations with email, role, expiry, and external token

**Key components:**
- `New(name)` constructor with 3-character minimum validation
- `Connect(org, user, role)` creates an `OrganisationUser` association
- `Service` for org CRUD, user membership management, invite lifecycle
- `Repository` wrapping GORM for database access

**Integrations:**
- `rbac.Role` defines the role assigned to each org member (admin, manager)
- `user.User` referenced in `OrganisationUser` for membership
- Release notes, widget configs, release page configs, metrics, and likes all scope to an organisation via `OrganisationID`
- `ExternalID` is the public-facing org identifier used in widget API endpoints and embed scripts

**Notes:**
- Package name is singular (`organisation`)
- `ExternalID` is auto-generated in `BeforeCreate` GORM hook
- Invites use an `ExternalID` string token (not UUID) for URL-safe invite links
