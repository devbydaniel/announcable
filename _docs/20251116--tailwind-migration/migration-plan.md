# Tailwind CSS Migration Plan

**Date**: 2025-11-16
**Status**: Planning
**Estimated Effort**: 2-3 days

## Executive Summary

This document outlines the migration from the current manual CSS architecture to a Tailwind CSS-based approach. The current system requires manually managing CSS file references in every template, leading to maintenance overhead and multiple HTTP requests. Tailwind will provide a single bundled CSS file, utility-first styling, and automatic purging of unused styles.

## Current State Analysis

### Problems with Current Approach

1. **Manual CSS Management**: Each page template requires explicit CSS file references via `define "page-css"` blocks
2. **Multiple HTTP Requests**: 15+ separate CSS files loaded across pages (base, components, layouts, pages)
3. **No Build Process**: Hand-written CSS with no minification, autoprefixing, or optimization
4. **Maintenance Burden**: Easy to forget required CSS files when creating new pages
5. **No Dead Code Elimination**: Unused CSS accumulates over time
6. **Inconsistent Spacing/Sizing**: Manual CSS leads to arbitrary values (padding: 14px vs 16px)

### Current CSS Structure

```
backend/static/css/
├── base/
│   ├── reset.css           # CSS reset
│   └── variables.css       # CSS custom properties
├── components/
│   ├── alert.css
│   ├── badge.css
│   ├── button.css
│   ├── card.css
│   ├── checkbox.css
│   ├── file-input.css
│   ├── form.css
│   ├── header.css
│   ├── menu.css
│   ├── modal.css
│   ├── nav.css
│   ├── popover.css
│   ├── skeleton.css
│   └── table.css
├── layouts/
│   ├── appframe.css
│   ├── fullscreenmessage.css
│   └── onboard.css
└── pages/
    ├── forgot-pw.css
    ├── landing-page-config.css
    ├── login.css
    ├── register.css
    ├── release-note-create-edit.css
    ├── release-notes-list.css
    ├── release-notes-website.css
    ├── settings.css
    ├── subscribe.css
    ├── subscription-success.css
    ├── user-list.css
    ├── verify-email.css
    └── widget.css
```

### Current Template Pattern

```html
{{ define "page-css" }}
<link rel="stylesheet" href="/static/css/components/card.css" />
<link rel="stylesheet" href="/static/css/components/form.css" />
<link rel="stylesheet" href="/static/css/components/button.css" />
<link rel="stylesheet" href="/static/css/pages/login.css" />
{{ end }}
```

## Target State

### Goals

1. **Single CSS Bundle**: One optimized CSS file served to all pages
2. **Utility-First Styling**: Tailwind utilities directly in HTML templates
3. **Automatic Optimization**: Purge unused CSS, minify for production
4. **Design System Enforcement**: Consistent spacing, colors, typography via Tailwind config
5. **Developer Experience**: No manual CSS file management
6. **Backward Compatibility**: Maintain existing functionality during migration

### Target CSS Structure

```
backend/static/css/
├── input.css               # Tailwind entry point
├── output.css              # Generated bundle (dev)
├── output.min.css          # Generated bundle (prod)
└── custom/                 # Custom CSS for edge cases
    └── legacy.css          # Temporary: styles not yet migrated
```

### Target Template Pattern

```html
{{ define "root" }}
<!doctype html>
<html lang="en">
  <head>
    ...
    <link rel="stylesheet" href="/static/css/output.css" />
    {{ block "page-css" . }}{{ end }}
    <!-- For legacy/custom CSS only -->
    ...
  </head>
  ... {{ end }}
</html>
```

No more per-page CSS blocks needed for standard pages.

## Migration Strategy

### Phase 1: Setup Tailwind Build Process

**Goal**: Get Tailwind building alongside existing CSS (no breaking changes)

**Steps**:

1. **Install Tailwind in backend**:

   ```bash
   cd backend
   npm init -y  # Create package.json
   npm install -D tailwindcss@latest postcss@latest autoprefixer@latest
   npm install -D @tailwindcss/forms @tailwindcss/typography
   ```

2. **Initialize Tailwind config**:

   ```bash
   npx tailwindcss init -p
   ```

3. **Configure `tailwind.config.js`**:

   ```js
   /** @type {import('tailwindcss').Config} */
   module.exports = {
     content: ["./templates/**/*.html", "./internal/handler/**/*.go"],
     theme: {
       extend: {
         // Map existing CSS variables to Tailwind
         colors: {
           // Extract from variables.css
         },
         spacing: {
           // Define consistent spacing scale
         },
       },
     },
     plugins: [
       require("@tailwindcss/forms"),
       require("@tailwindcss/typography"),
     ],
   };
   ```

4. **Create `backend/static/css/input.css`**:

   ```css
   @tailwind base;
   @tailwind components;
   @tailwind utilities;

   /* Custom component classes for complex components */
   @layer components {
     /* Migration: keep essential custom components here */
   }
   ```

5. **Add build scripts to `backend/package.json`**:

   ```json
   {
     "scripts": {
       "css:dev": "tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch",
       "css:build": "tailwindcss -i ./static/css/input.css -o ./static/css/output.min.css --minify",
       "format": "prettier --write 'templates/**/*.html'"
     },
     "devDependencies": {
       "tailwindcss": "^3.4.0",
       "postcss": "^8.4.0",
       "autoprefixer": "^10.4.0",
       "@tailwindcss/forms": "^0.5.0",
       "@tailwindcss/typography": "^0.5.0",
       "prettier": "^3.4.0",
       "prettier-plugin-go-template": "^0.0.15"
     }
   }
   ```

6. **Update `.air.toml` to run CSS build**:

   ```toml
   [build]
     cmd = "cd .. && npm run --prefix backend css:build && cd backend && go generate ./... && go build -o ./tmp/main ."
   ```

7. **Update `root.html` to include Tailwind**:

   ```html
   <link rel="stylesheet" href="/static/css/output.css" />
   <link rel="stylesheet" href="/static/css/base/reset.css" />
   <link rel="stylesheet" href="/static/css/base/variables.css" />
   ```

8. **Update `.gitignore`**:
   ```
   backend/node_modules/
   backend/static/css/output.css
   backend/static/css/output.min.css
   backend/package-lock.json
   ```

**Validation**: Run `npm run css:dev` and verify `output.css` is generated without errors.

### Phase 2: Create Tailwind Component Mappings

**Goal**: Document Tailwind equivalents for existing CSS components

**Steps**:

1. **Audit existing CSS components**: Review each file in `static/css/components/`

2. **Create mapping document**: `_docs/20251116--tailwind-migration/component-mappings.md`

3. **For each component, document**:
   - Original CSS class structure
   - Tailwind equivalent utility classes
   - Any custom CSS still needed

**Example mapping**:

````markdown
## Button Component

### Original (`components/button.css`)

```css
.button {
  padding: 0.75rem 1.5rem;
  border-radius: 0.375rem;
  font-weight: 600;
}

.button--primary {
  background: var(--color-primary);
  color: white;
}
```
````

### Tailwind Equivalent

```html
<button
  class="px-6 py-3 rounded-md font-semibold bg-blue-600 text-white hover:bg-blue-700"
></button>
```

### Custom CSS Needed

None - fully replaced by Tailwind utilities

````

4. **Identify components requiring custom CSS**:
   - Complex components (nav, modal, popover)
   - Components with animations
   - Components with specific business logic styling

**Deliverable**: Complete mapping document for all components

### Phase 3: Migrate Templates Page-by-Page

**Goal**: Convert templates from custom CSS to Tailwind utilities

**Strategy**: Start with simplest pages first (auth flows), then tackle complex pages

**Migration Order**:

1. **Simple pages** (no complex layouts):
   - login.html
   - register.html
   - forgot-pw.html
   - verify-email.html
   - subscription-confirm.html

2. **Medium complexity** (forms with validation):
   - settings-page.html
   - invite-accept.html
   - user-list.html

3. **Complex pages** (heavy interactivity):
   - release-note-create-edit.html
   - release-notes-list.html
   - landing-page-config.html
   - widget.html
   - admin-org-details.html

**Per-Page Migration Steps**:

1. **Open the page template** (e.g., `templates/pages/login.html`)

2. **Remove `page-css` block** (or comment out initially)

3. **Replace CSS classes with Tailwind utilities**:
   - `.card` → `bg-white shadow-md rounded-lg p-6`
   - `.form__group` → `mb-4`
   - `.form__label` → `block text-sm font-medium text-gray-700 mb-1`
   - `.form__input` → `w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500`
   - `.button.button--primary` → `w-full bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500`

4. **Test the page** in browser for visual regressions

5. **Commit the change**:
   ```bash
   git add templates/pages/login.html
   git commit -m "feat: migrate login page to Tailwind CSS"
````

6. **Once page is fully migrated**, delete corresponding CSS files:
   ```bash
   rm static/css/pages/login.css
   git commit -am "chore: remove legacy login.css"
   ```

**Example: Login Page Migration**

**Before**:

```html
{{ define "page-css" }}
<link rel="stylesheet" href="/static/css/components/card.css" />
<link rel="stylesheet" href="/static/css/components/form.css" />
<link rel="stylesheet" href="/static/css/components/button.css" />
<link rel="stylesheet" href="/static/css/pages/login.css" />
{{ end }} {{ define "main" }}
<div class="card" x-data>
  <h2 class="card__title">Login</h2>
  <form class="form" hx-post="/login">
    <div class="form__group">
      <label class="form__label" for="email">Email</label>
      <input
        class="form__input"
        type="email"
        id="email"
        name="email"
        required
        autofocus
      />
    </div>
    <div class="form__group">
      <label class="form__label" for="password">Password</label>
      <input
        class="form__input"
        type="password"
        id="password"
        name="password"
        required
      />
    </div>
    <div class="form__group">
      <button class="button button--block button--primary" type="submit">
        Login
      </button>
    </div>
  </form>
</div>
{{ end }}
```

**After**:

```html
{{ define "page-css" }}{{ end }} {{ define "main" }}
<div class="bg-white shadow-md rounded-lg p-6 max-w-md mx-auto" x-data>
  <h2 class="text-2xl font-semibold mb-6">Login</h2>
  <form class="space-y-4" hx-post="/login">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1" for="email"
        >Email</label
      >
      <input
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        type="email"
        id="email"
        name="email"
        required
        autofocus
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1" for="password"
        >Password</label
      >
      <input
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        type="password"
        id="password"
        name="password"
        required
      />
      <p class="text-sm text-gray-600 mt-1">
        <a href="/forgot-pw" class="text-blue-600 hover:underline"
          >Forgot password?</a
        >
      </p>
    </div>
    <div>
      <button
        class="w-full bg-blue-600 text-white px-4 py-2 rounded-md font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        type="submit"
      >
        Login
      </button>
    </div>
    <p class="text-center text-sm text-gray-600 mt-4">
      <a href="/register" class="text-blue-600 hover:underline"
        >No account yet? Register</a
      >
    </p>
  </form>
