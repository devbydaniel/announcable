# Migration Plan: React Widget to Lit (New widget-lit folder)

## Overview
Create a new `widget-lit` folder to build the Lit version of the Announcable widget alongside the existing React widget. This allows parallel development and testing without disrupting the production React widget.

## Phase 0: New Project Structure Setup

### 0.1 Create New Widget Directory
**Create: `widget-lit/` folder** at root level (same level as `widget/`)

### 0.2 Copy Base Configuration Files
**Copy from `widget/` to `widget-lit/`:**
- `components.json`
- `eslint.config.js` (will update for Lit)
- `index.html`
- `postcss.config.js`
- `tailwind.config.cjs`
- `test.dev.html`
- `test.prod.html`
- `tsconfig.json` (will update for Lit)
- `tsconfig.app.json` (will update for Lit)
- `tsconfig.node.json`

**Create new: `widget-lit/package.json`**
- Start fresh with Lit dependencies (see Phase 1.1)

**Create new: `widget-lit/.gitignore`**
```
node_modules/
dist/
*.local
.DS_Store
```

### 0.3 Directory Structure
```
widget-lit/
├── src/
│   ├── components/
│   │   ├── ui/
│   │   └── widget/
│   ├── hooks/          (will become services/tasks)
│   ├── lib/
│   ├── assets/
│   ├── main.ts
│   ├── app.ts
│   └── index.css
├── public/
├── dist/
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.cjs
├── postcss.config.js
├── index.html
├── test.dev.html
├── test.prod.html
└── README.md
```

## Phase 1: Project Setup & Dependencies

### 1.1 Create package.json for widget-lit
**New File: `widget-lit/package.json`**

```json
{
  "name": "devbydaniel--release-notes-widget-lit",
  "version": "1.0.0",
  "files": [
    "dist"
  ],
  "main": "./dist/release-beacon-widget.umd.js",
  "module": "./dist/release-beacon-widget.es.js",
  "exports": {
    ".": {
      "import": "./dist/release-beacon-widget.es.js",
      "require": "./dist/release-beacon-widget.umd.js"
    }
  },
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "build:test": "tsc -b && vite build --mode development",
    "lint": "eslint .",
    "preview": "vite preview"
  },
  "dependencies": {
    "lit": "^3.1.0",
    "@lit/context": "^1.1.0",
    "@lit/task": "^1.0.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.1",
    "tailwind-merge": "^2.5.2",
    "tailwindcss-animate": "^1.0.7",
    "vite-plugin-css-injected-by-js": "^3.5.1"
  },
  "devDependencies": {
    "@types/node": "^22.2.0",
    "@types/postcss-prefix-selector": "^1.16.3",
    "autoprefixer": "^10.4.20",
    "eslint": "^9.8.0",
    "globals": "^15.9.0",
    "postcss": "^8.4.41",
    "postcss-plugin-namespace": "^0.0.3",
    "tailwindcss": "^3.4.9",
    "typescript": "^5.5.3",
    "typescript-eslint": "^8.0.0",
    "vite": "^5.4.0"
  }
}
```

### 1.2 Create Vite Configuration
**New File: `widget-lit/vite.config.ts`**

```typescript
import path from "path";
import { defineConfig } from "vite";
import namespace from "postcss-plugin-namespace";
import autoprefixer from "autoprefixer";
import tailwindcss from "tailwindcss";

export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwindcss(), namespace(".announcable-widget"), autoprefixer],
    },
  },
  plugins: [],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    lib: {
      entry: "src/main.ts",
      formats: ["umd"],
      name: "ReleaseBeaconWidget",
      fileName: () => `widget.js`,
    },
    rollupOptions: {
      external: [],
    },
    outDir: "dist",
    emptyOutDir: true,
  },
  define: {
    "process.env": {},
  },
});
```

### 1.3 Update TypeScript Configuration
**File: `widget-lit/tsconfig.json`**

Add to compilerOptions:
```json
{
  "compilerOptions": {
    "experimentalDecorators": true,
    "useDefineForClassFields": false,
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"]
}
```

### 1.4 Install Dependencies
**Command to run:**
```bash
cd widget-lit
npm install
```

