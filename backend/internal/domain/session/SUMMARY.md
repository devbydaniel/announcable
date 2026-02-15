# Session Authentication

Session-based authentication for the web application.

The `session` package manages user sessions. Sessions link a user to a browser via a cookie containing an external session ID. Sessions have an expiry timestamp.

**Key components:**
- `Session` model with `UserID`, `ExpiresAt` (Unix millis), and `ExternalID` (cookie value)
- `New(userId, expiresAt, sessionId)` constructor
- `AuthCookieName` constant: `"announcable-session"`
- `Service` for session creation, lookup by external ID, and deletion
- `Repository` wrapping GORM for database access

**Integrations:**
- `middleware.Authenticate` reads the session cookie and validates against this module
- Used by login/logout handlers to create and destroy sessions
- `UserID` references `user.User`

**Notes:**
- Session ID stored in cookie is the `ExternalID`, not the database UUID
- Expiry is stored as Unix milliseconds
