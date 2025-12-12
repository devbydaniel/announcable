# Reverse Proxy Configuration with nginx (Subdomain Deployment)

This guide covers configuring nginx as a reverse proxy to serve Announcable at a subdomain (e.g., `release-notes.your-domain.com`).

## Overview

Announcable exposes several entry points that nginx will proxy:

| Entry Point | Path | Description |
|-------------|------|-------------|
| **Backend Admin UI** | `/`, `/login`, `/release-notes`, `/settings`, etc. | Dashboard for managing release notes |
| **Widget Script** | `/widget` | JavaScript file embedded on customer websites |
| **Widget API** | `/api/*` | REST endpoints the widget calls to fetch data |
| **Public Release Page** | `/s/{slug}` | Public changelog page visitors can view |

All paths are served by the same Announcable application, so nginx proxies everything to the internal app port.

---

## Prerequisites

Before starting, ensure you have:

- [ ] Announcable running via Docker Compose (on an internal port, e.g., 8080)
- [ ] nginx installed on your host server
- [ ] SSL certificate for your subdomain (Let's Encrypt/Certbot recommended)
- [ ] DNS A record pointing your subdomain to your server's IP address

---

## Environment Configuration

Update your `.env` file with the following settings:

```bash
# Your public subdomain (no trailing slash, no path)
BASE_URL=https://release-notes.your-domain.com

# Internal port the app listens on
PORT=8080
```

### What `BASE_URL` Affects

The `BASE_URL` environment variable is used to generate absolute URLs in:

- Password reset email links
- User invitation URLs
- Email template links
- OAuth callback URLs (if applicable)

**Important**: Always set this to your public-facing URL (the subdomain), not the internal port.

---

## nginx Configuration

Create the following nginx configuration file:

```nginx
server {
    listen 443 ssl http2;
    server_name release-notes.your-domain.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/release-notes.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/release-notes.your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # Proxy all requests to Announcable
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Allow larger uploads for release note images
    client_max_body_size 10M;
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name release-notes.your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

---

## Setup Steps

### 1. Configure DNS

Add an A record for your subdomain pointing to your server's IP address:

```
release-notes.your-domain.com → YOUR_SERVER_IP
```

Wait for DNS propagation (can take up to 48 hours, but usually much faster).

### 2. Obtain SSL Certificate

Use Certbot to obtain a free Let's Encrypt certificate:

```bash
sudo certbot certonly --nginx -d release-notes.your-domain.com
```

Or, if nginx isn't running yet:

```bash
sudo certbot certonly --standalone -d release-notes.your-domain.com
```

### 3. Create nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/announcable
```

Paste the nginx configuration from above, replacing `release-notes.your-domain.com` with your actual subdomain.

### 4. Enable the Site

```bash
# Create symlink to enable the site
sudo ln -s /etc/nginx/sites-available/announcable /etc/nginx/sites-enabled/

# Test the configuration
sudo nginx -t

# Reload nginx to apply changes
sudo systemctl reload nginx
```

### 5. Update `.env`

Edit your `.env` file:

```bash
BASE_URL=https://release-notes.your-domain.com
```

### 6. Restart Announcable

```bash
docker compose down && docker compose up -d
```

### 7. Verify

Visit `https://release-notes.your-domain.com` in your browser. You should see the Announcable login page.

---

## Widget Integration

Once deployed, customers can embed the widget on their websites:

```html
<script
  src="https://release-notes.your-domain.com/widget"
  data-org-id="YOUR_WIDGET_ID"
  data-anchor-query-selector="#changelog-button"
></script>
```

**Finding the Widget ID**: The Widget ID (also called Org ID) is found in the **Settings** page of the admin dashboard.

### Widget Attributes

| Attribute | Required | Description |
|-----------|----------|-------------|
| `data-org-id` | Yes | Your organization's widget ID |
| `data-anchor-query-selector` | No | CSS selector for the element that triggers the widget |

---

## Public Release Page

Announcable provides a public changelog page for each organization:

- **URL Pattern**: `https://release-notes.your-domain.com/s/{org-slug}`
- **Example**: `https://release-notes.your-domain.com/s/acme-corp`

The organization slug is configured in **Release Page Config** in the admin settings.

---

## Widget Backend URL Configuration (Required for Self-Hosted)

> ⚠️ **Important**: The widget has a hardcoded backend URL that must be updated for self-hosted deployments.

### File to Modify

`widget/src/lib/config.ts`

### Current Configuration

```typescript
export const backendUrl =
  import.meta.env.MODE === "production"
    ? "https://app.announcable.com"  // ← Default production URL
    : "http://localhost:3000";
```

### Update for Your Domain

```typescript
export const backendUrl =
  import.meta.env.MODE === "production"
    ? "https://release-notes.your-domain.com"  // ← Your subdomain
    : "http://localhost:3000";
```

### Steps to Update

1. **Edit the config file**:
   ```bash
   cd widget
   nano src/lib/config.ts
   ```

2. **Build the widget**:
   ```bash
   npm run build
   ```

3. **Copy to backend static folder**:
   ```bash
   cp dist/widget.js ../backend/static/widget/
   ```

4. **Rebuild and restart the app**:
   ```bash
   cd ..
   docker compose down
   docker compose build
   docker compose up -d
   ```

---

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| 502 Bad Gateway | App not running | Check `docker compose ps`, ensure app container is healthy |
| Static assets 404 | Incorrect proxy config | Verify `proxy_pass` URL and port match your `PORT` env var |
| Widget shows no data | Wrong backend URL in widget | Update `widget/src/lib/config.ts` and rebuild widget |
| Login redirects fail | Wrong `BASE_URL` | Update `.env` with correct subdomain and restart app |
| Email links broken | Wrong `BASE_URL` | Update `.env` with correct subdomain and restart app |
| Upload fails | Body size limit | Increase `client_max_body_size` in nginx config |
| SSL certificate errors | Cert path mismatch | Verify certificate paths in nginx config match Certbot output |
| Connection refused | Firewall blocking | Ensure ports 80 and 443 are open; verify app is listening on configured `PORT` |

### Debugging Commands

```bash
# Check if Announcable is running
docker compose ps

# View Announcable logs
docker compose logs -f app

# Check nginx status
sudo systemctl status nginx

# Test nginx config
sudo nginx -t

# Check if port is listening
sudo netstat -tlnp | grep 8080
```

---

## Security Recommendations

### Essential

- **Always use HTTPS** in production
- **Keep nginx updated** with security patches
- **Renew SSL certificates** before expiration (Certbot auto-renewal recommended)
- **Firewall configuration**: Only expose ports 80 and 443; never expose the internal app port (8080) to the public internet

### Optional Security Headers

Add these headers to your nginx configuration for enhanced security:

```nginx
server {
    # ... existing config ...
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # ... rest of config ...
}
```

### Rate Limiting (Recommended)

To protect your API endpoints from abuse, add rate limiting:

```nginx
# Add to http block (usually in /etc/nginx/nginx.conf)
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

# Add to your server block
location /api/ {
    limit_req zone=api_limit burst=20 nodelay;
    
    proxy_pass http://127.0.0.1:8080;
    # ... rest of proxy settings ...
}
```

### Certbot Auto-Renewal

Ensure Certbot's auto-renewal is active:

```bash
# Test renewal
sudo certbot renew --dry-run

# Check renewal timer
sudo systemctl status certbot.timer
```

---

## Full Configuration Example

Here's a complete nginx configuration with all recommended settings:

```nginx
# Rate limiting zone (add to http block in /etc/nginx/nginx.conf)
# limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

server {
    listen 443 ssl http2;
    server_name release-notes.your-domain.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/release-notes.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/release-notes.your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Allow larger uploads for release note images
    client_max_body_size 10M;
    
    # API rate limiting (uncomment if rate limit zone is configured)
    # location /api/ {
    #     limit_req zone=api_limit burst=20 nodelay;
    #     
    #     proxy_pass http://127.0.0.1:8080;
    #     proxy_http_version 1.1;
    #     proxy_set_header Host $host;
    #     proxy_set_header X-Real-IP $remote_addr;
    #     proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    #     proxy_set_header X-Forwarded-Proto $scheme;
    # }
    
    # Proxy all requests to Announcable
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name release-notes.your-domain.com;
    return 301 https://$server_name$request_uri;
}
```
