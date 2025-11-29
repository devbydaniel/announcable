# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in Announcable, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, please email security concerns to: **security@announcable.me**

Include as much information as possible:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fixes (optional)

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 1 week
- **Resolution timeline**: Depends on severity, but we aim to fix critical issues within 30 days

## Supported Versions

We provide security updates for the latest release only. We recommend always running the most recent version.

## Security Best Practices for Self-Hosters

When self-hosting Announcable:

1. **Use strong passwords** for all services (PostgreSQL, Minio, SMTP)
2. **Enable HTTPS** using a reverse proxy (nginx, Caddy, Traefik)
3. **Keep dependencies updated** by regularly pulling the latest release
4. **Restrict network access** to admin ports (database, Minio console)
5. **Back up your database** regularly
6. **Use environment variables** for secrets, never commit them to version control
