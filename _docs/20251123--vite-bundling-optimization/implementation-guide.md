# Implementation Guide: Vite CSS/JS Bundling

✅ **IMPLEMENTATION COMPLETE - 2025-11-23**

This guide was used to implement the CSS/JS bundling optimization. All phases have been completed successfully.

**Final State:**
- All CSS moved to `assets/css/` (base, components, layouts, pages)
- 36 bundles generated in `static/dist/`
- `static/css/` directory removed
- See `assets/css/AGENTS.md` for current architecture

---

**Original Implementation Guide (for reference):**

This guide provides detailed step-by-step instructions for implementing the CSS/JS bundling optimization.

## Prerequisites

- Node.js and npm installed
- Vite already in package.json (currently configured for Tailwind)
- Existing CSS architecture in `backend/static/css/`
- All templates using the base/layouts/pages structure

## Phase 1: Configure Vite

### Step 1.1: Update Vite Configuration

**File:** `backend/vite.config.js`

Replace the entire contents with:

```javascript
import { defineConfig } from 'vite'
import { resolve } from 'path'
import { readdirSync } from 'fs'

/**
 * Recursively find all files with given extension in directory
 */
function findEntryFiles(dir, basePath = '', extension = '.css') {
  const entries = {}
  
  try {
    const items = readdirSync(dir, { withFileTypes: true })
    
    items.forEach(item => {
      const fullPath = resolve(dir, item.name)
      const relativePath = basePath ? `${basePath}/${item.name}` : item.name
      
      if (item.isDirectory()) {
        Object.assign(entries, findEntryFiles(fullPath, relativePath, extension))
      } else if (item.name.endsWith(extension)) {
        // Remove extension from entry name
        const name = relativePath.replace(extension, '')
        entries[name] = fullPath
      }
    })
  } catch (err) {
    // Directory doesn't exist yet, that's okay
  }
  
  return entries
}

// Auto-discover all CSS and JS entry files
const cssEntries = findEntryFiles(resolve(__dirname, 'assets/css'), '', '.css')
const jsEntries = findEntryFiles(resolve(__dirname, 'assets/js'), '', '.js')

export default defineConfig({
  build: {
    outDir: 'static/dist',
    emptyOutDir: true,
    rollupOptions: {
      input: { ...cssEntries, ...jsEntries },
      output: {
        // JS files go to their paths
        entryFileNames: '[name].js',
        // CSS files keep their paths
        assetFileNames: '[name].css'
      }
    },
    minify: true,
    cssMinify: true,
  }
})
```

**What this does:**
- Auto-discovers all `.css` and `.js` files in `assets/` directory
- Preserves directory structure in output (e.g., `assets/css/pages/login.css` → `dist/pages/login.css`)
- Enables minification for both CSS and JS
- Outputs to `static/dist/` so it can be embedded in the Go binary

### Step 1.2: Update package.json

**File:** `backend/package.json`

Replace with:

```json
{
  "name": "announcable-backend",
  "version": "1.0.0",
  "scripts": {
    "dev": "vite build --watch",
    "build": "vite build"
  },
  "devDependencies": {
    "vite": "^7.2.2"
  }
}
```

**What changed:**
- Removed Tailwind dependencies (`@tailwindcss/vite`, `tailwindcss`)
- Removed Basecoat CSS
- Removed Alpine.js, HTMX (these come from CDN)
- Removed other unused packages
- Kept only Vite for bundling

**Run this:**
```bash
cd backend
rm -rf node_modules package-lock.json
npm install
```

### Step 1.3: Delete Old Tailwind Assets

**Delete these files:**
```bash
cd backend
rm -f assets/css/main.css
rm -f assets/js/main.js
rm -f tailwind.config.js
rm -f postcss.config.js
```

These were for the Tailwind integration you're abandoning.

---

## Phase 2: Create CSS Entry Files

### Step 2.1: Create Directory Structure

```bash
cd backend
mkdir -p assets/css/layouts
mkdir -p assets/css/pages
```