## Phase 2: Core Architecture Setup

### 2.1 Copy and Update Library Files
**Copy from `widget/src/lib/` to `widget-lit/src/lib/`:**
- `clientId.ts` (no changes)
- `config.ts` (no changes)
- `utils.ts` (no changes)
- `types.ts` (update to remove React types)

**Update: `widget-lit/src/lib/types.ts`**
- Remove any `React.ReactNode` references
- Keep all data interfaces as-is

### 2.2 Create Lit Base Component
**New File: `widget-lit/src/lib/base-component.ts`**

```typescript
import { LitElement } from 'lit';

export class BaseComponent extends LitElement {
  // Base utilities for all widget components
  // Can add common methods here
}
```

### 2.3 Create Context System
**New File: `widget-lit/src/lib/contexts.ts`**

```typescript
import { createContext } from '@lit/context';
import type { WidgetConfig, WidgetInit } from './types';

export const widgetConfigContext = createContext<WidgetConfig>('widget-config');
export const widgetInitContext = createContext<WidgetInit>('widget-init');
```

### 2.4 Copy CSS
**Copy: `widget/src/index.css` → `widget-lit/src/index.css`**
- No changes needed

### 2.5 Copy Assets
**Copy: `widget/src/assets/` → `widget-lit/src/assets/`**
- Copy any static assets

## Phase 3: State Management & Data Fetching

### 3.1 Create Tasks Directory
**New Directory: `widget-lit/src/tasks/`**

This replaces the `hooks/` directory from React widget.

### 3.2 Create Data Fetching Tasks
**New File: `widget-lit/src/tasks/release-notes.ts`**

```typescript
import { Task } from '@lit/task';
import { ReactiveController, ReactiveControllerHost } from 'lit';
import type { ReleaseNote } from '@/lib/types';
import { backendUrl } from '@/lib/config';

export class ReleaseNotesController implements ReactiveController {
  host: ReactiveControllerHost;
  
  task: Task<[string], ReleaseNote[]>;

  constructor(host: ReactiveControllerHost, orgId: string) {
    this.host = host;
    host.addController(this);
    
    this.task = new Task(
      host,
      async ([orgId]) => {
        const url = `${backendUrl}/api/release-notes/${orgId}?for=widget`;
        const res = await fetch(url, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        });
        const { data } = await res.json();
        return (data || []) as ReleaseNote[];
      },
      () => [orgId]
    );
  }

  hostConnected() {}
  hostDisconnected() {}
}
```

**New File: `widget-lit/src/tasks/widget-config.ts`**

Similar pattern for fetching widget configuration.

**New File: `widget-lit/src/tasks/release-note-status.ts`**

Similar pattern for release note status.

**New File: `widget-lit/src/tasks/release-note-likes.ts`**

For like functionality with state management.

**New File: `widget-lit/src/tasks/release-note-metrics.ts`**

For view tracking with IntersectionObserver.

### 3.3 Create State Controllers
**New File: `widget-lit/src/controllers/widget-toggle.ts`**

```typescript
import { ReactiveController, ReactiveControllerHost } from 'lit';

export class WidgetToggleController implements ReactiveController {
  host: ReactiveControllerHost;
  isOpen = false;
  lastOpened: string | null = null;

  constructor(host: ReactiveControllerHost, querySelector?: string) {
    this.host = host;
    host.addController(this);
    this.loadLastOpened();
    this.setupAnchorClickListeners(querySelector);
  }

  private loadLastOpened() {
    // localStorage logic
  }

  private setupAnchorClickListeners(querySelector?: string) {
    // Setup click handlers
  }

  setIsOpen(value: boolean) {
    this.isOpen = value;
    if (value) {
      this.lastOpened = Date.now().toString();
      localStorage.setItem('announcable_last_opened', this.lastOpened);
    }
    this.host.requestUpdate();
  }

  hostConnected() {}
  hostDisconnected() {}
}
```

## Phase 4: Component Migration

### 4.1 Create Main Entry Point
**New File: `widget-lit/src/main.ts`**

