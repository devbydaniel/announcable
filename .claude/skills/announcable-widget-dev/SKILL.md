---
name: announcable-widget-dev
description: Widget development in Announcable. Use when creating, modifying, or debugging the embeddable widget (Lit, TypeScript, Vite, UMD).
---

# Widget Development — Announcable

## Working Directory

**All commands run from `widget/`:**

```bash
cd widget
```

## Validation Sequence

Run after every change. Do NOT trust your own assessment — verify through observable behavior.

```bash
npm run lint                   # Must pass
npm run build                  # Must succeed (UMD bundle)
```

## Project Structure

```
widget/
├── src/
│   ├── main.ts                # Entry point, registers <announcable-widget> custom element
│   ├── app.ts                 # Main widget LitElement component
│   ├── index.css              # Base widget styles
│   ├── components/
│   │   ├── ui/                # Generic UI components
│   │   ├── widget/            # Widget-specific view components
│   │   └── icons/             # SVG icon components
│   ├── controllers/           # Lit reactive controllers
│   │   ├── anchors.ts         # Anchor link behavior
│   │   └── widget-toggle.ts   # Widget open/close toggle
│   ├── tasks/                 # Async data fetching (Lit Task pattern)
│   │   ├── release-notes.ts   # Fetch release notes from API
│   │   ├── release-note-likes.ts
│   │   ├── release-note-metrics.ts
│   │   ├── release-note-status.ts
│   │   └── widget-config.ts   # Fetch widget configuration
│   └── lib/                   # Utilities and shared code
│       ├── base-component.ts  # Base class for all components
│       ├── clientId.ts        # Anonymous client identification
│       ├── config.ts          # Backend URL configuration
│       ├── contexts.ts        # Lit context definitions
│       ├── storage.ts         # Local storage helpers
│       ├── types.ts           # Shared TypeScript types
│       └── utils.ts           # General utilities
├── vite.config.ts             # Builds as UMD library
└── package.json
```

## Architecture

The widget is a self-contained Lit web component application:

1. **Entry** (`main.ts`): Registers `<announcable-widget>` custom element
2. **App** (`app.ts`): Root LitElement that orchestrates data fetching and rendering
3. **Components**: Lit components for UI rendering
4. **Controllers**: Lit reactive controllers for reusable behavior (toggle, anchors)
5. **Tasks**: Lit `Task` pattern for async data fetching from the backend API
6. **Lib**: Shared utilities, types, and configuration

### Build Output

Built as **UMD bundle** (`widget.js`) with CSS injected by JS (no separate CSS file). Embedded on customer websites via:

```html
<script src="https://announcable.com/widget?id=ORG_ID"></script>
```

### API Communication

The widget fetches data from the backend API at `/api/`:

| Endpoint | Task | Purpose |
|----------|------|---------|
| `GET /api/release-notes/{orgId}` | `release-notes.ts` | Fetch published release notes |
| `GET /api/release-notes/{orgId}/status` | `release-note-status.ts` | Check for new notes |
| `POST /api/release-notes/{orgId}/metrics` | `release-note-metrics.ts` | Track views |
| `GET /api/release-notes/{orgId}/{id}/like` | `release-note-likes.ts` | Get like state |
| `POST /api/release-notes/{orgId}/{id}/like` | `release-note-likes.ts` | Toggle like |
| `GET /api/widget-config/{orgId}` | `widget-config.ts` | Get widget appearance config |

## Key Patterns

### Base Component

All components extend `BaseComponent` from `lib/base-component.ts` which provides shared functionality.

### Lit Context

Widget-wide state is shared via Lit contexts defined in `lib/contexts.ts`. The root `app.ts` provides context values consumed by child components.

### Lit Task Pattern

Data fetching uses Lit's `Task` class for declarative async operations:

```typescript
// tasks/release-notes.ts
export const releaseNotesTask = new Task(host, {
  task: async ([orgId]) => {
    const response = await fetch(`${config.backendUrl}/api/release-notes/${orgId}`);
    return response.json();
  },
  args: () => [host.orgId],
});
```

### CSS Namespacing

All widget CSS is namespaced to `.announcable-widget` to avoid conflicts when embedded on customer websites. Use Shadow DOM encapsulation where possible.

## Development Workflow

### With Backend

```bash
# Terminal 1: Start backend dev services
cd backend && make dev-services

# Terminal 2: Start backend
cd backend && make dev-air

# Terminal 3: Start widget dev server
cd widget && npm run dev
```

### Widget Only

```bash
npm run dev              # Dev server with hot-reload
npm run build            # Production UMD bundle
npm run build:test       # Dev/test build
npm run lint             # Lint check
```

## Adding a New Component

1. Create component file in appropriate directory (`components/ui/` or `components/widget/`)
2. Extend `BaseComponent` or `LitElement`
3. Use Lit decorators (`@customElement`, `@property`, `@state`)
4. Import and use in parent component

```typescript
import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('my-component')
export class MyComponent extends LitElement {
  @property() label = '';

  static styles = css`
    :host { display: block; }
  `;

  render() {
    return html`<div>${this.label}</div>`;
  }
}
```

## Adding a New API Task

1. Create task file in `src/tasks/`
2. Define the fetch logic using Lit `Task` pattern
3. Wire it into the component that needs the data
4. Handle loading, error, and success states

## Completion Checklist

- [ ] `npm run lint` passes
- [ ] `npm run build` succeeds
- [ ] Widget renders correctly in browser
- [ ] No console errors
- [ ] CSS doesn't leak outside widget (Shadow DOM / namespacing)
- [ ] API calls work against running backend
- [ ] No `any` types introduced

## Anti-Patterns

| Don't | Why | Instead |
|-------|-----|---------|
| Use global CSS | Leaks into host page | Use Shadow DOM or `.announcable-widget` namespace |
| Fetch directly in render | Re-fetches on every render | Use Lit Task pattern |
| Use `any` type | Hides errors | Use proper TypeScript types from `lib/types.ts` |
| Hardcode backend URL | Breaks in different environments | Use `lib/config.ts` |
| Modify DOM outside Shadow DOM | Unpredictable on host pages | Keep all DOM in component templates |
| Batch changes | Hard to identify breakage | One change → validate |