### Step 2.2: Create Layout Entry Files

**File:** `backend/assets/css/layouts/appframe.css`

```css
/* Base styles - loaded on every page */
@import '../../../static/css/base/reset.css';
@import '../../../static/css/base/variables.css';

/* Appframe layout */
@import '../../../static/css/layouts/appframe.css';

/* Layout-specific components */
@import '../../../static/css/components/nav.css';
@import '../../../static/css/components/header.css';
@import '../../../static/css/components/alert.css';
```

**File:** `backend/assets/css/layouts/onboard.css`

```css
/* Base styles */
@import '../../../static/css/base/reset.css';
@import '../../../static/css/base/variables.css';

/* Onboard layout */
@import '../../../static/css/layouts/onboard.css';
```

**File:** `backend/assets/css/layouts/fullscreen.css`

```css
/* Base styles */
@import '../../../static/css/base/reset.css';
@import '../../../static/css/base/variables.css';

/* Fullscreen message layout */
@import '../../../static/css/layouts/fullscreenmessage.css';
```

### Step 2.3: Create Page Entry Files

For each page template, create a corresponding CSS entry file that imports all the CSS it needs.

**File:** `backend/assets/css/pages/login.css`

```css
/* Components used by login page */
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';

/* Page-specific styles */
@import '../../../static/css/pages/login.css';
```

**File:** `backend/assets/css/pages/register.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/register.css';
```

**File:** `backend/assets/css/pages/forgot-pw.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/forgot-pw.css';
```

**File:** `backend/assets/css/pages/verify-email.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/modal.css';
@import '../../../static/css/pages/verify-email.css';
```

**File:** `backend/assets/css/pages/release-notes-list.css`

```css
@import '../../../static/css/components/button.css';
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/table.css';
@import '../../../static/css/components/badge.css';
@import '../../../static/css/pages/release-notes-list.css';
```

**File:** `backend/assets/css/pages/release-note-create-edit.css`

```css
@import '../../../static/css/components/button.css';
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/file-input.css';
@import '../../../static/css/components/popover.css';
@import '../../../static/css/components/modal.css';
@import '../../../static/css/pages/release-note-create-edit.css';
```

**File:** `backend/assets/css/pages/settings.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/components/file-input.css';
@import '../../../static/css/pages/settings.css';
```

**File:** `backend/assets/css/pages/user-list.css`

```css
@import '../../../static/css/components/button.css';
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/table.css';
@import '../../../static/css/components/menu.css';
@import '../../../static/css/pages/user-list.css';
```

**File:** `backend/assets/css/pages/landing-page-config.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/components/file-input.css';
@import '../../../static/css/pages/landing-page-config.css';
```

**File:** `backend/assets/css/pages/widget.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/widget.css';
```

**File:** `backend/assets/css/pages/subscribe.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/subscribe.css';
```

**File:** `backend/assets/css/pages/subscription-success.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/subscription-success.css';
```

**File:** `backend/assets/css/pages/subscription-confirm.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/button.css';
```

**File:** `backend/assets/css/pages/pw-reset-token-invalid.css`

```css
@import '../../../static/css/components/card.css';
```

**File:** `backend/assets/css/pages/reset-pw.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
```

**File:** `backend/assets/css/pages/invite-accept.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/form.css';
@import '../../../static/css/components/button.css';
```

**File:** `backend/assets/css/pages/invite-invalid.css`

```css
@import '../../../static/css/components/card.css';
```

**File:** `backend/assets/css/pages/not-found.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/button.css';
```

**File:** `backend/assets/css/pages/admin-dashboard.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/table.css';
```

**File:** `backend/assets/css/pages/admin-org-details.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/table.css';
@import '../../../static/css/components/form.css';
```

**File:** `backend/assets/css/pages/release-notes-website.css`

```css
@import '../../../static/css/components/card.css';
@import '../../../static/css/components/button.css';
@import '../../../static/css/pages/release-notes-website.css';
```

### Step 2.4: Test CSS Build

```bash
cd backend
npm run build
```