```typescript
import { html, render } from 'lit';
import styles from "./index.css?inline";
import type { WidgetInit } from "./lib/types";
import './app'; // Register app component

declare global {
  interface Window {
    announcable_init?: WidgetInit;
    AnnouncableWidget?: {
      init?: (config: WidgetInit) => void;
    };
    release_beacon_widget_init?: WidgetInit;
    ReleaseBeaconWidget?: {
      init?: (config: WidgetInit) => void;
    };
  }
}

function initialize(init: WidgetInit) {
  if (!init.org_id) {
    console.error("org_id is required to initialize release notes widget");
    return;
  }

  const widgetRoot = document.createElement("div");
  widgetRoot.id = "announcable-widget-root";
  document.body.appendChild(widgetRoot);

  // Create a closed Shadow DOM
  const shadowRoot = widgetRoot.attachShadow({ mode: "closed" });

  // Inject styles into Shadow DOM first
  const styleElement = document.createElement("style");
  styleElement.textContent = styles;
  shadowRoot.appendChild(styleElement);

  // Create a container for the Lit app inside the Shadow DOM
  const appContainer = document.createElement("div");
  appContainer.className = "announcable-widget";
  shadowRoot.appendChild(appContainer);

  // Render Lit app
  render(
    html\`<announcable-app .init=\${init}></announcable-app>\`,
    appContainer
  );
}

// Expose initialization function globally
window.AnnouncableWidget = {
  init: (config: WidgetInit) => {
    initialize(config);
  },
};

// Use release beacon legacy init for backwards compatibility
window.ReleaseBeaconWidget = {
  init: (config: WidgetInit) => {
    initialize(config);
  },
};

// Automatically initialize if config is present
if (window.announcable_init && window.AnnouncableWidget?.init) {
  console.log("AnnouncableWidget init");
  window.AnnouncableWidget.init(window.announcable_init);
} else if (
  window.release_beacon_widget_init &&
  window.ReleaseBeaconWidget?.init
) {
  console.log("ReleaseBeaconWidget init");
  window.ReleaseBeaconWidget.init(window.release_beacon_widget_init);
} else {
  console.error("No widget init config found");
}
```

### 4.2 Create App Component
**New File: `widget-lit/src/app.ts`**

Convert `widget/src/App.tsx` to Lit component:
- Change to class extending `LitElement`
- Use `@customElement` decorator for registration
- Use `@property()` for props
- Use `@state()` for internal state
- Replace JSX with `html` tagged template literals
- Replace `useEffect` with Lit lifecycle methods (`updated()`, `firstUpdated()`)
- Replace React hooks with Tasks and reactive properties

Key patterns:
- `export default function App({ init }: Props)` → `@customElement('announcable-app') export class App extends LitElement`
- JSX → `render() { return html\`...\` }`
- `useState` → `@state()`
- `useEffect` → `updated()` or `willUpdate()`
- Event handlers: `@click=${this.handleClick}`
- Controllers for data fetching and state management

### 4.3 Migrate Widget Type Components
**New Files in `widget-lit/src/components/widget/`:**
- `popover.ts` - Convert from popover.tsx
- `modal.ts` - Convert from modal.tsx  
- `sidebar.ts` - Convert from sidebar.tsx
- `index.ts` - Main widget component selector

Each becomes a LitElement with:
- `@customElement` decorator for registration
- `@property()` decorators for reactive properties
- Lit template rendering with `html` tagged templates
- Style encapsulation via `static styles` property

Example structure:
```typescript
import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import type { WidgetConfig, WidgetInit } from '@/lib/types';

@customElement('widget-popover')
export class PopoverWidget extends LitElement {
  @property({ type: Object }) config!: WidgetConfig;
  @property({ type: Object }) init!: WidgetInit;
  @property({ type: Boolean }) isOpen = false;
  @property({ type: Function }) onClose!: () => void;

  static styles = css`
    /* Component styles */
  `;

  render() {
    return html`
      <!-- Convert JSX to Lit template -->
    `;
  }
}
```

### 4.4 Migrate UI Components
**New Directory: `widget-lit/src/components/ui/`**