</div>
{{ end }}
```

### Phase 4: Migrate Shared Components

**Goal**: Convert reusable partials and layout components

**Components to Migrate**:

1. **`partials/nav.html`**:
   - Complex component with active states
   - May need custom CSS for active link styling
   - Use Tailwind for spacing, colors, transitions

2. **`partials/header.html`**:
   - Simpler component
   - Fully migrate to Tailwind

3. **`layouts/appframe.html`**:
   - Grid/flexbox layout
   - Tailwind works great for this

4. **`layouts/onboard.html`**:
   - Centered layout
   - Straightforward Tailwind migration

**Approach**:

1. Migrate one component at a time
2. Test across all pages that use the component
3. Keep component CSS temporarily if needed for complex styles
4. Document any remaining custom CSS requirements

### Phase 5: Handle Complex Components

**Goal**: Migrate or preserve complex component styles

**Complex Components Requiring Attention**:

1. **Modal/Dialog**: May need custom CSS for animations and positioning
2. **Popover/Menu**: Positioning logic often requires custom CSS
3. **Table**: Complex table styling might need custom classes
4. **Skeleton**: Animation keyframes require custom CSS

**Strategy**:

1. **Use `@layer components`** in `input.css` for complex components:

   ```css
   @layer components {
     .modal {
       @apply fixed inset-0 z-50 flex items-center justify-center;
       /* Custom animation CSS */
     }
   }
   ```

2. **Or keep as separate CSS file** if very complex:
   - Move to `static/css/custom/modal.css`
   - Include in templates that need it

3. **Prefer Headless UI or Radix UI** for complex interactive components:
   - These work well with Tailwind
   - But may be overkill for this project

### Phase 6: Clean Up and Optimize

**Goal**: Remove legacy CSS, optimize build process

**Steps**:

1. **Delete unused CSS files**:

   ```bash
   rm -rf static/css/components/
   rm -rf static/css/layouts/
   rm -rf static/css/pages/
   ```

2. **Remove CSS block definitions from templates**:
   - Remove `{{ define "layout-css" }}` blocks from layouts
   - Remove `{{ define "page-css" }}` blocks from pages
   - Or keep them empty for future custom CSS

3. **Update `root.html`**:

   ```html
   <head>
     ...
     <link rel="stylesheet" href="/static/css/output.css" />
     {{ block "custom-css" . }}{{ end }}
     <!-- Only for truly custom styles -->
     ...
   </head>
   ```

4. **Optimize Tailwind config**:

   ```js
   module.exports = {
     content: ["./templates/**/*.html", "./internal/handler/**/*.go"],
     theme: {
       extend: {
         colors: {
           primary: "#3b82f6", // blue-600
           secondary: "#64748b", // slate-500
           danger: "#ef4444", // red-500
           success: "#10b981", // green-500
         },
       },
     },
     plugins: [require("@tailwindcss/forms")],
   };
   ```

5. **Configure PurgeCSS for production**:
   - Tailwind automatically purges unused styles
   - Ensure `content` paths cover all templates

6. **Update Dockerfile to build CSS**:

   ```dockerfile
   # Install Node.js for Tailwind build
   FROM node:20-alpine AS css-builder
   WORKDIR /app/backend
   COPY backend/package*.json ./
   RUN npm ci
   COPY backend/static/css/input.css ./static/css/
   COPY backend/templates ./templates
   RUN npm run css:build

   # Go build stage
   FROM golang:1.23-alpine AS go-builder
   # ... existing Go build ...
   COPY --from=css-builder /app/backend/static/css/output.min.css ./static/css/
   ```

7. **Add CSS build to CI/CD**:

   ```yaml
   # .github/workflows/deploy.yml
   - name: Build CSS
     run: |
       cd backend
       npm ci
       npm run css:build
   ```

8. **Update CLAUDE.md** with new CSS architecture

### Phase 7: Documentation and Training

**Goal**: Document new workflow for future development

**Deliverables**:

1. **Update CLAUDE.md**:
   - Remove references to manual CSS management
   - Document Tailwind workflow
   - Explain how to add new pages

2. **Create `_docs/20251116--tailwind-migration/tailwind-guide.md`**:
   - Common Tailwind patterns for this project
   - Component composition examples
   - When to use custom CSS vs Tailwind

3. **Document gotchas**:
   - Go template syntax conflicts (if any)
   - HTMX-specific styling needs
   - Alpine.js integration patterns

## Risk Mitigation

### Potential Risks

1. **Visual Regressions**: Styles don't match pixel-perfect
   - **Mitigation**: Migrate one page at a time, test thoroughly
   - **Acceptance**: Minor visual differences are acceptable if design improves

2. **Build Process Complexity**: Adding Node.js dependency to Go project
   - **Mitigation**: Well-documented build scripts, update Dockerfile
   - **Fallback**: Keep generated CSS in git temporarily during transition

3. **Developer Unfamiliarity**: Team doesn't know Tailwind
   - **Mitigation**: Provide documentation and examples
   - **Benefit**: Tailwind is widely used, easy to learn

4. **Template File Size**: Inline utilities make HTML verbose
   - **Mitigation**: Use `@apply` for repeated patterns
   - **Acceptance**: Verbosity is acceptable tradeoff for maintainability

5. **Incomplete Migration**: Some pages left in old CSS
   - **Mitigation**: Track progress in this document
   - **Strategy**: Keep both systems working during transition

### Rollback Plan

If migration needs to be rolled back:

1. **Revert template changes** via git
2. **Re-add CSS file links** to templates
3. **Keep legacy CSS files** in git until fully validated
4. **Remove Tailwind config** and dependencies

## Success Criteria

Migration is complete when:

- [ ] All templates migrated to Tailwind utilities
- [ ] Legacy CSS files deleted (`components/`, `layouts/`, `pages/`)
- [ ] Single CSS bundle served (`output.css` or `output.min.css`)
- [ ] Build process runs in dev (Air) and production (Dockerfile)
- [ ] No visual regressions on key pages
- [ ] Documentation updated (CLAUDE.md, tailwind-guide.md)
- [ ] CSS bundle size < 50KB gzipped (with purging)

## Timeline Estimate

| Phase                              | Effort       | Duration            |
| ---------------------------------- | ------------ | ------------------- |
| Phase 1: Setup                     | 2 hours      | Day 1 AM            |
| Phase 2: Component Mappings        | 3 hours      | Day 1 PM            |
| Phase 3: Migrate Simple Pages      | 4 hours      | Day 1 PM - Day 2 AM |
| Phase 4: Migrate Shared Components | 3 hours      | Day 2 PM            |
| Phase 5: Complex Components        | 4 hours      | Day 2 PM - Day 3 AM |
| Phase 6: Clean Up                  | 2 hours      | Day 3 PM            |
| Phase 7: Documentation             | 2 hours      | Day 3 PM            |
| **Total**                          | **20 hours** | **2.5 days**        |

## Post-Migration Benefits

1. **Single CSS File**: Reduced from 15+ files to 1, faster page loads
2. **Smaller Bundle Size**: Purging removes unused CSS (~50KB vs ~200KB estimated)
3. **Consistent Design**: Tailwind's spacing/color scales enforce consistency
4. **Faster Development**: No more manually managing CSS files per page
5. **Better DX**: IntelliSense for Tailwind classes in VSCode
6. **Future-Proof**: Easier to add new pages, just use utilities

## Appendix

### Alternative Approaches Considered

1. **CSS-in-JS**: Not suitable for Go templates
2. **SASS/SCSS**: Adds build complexity without utility-first benefits
3. **CSS Modules**: Not well-supported in Go ecosystem
4. **Keep Current Approach**: Technical debt continues to grow

### Resources

- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Tailwind with HTMX Examples](https://htmx.org/examples/)
- [Tailwind Forms Plugin](https://github.com/tailwindlabs/tailwindcss-forms)
- [Go Template Syntax](https://pkg.go.dev/html/template)

### Questions for Review

1. Should we keep `variables.css` for CSS custom properties or migrate fully to Tailwind config?
2. Any pages that absolutely cannot be migrated (e.g., widget embed page)?
3. Should we use Prettier for HTML formatting consistency?
4. Keep `page-css` blocks for future extensibility or remove entirely?