**Expected output:**
- `static/dist/layouts/appframe.css`
- `static/dist/layouts/onboard.css`
- `static/dist/layouts/fullscreen.css`
- `static/dist/pages/*.css` (one per page)

Open one of the files - it should be minified CSS with all imports combined.

---

## Phase 3: Create JS Entry Files (Optional)

If you want to bundle JavaScript too, create entry files similar to CSS.

### Step 3.1: Create Directory Structure

```bash
cd backend
mkdir -p assets/js/pages
```

### Step 3.2: Create App-Level JS Entry

**File:** `backend/assets/js/app.js`

```javascript
// Core app utilities loaded on every page
import '../../static/js/app/confirmDialog.js';
import '../../static/js/app/successMsg.js';
import '../../static/js/components/toast.js';
```

### Step 3.3: Create Page-Specific JS Entries

**File:** `backend/assets/js/pages/release-note-create-edit.js`

```javascript
import '../../../static/js/pages/release-note-create-edit.js';
import '../../../static/js/components/file-input.js';
import '../../../static/js/components/popover.js';
```

**File:** `backend/assets/js/pages/settings.js`

```javascript
import '../../../static/js/pages/settings.js';
```

**File:** `backend/assets/js/pages/landing-page-config.js`

```javascript
import '../../../static/js/pages/landing-page-config.js';
import '../../../static/js/components/file-input.js';
```

**File:** `backend/assets/js/pages/widget.js`

```javascript
import '../../../static/js/pages/widget.js';
```

### Step 3.4: Test JS Build

```bash
cd backend
npm run build
```

**Expected output:**
- `static/dist/app.js`
- `static/dist/pages/*.js` (for pages with specific JS)

---

## Phase 4: Update Templates

### Step 4.1: Update Root Layout

**File:** `backend/templates/layouts/root.html`

Find these lines:
```html
<link rel="stylesheet" href="/static/css/base/reset.css" />
<link rel="stylesheet" href="/static/css/base/variables.css" />
{{ block "layout-css" . }}{{ end }}
{{ block "page-css" . }}{{ end }}
<script src="/static/js/app/confirmDialog.js"></script>
<script src="/static/js/app/successMsg.js"></script>
```

Replace with:
```html
<!-- Layout CSS bundle -->
{{ block "layout-css" . }}{{ end }}

<!-- Page CSS bundle -->
{{ block "page-css" . }}{{ end }}
```

And at the bottom, find:
```html
<script src="/static/js/components/toast.js"></script>
```

Replace with:
```html
<!-- App JS bundle -->
<script src="/static/dist/app.js"></script>
```

**Full updated root.html:**

```html
{{ define "root" }}
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    
    <!-- Favicons -->
    <link rel="apple-touch-icon" sizes="180x180" href="/static/media/favicon/apple-touch-icon.png" />
    <link rel="icon" type="image/png" sizes="32x32" href="/static/media/favicon/favicon-32x32.png" />
    <link rel="icon" type="image/png" sizes="16x16" href="/static/media/favicon/favicon-16x16.png" />
    <link rel="manifest" href="/static/media/favicon/site.webmanifest" />
    
    <!-- Layout CSS bundle (one file) -->
    {{ block "layout-css" . }}{{ end }}
    
    <!-- Page CSS bundle (one file) -->
    {{ block "page-css" . }}{{ end }}
    
    <!-- External dependencies -->
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
    <script src="https://unpkg.com/htmx.org@2.0.4" integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/toastify-js/src/toastify.min.css" />
    <script src="https://unpkg.com/sweetalert/dist/sweetalert.min.js"></script>
    <script src="https://unpkg.com/feather-icons"></script>
    <script defer data-domain="app.announcable.me" src="https://plausible.danielbenner.de/js/script.js"></script>

    <title>{{ .Title }} | Announcable</title>
  </head>
  <body>
    {{ block "body" . }}{{ end }}
    
    <!-- App JS bundle -->
    <script src="/static/dist/app.js"></script>
    
    <!-- Feather Icons -->
    <script>feather.replace();</script>
    
    <!-- Toastify -->
    <script src="https://cdn.jsdelivr.net/npm/toastify-js"></script>
    
    <!-- Layout-specific JS -->
    {{ block "layout-js" . }}{{ end }}
    
    <!-- Page-specific JS -->
    {{ block "page-js" . }}{{ end }}
  </body>
</html>
{{ end }}
```

