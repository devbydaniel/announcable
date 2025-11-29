# Announcable

A self-hosted release notes and changelog platform. Create beautiful release notes, embed them in your product with a customizable widget, and host a public changelog page.

## Features

### Release Notes Management

- Create and edit release notes with title, description, release date, and media (images or YouTube/Loom embeds)
- Publish/unpublish release notes to control visibility
- Track engagement metrics: views, likes, and CTA clicks
- Configure per-release-note options:
  - Custom call-to-action text and URL
  - Attention mechanisms (indicator dot or instant-open on page load)
  - Visibility controls (hide on widget, hide on release page)

### Embeddable Widget

- Embed a customizable widget in your product with a single script tag
- Three widget types: Modal, Popover, or Sidebar
- Full styling control: colors, borders, border radius
- "New release" indicator shows users when there are updates they haven't seen
- Optional like/reaction feature for user engagement
- Works with any website—just add the script and trigger element

### Public Release Page

- Hosted changelog page at `yourinstance.com/s/your-org-slug`
- Customizable branding: logo, colors, title, description
- Optional back link to your main site
- Or disable the hosted page and build your own using the API

### Team Management

- Multi-user support with role-based access control
- Invite team members via email with Admin or Member roles
- Organization-based data isolation

## Tech Stack

- **Backend**: Go with Chi router
- **Database**: PostgreSQL with GORM
- **Frontend**: Server-rendered HTML templates with HTMX and Alpine.js
- **Widget**: Lit Web Components (lightweight, framework-agnostic)
- **Object Storage**: Minio (S3-compatible)
- **Email**: SMTP (any provider)

## Self-Hosting

Announcable is designed for self-hosting. No external services required except email delivery.

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for development)
- Node.js 18+ (for widget development)

### Quick Start

1. **Clone the repository**

   ```bash
   git clone https://github.com/devbydaniel/announcable.git
   cd announcable
   ```

2. **Configure environment**

   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Start services**

   ```bash
   cd backend
   make dev-services  # Starts PostgreSQL, Minio, Mailcatcher
   ```

4. **Run database migrations**

   ```bash
   make migrations-up-all
   ```

5. **Start the application**

   ```bash
   make dev-air  # Starts Go backend with hot-reload
   ```

6. **Access the application**
   - App: http://localhost:3000
   - Email (Mailcatcher): http://localhost:1080
   - Minio Console: http://localhost:9001

### Environment Variables

Create a `.env` file in the repository root:

```bash
# App Configuration
ENV=development
APP_ENVIRONMENT=self-hosted  # Use 'self-hosted' for self-hosted deployments
BASE_URL=localhost:3000
PORT=3000

# Database (PostgreSQL)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-secure-password
POSTGRES_NAME=announcable

# Object Storage (Minio/S3)
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=your-secure-secret
MINIO_ENDPOINT=localhost:9000
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# Email (SMTP)
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USER=your-smtp-user
SMTP_PASS=your-smtp-password
SMTP_TLS=true

# Legal Document Versions (increment when you update TOS/Privacy Policy)
TOS_VERSION=1
PP_VERSION=1
```

For production deployments, also set:

- `ENV=production`
- Use secure passwords
- Configure a production SMTP service
- Set up SSL/TLS termination (e.g., with a reverse proxy)

### Production Deployment

For production, use the provided Docker Compose files:

```bash
docker compose -f compose-base.yml -f compose-prod.yml up -d
```

You'll want to:

- Put a reverse proxy (nginx, Caddy, Traefik) in front for SSL termination
- Use a managed PostgreSQL or ensure database backups
- Configure production SMTP credentials
- Set strong, unique passwords for all services

## Widget Integration

After setting up Announcable and creating your first release notes:

1. Go to **Settings** in your dashboard
2. Copy your **Widget ID**

3. Add the widget to your site:

```html
<script
  src="https://your-announcable-instance.com/widget"
  data-org-id="YOUR_WIDGET_ID"
  data-anchor-query-selector="#changelog-button"
></script>
```

The widget attaches to elements matching your query selector. When users click those elements, the widget opens.

### Widget Options

| Attribute                    | Description                                                              |
| ---------------------------- | ------------------------------------------------------------------------ |
| `data-org-id`                | Required. Your widget ID from settings                                   |
| `data-anchor-query-selector` | CSS selector for trigger element(s). If not set, shows a floating button |
| `data-hide-indicator`        | Set to `"true"` to hide the "new release" indicator dot                  |

### Customization

Use the **Widget** page in your dashboard to customize:

- Widget type (modal, popover, sidebar)
- Colors and borders
- Title and description
- Like button text
- Call-to-action button text

## API

The widget fetches data from these public endpoints:

- `GET /api/release-notes/{orgId}` - Get published release notes
- `GET /api/widget-config/{orgId}` - Get widget configuration
- `GET /s/{orgSlug}` - Public release page

## Development

### Backend Development

```bash
cd backend

# Start supporting services
make dev-services

# Run with hot-reload
make dev-air

# In another terminal, for CSS/JS hot-reload:
npm run dev

# Database migrations
make migrations-new name=add_new_feature  # Create migration
make migrations-up                         # Run one migration
make migrations-up-all                     # Run all migrations
make migrations-down                       # Rollback one migration
```

### Widget Development

```bash
cd widget

npm install
npm run dev      # Development with hot-reload
npm run build    # Production build
npm run lint     # Run linter
```

### Project Structure

```
announcable/
├── backend/
│   ├── internal/
│   │   ├── handler/     # HTTP handlers
│   │   ├── domain/      # Domain models
│   │   ├── database/    # Database layer and migrations
│   │   ├── middleware/  # Auth, RBAC middleware
│   │   └── objstore/    # Minio wrapper
│   ├── templates/       # Go HTML templates
│   ├── static/          # Static assets (CSS, JS)
│   └── assets/          # Source CSS/JS (built by Vite)
├── widget/              # Lit web component embeddable widget
├── compose-base.yml     # Base Docker services
├── compose-dev.yml      # Development overrides
└── compose-prod.yml     # Production configuration
```

## License

This project is licensed under the [GNU Affero General Public License v3.0 (AGPL-3.0)](LICENSE.md).

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

To report security vulnerabilities, please see [SECURITY.md](SECURITY.md).