Convert each React component to Lit:
- `button.ts`
- `card.ts`
- `dialog.ts` - Use native `<dialog>` or custom implementation
- `scroll-area.ts` - Use CSS overflow instead of Radix
- `skeleton.ts`
- `indicator.ts`
- `error-panel.ts`
- `release-notes-list.ts`

**Key Radix Replacements:**
- `@radix-ui/react-dialog` → Use native `<dialog>` element or custom Lit component
- `@radix-ui/react-scroll-area` → Use CSS `overflow: auto` with custom scrollbar styling
- `@radix-ui/react-slot` → Use Lit's slot system (`<slot></slot>`)

**Example: `widget-lit/src/components/ui/button.ts`**
```typescript
import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { classMap } from 'lit/directives/class-map.js';
import { cva } from 'class-variance-authority';

const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        ghost: "hover:bg-accent hover:text-accent-foreground",
      },
      size: {
        default: "h-10 px-4 py-2",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

@customElement('ui-button')
export class Button extends LitElement {
  @property() variant: 'default' | 'ghost' = 'default';
  @property() size: 'default' | 'icon' = 'default';

  static styles = css`
    :host {
      display: inline-block;
    }
  `;

  render() {
    const classes = buttonVariants({ variant: this.variant, size: this.size });
    return html`
      <button class=${classes}>
        <slot></slot>
      </button>
    `;
  }
}
```

### 4.5 Create Icon Components
**New Directory: `widget-lit/src/components/icons/`**

**Option 1: Inline SVG Components**
```typescript
import { LitElement, html, svg } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('icon-gift')
export class GiftIcon extends LitElement {
  @property({ type: Number }) size = 24;

  render() {
    return svg`
      <svg width=${this.size} height=${this.size} viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <!-- SVG path data -->
      </svg>
    `;
  }
}
```

**Create icon components for:**
- `gift.ts` - For main button
- `thumbs-up.ts` - For likes
- `external-link.ts` - For CTAs
- `x.ts` - For close button

## Phase 5: Advanced Features

### 5.1 IntersectionObserver for Metrics
**In release-note-metrics controller:**
- Use Lit lifecycle methods
- Attach observer in `hostConnected()`
- Clean up in `hostDisconnected()`
- Track when release notes come into view
- Send metrics to backend API

### 5.2 LocalStorage Management
**Create: `widget-lit/src/lib/storage.ts`**

Utility functions for:
- Storing/retrieving last opened timestamp
- Storing/retrieving client ID
- Managing liked release notes
- Other persistent state

### 5.3 Anchor Element Management
**In widget-toggle controller:**
- Query anchor elements on initialization using `querySelector`
- Attach click event listeners to open widget
- Update `data-new-release-notes` dataset attribute
- Update `data-instant-open` dataset attribute
- Handle instant open logic for new releases
- Clean up listeners on disconnect

## Phase 6: Styling

### 6.1 Tailwind Setup
**File: `widget-lit/tailwind.config.cjs`**
- Copy from React widget
- No changes needed
- Keeps all utility classes and animations

### 6.2 Component Styles
For each Lit component:
- Add `static styles` property for scoped styles using `css` tagged template
- Use for component-specific styles
- Reference Tailwind classes in `render()` templates via `class=""` attributes
- Inline critical styles, reference Tailwind for utilities

Example:
```typescript
import { LitElement, html, css } from 'lit';

export class MyComponent extends LitElement {
  static styles = css`
    :host {
      display: block;
    }
    
    /* Component-specific styles */
  `;
  
  render() {
    return html`
      <div class="flex flex-col gap-4">
        <!-- Use Tailwind classes -->
      </div>
    `;
  }
}
```

### 6.3 Test Styling Isolation
- Verify Shadow DOM prevents style leakage from host page
- Test widget on various host pages with different CSS
- Ensure custom font families work via `init.font_family`
- Verify Tailwind namespace `.announcable-widget` works correctly

## Phase 7: Testing & Integration

### 7.1 Update Test HTML Files
**File: `widget-lit/test.dev.html`**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Widget Lit Test - Development</title>
</head>
<body>
  <h1>Widget Lit Test - Development</h1>
  
  <button data-announcable>What's New?</button>

  <script>
    window.announcable_init = {
      org_id: 'YOUR_ORG_ID',
      anchor_query_selector: '[data-announcable]',
      hide_indicator: false,
      font_family: ['Inter', 'system-ui', 'sans-serif']
    };
  </script>
  <script type="module" src="/src/main.ts"></script>