### Step 4.2: Update Layout Templates

**File:** `backend/templates/layouts/appframe.html`

Replace the `layout-css` block:

```html
{{ define "layout-css" }}
<link rel="stylesheet" href="/static/dist/layouts/appframe.css" />
{{ end }}

{{ define "layout-js" }}{{ end }}

{{ define "body" }}
<div class="app">
  <div class="app__nav">{{ template "nav" . }}</div>
  <main class="app__main">
    {{ template "header" . }}
    <div class="app__content">{{ template "main" . }}</div>
  </main>
  {{ block "right-panel" . }}{{ end }}
</div>
{{ end }}
```

**File:** `backend/templates/layouts/onboard.html`

```html
{{ define "layout-css" }}
<link rel="stylesheet" href="/static/dist/layouts/onboard.css" />
{{ end }}

{{ define "layout-js" }}{{ end }}

{{ define "body" }}
<div class="onboard">
  {{ template "main" . }}
</div>
{{ end }}
```

**File:** `backend/templates/layouts/fullscreenmessage.html`

```html
{{ define "layout-css" }}
<link rel="stylesheet" href="/static/dist/layouts/fullscreen.css" />
{{ end }}

{{ define "layout-js" }}{{ end }}

{{ define "body" }}
<div class="fullscreen">
  {{ template "main" . }}
</div>
{{ end }}
```

### Step 4.3: Update Page Templates

For EACH page template, replace the `page-css` block with a single link to the bundled CSS.

**Pattern:**

Before:
```html
{{ define "page-css" }}
<link rel="stylesheet" href="/static/css/components/card.css" />
<link rel="stylesheet" href="/static/css/components/form.css" />
<link rel="stylesheet" href="/static/css/components/button.css" />
<link rel="stylesheet" href="/static/css/pages/login.css" />
{{ end }}
```

After:
```html
{{ define "page-css" }}
<link rel="stylesheet" href="/static/dist/pages/login.css" />
{{ end }}
```

**Apply this to all page templates:**

1. `templates/pages/login.html` → `/static/dist/pages/login.css`
2. `templates/pages/register.html` → `/static/dist/pages/register.css`
3. `templates/pages/forgot-pw.html` → `/static/dist/pages/forgot-pw.css`
4. `templates/pages/verify-email.html` → `/static/dist/pages/verify-email.css`
5. `templates/pages/release-notes-list.html` → `/static/dist/pages/release-notes-list.css`
6. `templates/pages/release-note-create-edit.html` → `/static/dist/pages/release-note-create-edit.css`
7. `templates/pages/settings-page.html` → `/static/dist/pages/settings.css`
8. `templates/pages/user-list.html` → `/static/dist/pages/user-list.css`
9. `templates/pages/landing-page-config.html` → `/static/dist/pages/landing-page-config.css`
10. `templates/pages/widget.html` → `/static/dist/pages/widget.css`
11. `templates/pages/subscribe.html` → `/static/dist/pages/subscribe.css`
12. `templates/pages/subscription-success.html` → `/static/dist/pages/subscription-success.css`
13. `templates/pages/subscription-confirm.html` → `/static/dist/pages/subscription-confirm.css`
14. `templates/pages/reset-pw.html` → `/static/dist/pages/reset-pw.css`
15. `templates/pages/pw-reset-token-invalid.html` → `/static/dist/pages/pw-reset-token-invalid.css`
16. `templates/pages/invite-accept.html` → `/static/dist/pages/invite-accept.css`
17. `templates/pages/invite-invalid.html` → `/static/dist/pages/invite-invalid.css`
18. `templates/pages/not-found.html` → `/static/dist/pages/not-found.css`
19. `templates/pages/admin-dashboard.html` → `/static/dist/pages/admin-dashboard.css`
20. `templates/pages/admin-org-details.html` → `/static/dist/pages/admin-org-details.css`
21. `templates/pages/release-notes-website.html` → `/static/dist/pages/release-notes-website.css`

