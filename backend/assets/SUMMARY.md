# Assets (Source Files)

This directory contains source CSS and JS files that are built by Vite.

**DO NOT edit files in `../static/dist/` - they are auto-generated!**

## Build Process

```
Edit source → Vite builds → Go embeds → Users receive
(assets/)     (→ static/dist/)  (static.go)    (HTTP)
```

## Commands

- **Development**: `npm run dev` (watch mode with hot-reload)
- **Production build**: `npm run build` → outputs to `../static/dist/`

## Architecture

- **`assets/`** = Source files (you edit these, under version control)
- **`static/dist/`** = Built output (auto-generated, should be in .gitignore)
- **`static/media/`** = True static assets (brand logos, favicons)
- **`static/widget/`** = Embeddable React widget (built separately from `/widget` project)

For detailed CSS architecture, see `css/AGENTS.md`  
For detailed JS architecture, see `js/AGENTS.md`