</body>
</html>
```

**File: `widget-lit/test.prod.html`**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Widget Lit Test - Production</title>
</head>
<body>
  <h1>Widget Lit Test - Production</h1>
  
  <button data-announcable>What's New?</button>

  <script>
    window.announcable_init = {
      org_id: 'YOUR_ORG_ID',
      anchor_query_selector: '[data-announcable]',
      hide_indicator: false,
      font_family: ['Inter', 'system-ui', 'sans-serif']
    };
  </script>
  <script src="/dist/widget.js"></script>
</body>
</html>
```

### 7.2 Development Testing
**Start dev server:**
```bash
cd widget-lit
npm run dev
```

**Test in browser:**
- Open `http://localhost:5173/test.dev.html`
- Test all functionality with hot-reload

### 7.3 Production Build Testing
**Build and preview:**
```bash
cd widget-lit
npm run build
npm run preview
```

**Test in browser:**
- Open `http://localhost:4173/test.prod.html`
- Test minified production bundle

### 7.4 Feature Testing Checklist
Test all features:
- [ ] Widget initialization via script tag
- [ ] `announcable_init` initialization
- [ ] `release_beacon_widget_init` initialization (legacy)
- [ ] Popover widget type
- [ ] Modal widget type
- [ ] Sidebar widget type
- [ ] Release notes fetching and display
- [ ] Like/unlike functionality
- [ ] Like state persistence
- [ ] CTA click tracking
- [ ] View metrics (IntersectionObserver)
- [ ] Indicator on anchor elements (shows when unseen notes)
- [ ] Indicator on floating button (shows when unseen notes)
- [ ] Instant open mechanism (for attention_mechanism="instant_open")
- [ ] Custom anchor query selector
- [ ] Floating button (when no anchor provided)
- [ ] Image display
- [ ] iframe/video embed display
- [ ] Skeleton loading states
- [ ] Error states
- [ ] Error panel display
- [ ] Shadow DOM style isolation
- [ ] Custom font family via init.font_family
- [ ] Responsive design on mobile/tablet/desktop
- [ ] Close button functionality
- [ ] External link button
- [ ] Widget positioning (popover: bottom-right, sidebar: right side)
- [ ] Scroll functionality in long content
- [ ] Widget config customization (colors, borders, etc.)
- [ ] Hide indicator option (init.hide_indicator)
- [ ] Multiple anchor elements support
- [ ] Dataset updates on anchor elements

## Phase 8: Documentation

### 8.1 Create README
**New File: `widget-lit/README.md`**

Document:
- **Project Overview**: Lit-based release notes widget
- **Why Lit vs React**: Bundle size, performance, Web Components standard
- **Development Setup**: Installation and running dev server
- **Build Commands**: 
  - `npm run dev` - Development server with hot-reload
  - `npm run build` - Production build
  - `npm run build:test` - Development build for testing
  - `npm run preview` - Preview production build
- **Component Architecture**: 
  - App component structure
  - Widget types (popover, modal, sidebar)
  - UI components
  - Controllers and tasks
  - State management patterns
- **Custom Elements List**: Document all custom elements created
- **Usage Examples**: How to embed and initialize widget
- **Comparison with React widget**: Features, bundle size, performance
- **Development Patterns**: Lit patterns, decorators, reactive controllers

### 8.2 Add Code Comments
- Document complex controllers and their responsibilities
- Explain Lit patterns for developers unfamiliar with Lit
- Note differences from React version where relevant
- Add JSDoc comments for public APIs
- Comment tricky state management logic

### 8.3 Create Migration Notes
**New File: `widget-lit/MIGRATION-NOTES.md`**

Document:
- **Changes from React version**:
  - React hooks → Reactive controllers
  - React Context → @lit/context
  - TanStack Query → @lit/task
  - Radix UI → Native elements
  - JSX → Lit html templates