### Step 4.4: Update Page JS References

For pages that have `page-js` blocks with multiple script tags, replace with bundled JS.

**Example:** `templates/pages/release-note-create-edit.html`

Before:
```html
{{ define "page-js" }}
<script src="/static/js/pages/release-note-create-edit.js"></script>
<script src="/static/js/components/file-input.js"></script>
<script src="/static/js/components/popover.js"></script>
{{ end }}
```

After:
```html
{{ define "page-js" }}
<script src="/static/dist/pages/release-note-create-edit.js"></script>
{{ end }}
```

**Apply to:**
- `templates/pages/release-note-create-edit.html`
- `templates/pages/settings-page.html`
- `templates/pages/landing-page-config.html`
- `templates/pages/widget.html`
- `templates/pages/subscribe.html`

For pages with NO JavaScript (empty `page-js` block), leave them empty.

---

## Phase 5: Update Static Embed

### Step 5.1: Update static.go

**File:** `backend/static/static.go`

Update the embed directive to include the `dist` folder:

```go
package static

import "embed"

//go:embed css/**/* js/**/* media/* dist/*
var Assets embed.FS

//go:embed widget/widget.js
var Widget []byte
```

The `dist/*` pattern tells Go to embed everything in the `static/dist/` directory.

---

## Phase 6: Update .gitignore

### Step 6.1: Ignore Generated Files

**File:** `backend/.gitignore`

Add these lines:

```
# Vite build output
static/dist/

# Dependencies
node_modules/
package-lock.json

# Tailwind/Basecoat (if leftover)
static/dist/
```

---

## Phase 7: Test Everything

### Step 7.1: Build Assets

```bash
cd backend
npm run build
```

Verify output in `static/dist/`:
- `layouts/appframe.css` exists and is minified
- `layouts/onboard.css` exists and is minified
- `layouts/fullscreen.css` exists and is minified
- `pages/*.css` exists (one per page)
- `app.js` exists and is minified
- `pages/*.js` exists (for pages with JS)

### Step 7.2: Build and Run Go App

```bash
cd backend
go build -o tmp/main .
./tmp/main
```

Or use Air:
```bash
make dev-air
```

### Step 7.3: Test Each Page Type

Visit and verify styling:

**Authenticated Pages (appframe layout):**
- http://localhost:8080/release-notes (list)
- http://localhost:8080/release-notes/new (create)
- http://localhost:8080/settings
- http://localhost:8080/users

**Auth Pages (onboard layout):**
- http://localhost:8080/login
- http://localhost:8080/register
- http://localhost:8080/forgot-pw

**Fullscreen Pages:**
- http://localhost:8080/not-found (trigger a 404)

### Step 7.4: Verify in Browser Dev Tools

Open Network tab and check:
- ✅ Only 1-2 CSS files load per page (layout + page)
- ✅ CSS files are minified (view source)
- ✅ Total CSS size is ~80-100KB per page
- ✅ No 404s for missing CSS files
- ✅ JS files load correctly
- ✅ All interactive elements work (forms, buttons, HTMX, Alpine)

### Step 7.5: Test Interactive Features

- ✅ Forms submit correctly
- ✅ HTMX partial updates work
- ✅ Alpine.js interactions work
- ✅ Toast notifications appear
- ✅ Modals open/close
- ✅ File uploads work
- ✅ Popovers work

---

## Phase 8: Update Development Workflow

### Step 8.1: Update Makefile Documentation

**File:** `backend/Makefile`

Update the comments at the top:

```makefile
# Development workflow:
# Terminal 1: make dev-services (start postgres, mail, minio, etc.)
# Terminal 2: npm run dev (Vite watch mode - rebuilds CSS/JS on changes)
# Terminal 3: make dev-air (Go hot-reload with env loaded)
```

