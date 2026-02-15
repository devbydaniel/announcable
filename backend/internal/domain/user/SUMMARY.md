# User Accounts

User accounts and credentials management.

The `user` package manages user accounts for the platform. The core entity is `User`, which holds email, hashed password, and email verification status. All users embed `database.BaseModel` (UUID ID, timestamps).

**Key components:**
- `User` model with email uniqueness constraint and email verification flag
- `New(email, password)` constructor with email format validation
- `Service` for user CRUD operations (create, find by ID, find by email, update password, delete)
- `Repository` wrapping GORM for database access

**Integrations:**
- Referenced by `organisation.OrganisationUser` for org membership
- Referenced by `session.Session` via `UserID` for authentication
- Password hashing handled by the `password` utility package
- Email verification tokens managed in auth handlers

**Notes:**
- Package name is singular (`user`) — handler package uses plural (`users`)
- Email validation uses regex in the `New()` constructor
- Password is stored already hashed; hashing happens at the handler/service level
