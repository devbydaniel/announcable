# Announcable Widget - Lit Version

A lightweight, Web Components-based release notes widget built with Lit.

## Why Lit?

This is a migration from the React-based widget to Lit, providing:

- **~66% smaller bundle size**: From ~300KB (React) to ~100KB (Lit)
- **Better performance**: Faster initial load and runtime
- **Web Standards**: Built on native Web Components
- **Better encapsulation**: True Shadow DOM isolation
- **Framework agnostic**: Can be used anywhere

## Development

### Install Dependencies

```bash
npm install
```

### Development Server

```bash
npm run dev
```

Then open `http://localhost:5173/test.dev.html` to test the widget.

### Build

```bash
# Production build
npm run build

# Development build (for testing)
npm run build:test
```

### Preview Production Build

```bash
npm run preview
```

Then open `http://localhost:4173/test.prod.html` to test the built widget.

## Project Structure

```
src/
├── components/
│   ├── ui/           # UI components (button, card, etc.)
│   ├── widget/       # Widget type components (popover, modal, sidebar)
│   └── icons/        # Icon components
├── controllers/      # Reactive Controllers for state management
├── tasks/            # Data fetching tasks (replaces React Query)
├── lib/              # Utilities, types, contexts
├── main.ts           # Entry point
├── app.ts            # Main app component
└── index.css         # Global styles
```

## Lit Patterns

### Custom Elements

Components are defined as custom elements using decorators:

```typescript
import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

@customElement('my-component')
export class MyComponent extends LitElement {
  @property({ type: String }) name = '';
  @state() private count = 0;

  static styles = css`
    :host { display: block; }
  `;

  render() {
    return html`
      <div>Hello, ${this.name}! Count: ${this.count}</div>
    `;
  }
}
```

### Reactive Controllers

Controllers manage complex state and side effects:

```typescript
import { ReactiveController, ReactiveControllerHost } from 'lit';

export class MyController implements ReactiveController {
  host: ReactiveControllerHost;

  constructor(host: ReactiveControllerHost) {
    this.host = host;
    host.addController(this);
  }

  hostConnected() {
    // Setup when component connects to DOM
  }

  hostDisconnected() {
    // Cleanup when component disconnects
  }
}
```

### Tasks

Tasks handle async operations (like data fetching):

```typescript
import { Task } from '@lit/task';

this.task = new Task(
  this,
  async ([arg]) => {
    const response = await fetch(`/api/data/${arg}`);
    return response.json();
  },
  () => [this.someArg]
);

// In render:
${this.task.render({
  pending: () => html`Loading...`,
  complete: (data) => html`Data: ${data}`,
  error: (e) => html`Error: ${e}`,
})}
```

## Usage

The widget can be embedded in any website:

```html
<script>
  window.announcable_init = {
    org_id: 'YOUR_ORG_ID',
    anchor_query_selector: '[data-announcable]', // Optional
    hide_indicator: false, // Optional
    font_family: ['Inter', 'system-ui', 'sans-serif'] // Optional
  };
</script>
<script src="/path/to/widget.js"></script>
```

## Build Output

The build creates a UMD bundle at `dist/widget.js` that includes:
- All Lit code
- All component code
- Inlined CSS (Tailwind + custom styles)
- Everything needed to run standalone

## Comparison to React Widget

| Feature | React Widget | Lit Widget |
|---------|-------------|------------|
| Bundle Size | ~300KB | ~100KB |
| Initial Load | Slower | Faster |
| Framework | React 18 | Lit 3 (Web Components) |
| Dependencies | React, ReactDOM, Radix, TanStack | Lit, minimal utilities |
| Encapsulation | Shadow DOM | Shadow DOM |
| API | Same | Same (100% compatible) |

## Migration Status

See `_docs/20251123--widget-lit-migration/README.md` for the full migration plan.

Current status: **Phase 4 Complete** ✅

- [x] Phase 0: Project structure setup
- [x] Phase 1: Dependencies and configuration
- [x] Phase 2: Core architecture
- [x] Phase 3: State management & data fetching
- [x] Phase 4: Component migration ✅
- [ ] Phase 5: Advanced features
- [ ] Phase 6: Styling
- [ ] Phase 7: Testing
- [ ] Phase 8: Documentation
- [ ] Phase 9: Backend integration