### Step 8.2: Daily Development Workflow

```bash
# Start Docker services
cd backend
make dev-services

# Start Vite watch mode (rebuilds CSS/JS automatically)
cd backend
npm run dev

# Start Go app with Air (rebuilds on Go changes)
cd backend
make dev-air
```

Now when you edit any CSS file in `static/css/`, Vite will detect the change through the imports in `assets/css/` and rebuild automatically.

---

## Phase 9: Update Documentation

### Step 9.1: Update CLAUDE.md

Add a new section about asset bundling:

```markdown
### Asset Bundling

The application uses Vite to bundle CSS and JavaScript:

**Structure:**
- `backend/static/css/` - Source CSS files (components, layouts, pages)
- `backend/static/js/` - Source JavaScript files
- `backend/assets/css/` - CSS entry points that import from static/css
- `backend/assets/js/` - JS entry points that import from static/js
- `backend/static/dist/` - Built and minified output (git-ignored)

**Adding CSS to a Page:**
1. Edit the CSS entry file: `assets/css/pages/my-page.css`
2. Add `@import` statements for components you need
3. Vite will automatically rebuild in watch mode
4. The template loads one bundled CSS file

**Build Commands:**
- `npm run dev` - Watch mode (rebuilds on file changes)
- `npm run build` - Production build (minified)

**Templates:**
Each template loads exactly ONE CSS file from `/static/dist/pages/[page-name].css`
Layouts load ONE CSS file from `/static/dist/layouts/[layout-name].css`
```

---

## Troubleshooting

### Problem: Vite build fails with "cannot find module"

**Solution:** Check that all `@import` paths in your entry files are correct. They should be relative paths from the entry file to the source CSS file.

Example: If you're in `assets/css/pages/login.css` and want to import `static/css/components/button.css`, the path is:
```css
@import '../../../static/css/components/button.css';
```

### Problem: CSS doesn't load in the browser (404)

**Solution:** 
1. Check that you ran `npm run build`
2. Check that `static/dist/` contains the expected files
3. Verify the embed directive in `static/static.go` includes `dist/*`
4. Rebuild the Go binary: `go build`

### Problem: Styles don't look right

**Solution:**
1. Check browser console for CSS errors
2. Verify all required components are imported in the entry file
3. Check that you're not missing any components that were previously in the HTML
4. Use browser dev tools to inspect which CSS file is loaded

### Problem: Changes to CSS not reflected

**Solution:**
1. Make sure Vite watch mode is running (`npm run dev`)
2. Refresh the browser (the Go app needs to be restarted to pick up new embedded assets)
3. Or use Air which will restart on file changes

### Problem: Duplicate CSS rules

**Solution:** Vite automatically deduplicates CSS. If you see duplicates, it's because you imported the same file twice in one entry file.

---

## Rollback Plan

If something goes wrong and you need to rollback:

1. **Revert template changes** - Use git to restore HTML templates
2. **Remove dist folder** - `rm -rf backend/static/dist`
3. **Revert static.go** - Remove `dist/*` from embed directive
4. **Rebuild Go app** - `go build`

The source CSS files in `static/css/` are unchanged, so you can always go back to loading them individually.

---

## Success Criteria

✅ All pages load and look correct  
✅ All interactive features work (forms, modals, HTMX, Alpine)  
✅ Browser Network tab shows 1-2 CSS files per page  
✅ CSS files are minified  
✅ Total CSS size reduced by ~60%  
✅ Vite watch mode rebuilds automatically  
✅ Production build works  

---

## Next Steps After Implementation

1. **Monitor Performance** - Use browser dev tools to verify faster page loads
2. **Update Documentation** - Add bundling info to CLAUDE.md
3. **Add Cache Headers** - Configure Go HTTP server to send proper cache headers for `/static/dist/*`
4. **Consider Content Hashing** - In the future, you might want to add content hashes to filenames for better caching (e.g., `login-abc123.css`)