- **Feature Parity**: List what's the same and what's different
- **Bundle Size Improvements**: 
  - React widget: ~300KB (React + ReactDOM + Radix + TanStack)
  - Lit widget: ~100KB (Lit core is only ~15KB)
  - ~66% reduction in bundle size
- **Performance Benchmarks**:
  - Initial load time
  - Time to interactive
  - Runtime performance
  - Memory usage
- **Breaking Changes**: Document any API differences (should be none)
- **Migration Timeline**: When this will replace React widget

## Phase 9: Backend Integration (Future)

**When ready to switch from React to Lit widget:**

### 9.1 Update Backend Build
**File: `backend/Makefile`** (if widget build is automated)
- Add target for building widget-lit:
  ```makefile
  widget-lit-build:
  	cd ../widget-lit && npm install && npm run build
  	cp ../widget-lit/dist/widget.js ./static/widget-lit.js
  ```
- Or switch from `widget` to `widget-lit` in existing target

### 9.2 Update Static Assets
**File: `backend/static/static.go`**
- Option 1: Replace React widget entirely
  - Point embed to `widget-lit/dist` instead of `widget/dist`
- Option 2: Serve both and use feature flag
  - Embed both `widget/dist` and `widget-lit/dist`
  - Serve at different paths: `/static/widget.js` (React) and `/static/widget-lit.js` (Lit)

### 9.3 Update Templates
**Files: `backend/templates/pages/*`** (where widget is embedded)
- No changes needed to widget initialization API (100% compatible)
- If serving different path, update script src:
  ```html
  <!-- Old React widget -->
  <script src="/static/widget.js"></script>
  
  <!-- New Lit widget -->
  <script src="/static/widget-lit.js"></script>
  ```

### 9.4 Feature Flag (Optional)
**For gradual rollout:**

**Add to `backend/config/config.go`:**
```go
type Config struct {
  // ... existing config
  UseLitWidget bool `env:"USE_LIT_WIDGET" envDefault:"false"`
}
```

**Update template rendering to conditionally load widget:**
```go
// In template data
data := BaseTemplateData{
  // ... existing data
  WidgetScriptSrc: getWidgetScriptSrc(cfg.UseLitWidget),
}

func getWidgetScriptSrc(useLit bool) string {
  if useLit {
    return "/static/widget-lit.js"
  }
  return "/static/widget.js"
}
```

This allows:
- A/B testing both versions
- Gradual rollout by organization
- Quick rollback if issues arise
- Performance comparison in production

## Implementation Order Summary

1. **Setup** (Phase 0-1)
   - Create `widget-lit/` folder structure
   - Copy configuration files from `widget/`
   - Create `package.json` with Lit dependencies
   - Update `vite.config.ts` (remove React plugin)
   - Update `tsconfig.json` (add decorator support)
   - Run `npm install`

2. **Foundation** (Phase 2-3)
   - Copy `lib/` files (clientId, config, utils, types)
   - Update `types.ts` to remove React types
   - Create contexts using `@lit/context`
   - Create controllers for state management
   - Create tasks for data fetching (replace TanStack Query)
   - Setup widget-toggle, release-notes, widget-config tasks

3. **Core Components** (Phase 4)
   - Create `main.ts` entry point (Shadow DOM setup)
   - Create `app.ts` component (main app logic)
   - Migrate widget types to Lit (popover, modal, sidebar)
   - Migrate UI components to Lit (button, card, dialog, etc.)
   - Create icon components (inline SVG)
   - Replace Radix components with native/custom elements

4. **Features** (Phase 5)
   - Implement IntersectionObserver for view metrics
   - Setup localStorage utilities
   - Handle anchor element management
   - Implement instant open mechanism
   - Add like/unlike functionality
   - Add CTA tracking

5. **Polish** (Phase 6)
   - Configure Tailwind (copy from React widget)
   - Add component-specific styles with `css` tagged templates
   - Test Shadow DOM isolation
   - Ensure custom fonts work
   - Verify all styling matches React widget

6. **Testing** (Phase 7)
   - Create `test.dev.html` for development testing
   - Create `test.prod.html` for production testing
   - Test all features against checklist
   - Test on different browsers
   - Test on different host page styles
   - Verify bundle size reduction

