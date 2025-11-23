# JavaScript Architecture Guide

This document explains the JavaScript bundling architecture for AI agents and developers.

## Overview

All JavaScript source files are stored in `assets/js/` and bundled using Vite into `static/dist/`. This mirrors the CSS bundling architecture.

## Directory Structure

```
assets/js/
  app/                 # App-level utilities (loaded on every page)
    confirmDialog.js   # HTMX confirm dialog handler
    successMsg.js      # URL param success message handler
  components/          # Reusable components
    toast.js           # Toast notification functions
    file-input.js      # Alpine.js file input component
    popover.js         # Alpine.js popover component
  pages/               # Page-specific scripts
    release-note-create-edit.js
    settings.js
    landing-page-config.js
    widget.js
  app.js                          # Entry: App-level bundle
  pages/release-note-create-edit-bundle.js  # Entry: Release note page bundle
  pages/settings-bundle.js        # Entry: Settings page bundle
  pages/landing-page-config-bundle.js      # Entry: Landing page config bundle
  pages/widget-bundle.js          # Entry: Widget config bundle
```

## Source Files and Build Process

**Source Files**: Individual JS modules with specific functionality (e.g., `toast.js`, `confirmDialog.js`)

Vite processes all source files and minifies them individually into `static/dist/`, preserving the directory structure.

## File Mapping

| Template | Minified JS Files | Source Files |
|----------|------------------|--------------|
| root.html (all pages) | `/static/dist/app/confirmDialog.js`<br>`/static/dist/app/successMsg.js`<br>`/static/dist/components/toast.js` | confirmDialog.js, successMsg.js, toast.js |
| release-note-create-edit.html | `/static/dist/pages/release-note-create-edit.js`<br>`/static/dist/components/file-input.js`<br>`/static/dist/components/popover.js` | release-note-create-edit.js, file-input.js, popover.js |
| settings-page.html | `/static/dist/pages/settings.js` | settings.js |
| subscribe.html | `/static/dist/pages/settings.js` | settings.js |
| landing-page-config.html | `/static/dist/pages/landing-page-config.js`<br>`/static/dist/components/file-input.js` | landing-page-config.js, file-input.js |
| widget.html | `/static/dist/pages/widget.js` | widget.js |

## How to Add JS to a Page

### 1. Add to Existing Source File

Edit the appropriate source file in `assets/js/`:
- App-level utilities → `assets/js/app/`
- Reusable components → `assets/js/components/`
- Page-specific logic → `assets/js/pages/`

### 2. Create New JS for New Page

1. Create source file: `assets/js/pages/my-page.js`
2. Add your JavaScript code (Alpine.js components, event listeners, etc.)
3. Add to template:
   ```html
   {{ define "page-js" }}
     <script src="/static/dist/pages/my-page.js"></script>
     <!-- Include any required components -->
     <script src="/static/dist/components/some-component.js"></script>
   {{ end }}
   ```
4. Build: `npm run build`

## How to Create a New Component

1. Create `assets/js/components/my-component.js`
2. Add component logic (typically Alpine.js or vanilla JS)
3. Import it in the appropriate entry file
4. Build: `npm run build`

Example component using Alpine.js:
```javascript
document.addEventListener('alpine:init', () => {
  Alpine.data('myComponent', () => ({
    // Component state and methods
    isOpen: false,
    toggle() {
      this.isOpen = !this.isOpen;
    }
  }));
});
```

## Development Workflow

### Development Mode (with hot-reload)
```bash
cd backend
npm run dev
```
This starts Vite in watch mode. Changes to JS files trigger automatic rebuilds.

### Production Build
```bash
cd backend
npm run build
```
This bundles and minifies all JS files to `static/dist/`.

### Full Stack Development
```bash
# Terminal 1: Docker services
make dev-services

# Terminal 2: Go backend with Air
make dev-air

# Terminal 3: Vite for JS/CSS hot-reload
npm run dev
```

## Code Patterns

### Vanilla JavaScript - No Module Syntax
Source files are vanilla JavaScript - no `export` or `import` needed. They define global functions or Alpine.js components:

```javascript
// toast.js - defines global functions
function toastSuccess(message) {
  Toastify({
    text: message,
    className: "toast toast--success",
    duration: 3000
  }).showToast();
}

// file-input.js - defines Alpine.js component
document.addEventListener('alpine:init', () => {
  Alpine.data('fileInput', (initialUrl) => ({
    imgUrl: initialUrl,
    // ... component logic
  }));
});
```

## External Dependencies

All JS code depends on CDN-loaded libraries defined in `root.html`:
- **Alpine.js**: Reactive UI framework
- **HTMX**: AJAX and partial page updates
- **SweetAlert**: Modal dialogs (used by confirmDialog.js)
- **Toastify**: Toast notifications (used by toast.js)
- **Feather Icons**: Icon library

These libraries are available as global objects - no imports needed.

## Debugging Tips

1. **Check browser console** for JS errors
2. **Network tab** shows which bundle files are loaded
3. **Verify bundle exists** in `static/dist/` after build
4. **Check template** references correct bundle path
5. **Rebuild after changes**: `npm run build` or use `npm run dev`
6. **Empty bundle warning** is normal if imports produce no output

## Build Output

After `npm run build`, you'll see:
```
static/dist/
  app/
    confirmDialog.js                 # Minified (~0.24 kB)
    successMsg.js                    # Minified (~0.25 kB)
  components/
    toast.js                         # Minified
    file-input.js                    # Minified (~0.26 kB)
    popover.js                       # Minified (~0.23 kB)
  pages/
    release-note-create-edit.js      # Minified (~0.46 kB)
    settings.js                      # Minified (~1.08 kB)
    landing-page-config.js           # Minified (~2.50 kB)
    widget.js                        # Minified (~4.37 kB)
```

Templates load the minified files directly from `/static/dist/`.

## Key Principles

1. **Source files in assets/js/**: Never edit files in `static/dist/`
2. **Vite handles minification**: Auto-discovers all `.js` files in `assets/js/`
3. **Templates load minified files**: Reference `/static/dist/` paths
4. **Global scope**: Functions and components are globally accessible
5. **App-level JS loads first**: Toast, confirmDialog, successMsg available everywhere
6. **Multiple script tags OK**: Load dependencies in correct order

## Common Tasks

### Add toast notification to a new page
Toast is in app.js bundle (loaded on all pages), so just use:
```javascript
toastSuccess('Operation successful!');
toastError('Something went wrong!');
toastInfo('Helpful information');
```

### Add Alpine.js component to a page
1. Create component in `assets/js/components/my-component.js`
2. Import in page's entry bundle
3. Use in template: `<div x-data="myComponent">`

### Share code between pages
1. Create shared module in `assets/js/components/`
2. Import in multiple entry bundles
3. Vite will include it in each bundle (code duplication is acceptable for simplicity)

## Architecture Consistency

This JS bundling architecture matches the CSS architecture:
- Source files in `assets/`
- Bundled output in `static/dist/`
- Entry files for grouping
- Vite for building
- Progressive enhancement (app-level → page-specific)
