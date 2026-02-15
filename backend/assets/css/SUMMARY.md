# CSS Architecture for AI Assistants

## Structure

```
assets/css/
  base/          # CSS reset + design tokens (variables)
  components/    # Reusable UI components
  layouts/       # Layout entry files (inline styles + imports)
  pages/         # Page entry files (inline styles + imports)
```

**Build output:** `static/dist/` (mirrored structure, bundled & minified by Vite)

## How It Works

1. **Entry files** (`layouts/` and `pages/`) contain:

   - All `@import` statements at the top
   - Inline page/layout-specific CSS below imports
   - Imports use relative paths: `../base/`, `../components/`

2. **Vite** bundles entry files → `static/dist/`

3. **Templates** load bundled CSS:
   - Layout CSS: `/static/dist/layouts/{layout}.css`
   - Page CSS: `/static/dist/pages/{page}.css`

## Adding CSS to a Page

**Example: Adding modal to login page**

Edit `assets/css/pages/login.css`:

```css
/* All @import statements must come first */
@import '../components/card.css';
@import '../components/form.css';
@import '../components/button.css';
@import '../components/modal.css';  /* ADD THIS */

/* Login page styles */
.register-hint {
  text-align: center;
  font-size: var(--font-size-sm);
}
```

Run `npm run build` → outputs to `static/dist/pages/login.css`

## Creating a New Page

1. Create entry file: `assets/css/pages/my-page.css`
2. Add imports + page-specific styles
3. Template uses: `/static/dist/pages/my-page.css`
4. Vite auto-discovers and builds it

## Creating a New Component

1. Create: `assets/css/components/my-component.css`
2. Import in entry files that need it
3. Vite bundles it into consuming pages

## Rules

- **ALWAYS** put `@import` statements at the top (CSS requirement)
- **DON'T** import from `static/` - everything is in `assets/`
- **DON'T** edit files in `static/dist/` - they're auto-generated
- **DO** use relative imports: `../components/`, not `../../static/`
- **DO** run `npm run build` after changes
- **DO** rebuild Go app (`go build`) to embed new assets

## Development Workflow

```bash
# Terminal 1: Vite watch mode (auto-rebuild CSS)
npm run dev

# Terminal 2: Go app with Air (auto-restart on changes)
make dev-air
```

## Import Order Pattern

```css
/* Entry file template */

/* 1. Base (if needed - usually in layouts only) */
@import '../base/reset.css';
@import '../base/variables.css';

/* 2. Components (alphabetical) */
@import '../components/button.css';
@import '../components/card.css';
@import '../components/form.css';

/* 3. Page/layout-specific styles below */
.my-page {
  /* styles */
}
```

## Debugging

- **CSS not loading?** Check `static/dist/` - file should exist
- **Wrong styles?** Check imports in entry file
- **404 on CSS?** Rebuild Go app to embed new dist files
- **Changes not reflected?** Run `npm run build` first

## File Mapping

```
Template                    → Entry File                    → Bundle
templates/pages/login.html  → assets/css/pages/login.css   → static/dist/pages/login.css
templates/layouts/root.html → (none - loaded via blocks)   → (loaded from page/layout bundles)
```