7. **Documentation** (Phase 8)
   - Write comprehensive `README.md`
   - Add code comments throughout
   - Create `MIGRATION-NOTES.md`
   - Document all custom elements
   - Document differences from React version

8. **Integration** (Phase 9 - Later)
   - When ready, integrate with backend
   - Consider feature flag for gradual rollout
   - Update Makefile for automated builds
   - Update static asset embedding
   - Deploy and monitor

## Expected Benefits

### Bundle Size
- **React Widget**: ~300KB total
  - React: ~45KB
  - ReactDOM: ~130KB
  - Radix UI: ~50KB
  - TanStack Query: ~40KB
  - Other dependencies: ~35KB
- **Lit Widget**: ~100KB total
  - Lit core: ~15KB
  - Other dependencies: ~85KB (Tailwind, utilities)
- **Reduction**: ~66% smaller bundle size

### Performance
- **Faster Initial Load**: Smaller bundle = faster download and parse
- **Better Runtime Performance**: Lit is closer to native DOM manipulation
- **Lower Memory Usage**: No virtual DOM overhead
- **Faster Time to Interactive**: Less JavaScript to execute

### Standards & Compatibility
- **Native Web Components**: Built on browser standards
- **Better Encapsulation**: Shadow DOM prevents style leakage
- **No Framework Lock-in**: Can be used with any framework
- **Future-proof**: Built on stable web platform APIs

### Maintenance
- **Simpler Dependency Tree**: Fewer dependencies to update
- **Less Framework Churn**: Web Components are stable
- **Easier Debugging**: Closer to native browser APIs
- **Better TypeScript Support**: Decorators for reactive properties

### Parallel Development
- **No Disruption**: React widget continues working
- **Easy Comparison**: Can test both side-by-side
- **Gradual Migration**: Backend can switch when ready
- **Risk Mitigation**: Can roll back easily if needed

## Risks & Mitigations

### Risk: Breaking Changes in Widget API
**Mitigation**: 
- Keep exact same initialization API
- Test backward compatibility with legacy `release_beacon_widget_init`
- Maintain same global object structure
- No changes to embed code required

### Risk: Missing Radix UI Accessibility Features
**Mitigation**: 
- Use native `<dialog>` element which has built-in accessibility
- Add proper ARIA attributes manually
- Test with screen readers
- Follow WCAG guidelines for all interactive elements

### Risk: Complex State Management Without React Hooks
**Mitigation**: 
- Use `@lit/context` for context (similar to React Context)
- Use `@lit/task` for async data fetching (similar to TanStack Query)
- Use Reactive Controllers for complex state logic
- Patterns are similar, just different syntax

### Risk: Learning Curve for Team
**Mitigation**: 
- Lit is actually simpler than React
- Good documentation at lit.dev
- Similar patterns to React (lifecycle, reactivity)
- Create this comprehensive migration guide
- Add extensive code comments

### Risk: Bundle Size Not As Small As Expected
**Mitigation**: 
- Measure actual bundle size after build
- Use tree-shaking effectively
- Only include necessary Lit features
- Minimize Tailwind output with purge

### Risk: Shadow DOM Compatibility Issues
**Mitigation**: 
- Test on all major browsers
- Test on various host pages
- Verify CSS isolation works
- Test with different CSP policies

## Success Metrics

Track these metrics to measure success:

### Technical Metrics
- [ ] Bundle size reduced by >50%
- [ ] Initial load time improved by >30%
- [ ] Time to interactive improved by >25%
- [ ] Memory usage reduced by >20%
- [ ] All features work identically to React version

### Quality Metrics
- [ ] Zero breaking changes to API
- [ ] All tests pass
- [ ] No accessibility regressions
- [ ] Works on all supported browsers
- [ ] Shadow DOM isolation verified

### Development Metrics
- [ ] Complete, clear documentation
- [ ] All code properly commented
- [ ] Easy to build and test locally
- [ ] CI/CD pipeline updated
- [ ] Team trained on Lit patterns

## Timeline Estimate

**Phase 0-1 (Setup)**: 1 day
- Create folder structure
- Setup dependencies
- Configure build tools

**Phase 2-3 (Foundation)**: 2-3 days
- Copy and adapt libraries
- Create controllers
- Create tasks for data fetching
- Setup state management

**Phase 4 (Components)**: 5-7 days
- Main entry point and app component: 1 day
- Widget types (popover, modal, sidebar): 1 day
- UI components (8 components): 2-3 days
- Icons: 1 day
- Integration and debugging: 1-2 days

**Phase 5 (Features)**: 2-3 days
- IntersectionObserver: 0.5 day
- LocalStorage: 0.5 day
- Anchor management: 1-2 days

**Phase 6 (Styling)**: 1-2 days
- Tailwind setup: 0.5 day
- Component styles: 0.5 day
- Testing isolation: 0.5-1 day

**Phase 7 (Testing)**: 2-3 days
- Setup test files: 0.5 day
- Development testing: 1 day
- Production testing: 0.5 day
- Feature checklist: 1 day

**Phase 8 (Documentation)**: 1-2 days
- README: 0.5 day
- Code comments: 0.5 day
- Migration notes: 0.5 day

**Total Estimate**: 15-22 days (3-4 weeks)

**Phase 9 (Integration)**: When ready to deploy
- Backend integration: 1 day
- Feature flag setup: 0.5 day
- Deploy and monitor: 0.5-1 day

## Final Directory Structure

```
mono/
├── widget/                    # Existing React widget (unchanged)
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── lib/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── dist/
│   ├── package.json
│   └── vite.config.ts
│
├── widget-lit/                # New Lit widget
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/            # UI components as LitElements
│   │   │   │   ├── button.ts
│   │   │   │   ├── card.ts
│   │   │   │   ├── dialog.ts
│   │   │   │   ├── scroll-area.ts
│   │   │   │   ├── skeleton.ts
│   │   │   │   ├── indicator.ts
│   │   │   │   ├── error-panel.ts
│   │   │   │   └── release-notes-list.ts
│   │   │   ├── widget/        # Widget type components
│   │   │   │   ├── popover.ts
│   │   │   │   ├── modal.ts
│   │   │   │   ├── sidebar.ts
│   │   │   │   └── index.ts
│   │   │   └── icons/         # Icon components
│   │   │       ├── gift.ts
│   │   │       ├── thumbs-up.ts
│   │   │       ├── external-link.ts
│   │   │       └── x.ts
│   │   ├── controllers/       # Reactive Controllers
│   │   │   └── widget-toggle.ts
│   │   ├── tasks/             # Data fetching tasks
│   │   │   ├── release-notes.ts
│   │   │   ├── widget-config.ts
│   │   │   ├── release-note-status.ts
│   │   │   ├── release-note-likes.ts
│   │   │   └── release-note-metrics.ts
│   │   ├── lib/               # Utilities and types
│   │   │   ├── base-component.ts
│   │   │   ├── clientId.ts
│   │   │   ├── config.ts
│   │   │   ├── contexts.ts
│   │   │   ├── storage.ts
│   │   │   ├── types.ts
│   │   │   └── utils.ts
│   │   ├── assets/            # Static assets
│   │   ├── main.ts            # Entry point
│   │   ├── app.ts             # Main app component
│   │   └── index.css          # Global styles
│   ├── dist/
│   │   └── widget.js          # Built bundle
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tailwind.config.cjs
│   ├── test.dev.html
│   ├── test.prod.html
│   ├── README.md
│   └── MIGRATION-NOTES.md
│
└── backend/                   # Backend serves widget
    └── static/
        ├── widget.js          # React widget (current)
        └── widget-lit.js      # Lit widget (future)
```

## Conclusion

This migration plan provides a comprehensive roadmap for creating a Lit-based version of the Announcable widget alongside the existing React implementation. The approach prioritizes:

1. **Safety**: No changes to existing React widget
2. **Compatibility**: Same API, no breaking changes
3. **Performance**: Significant bundle size reduction
4. **Standards**: Built on Web Components
5. **Maintainability**: Simpler dependencies, better encapsulation

By following this plan, the team can confidently migrate to Lit while maintaining full backward compatibility and having the ability to roll back if needed.
