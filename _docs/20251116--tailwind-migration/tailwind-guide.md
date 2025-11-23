# Tailwind CSS Development Guide

This guide explains how to work with Tailwind CSS in the Announcable backend after migration.

## Quick Start

### Running the CSS Build

**Development** (watch mode with hot reload):
```bash
cd backend
npm run css:dev
```

**Production** (minified):
```bash
cd backend
npm run css:build
```

The CSS build runs automatically in Docker/Air, but you may need to run it manually during local development.

## Common Patterns

### Layout Patterns

#### Centered Content Container
```html
<div class="max-w-7xl mx-auto px-4 py-8">
  <!-- content -->
</div>
```

#### Flexbox Row with Gap
```html
<div class="flex items-center gap-4">
  <span>Item 1</span>
  <span>Item 2</span>
</div>
```

#### Grid Layout
```html
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
  <!-- grid items -->
</div>
```

#### Sticky Header
```html
<header class="sticky top-0 z-10 bg-white shadow">
  <!-- header content -->
</header>
```

---

### Form Patterns

#### Standard Form with Spacing
```html
<form class="space-y-6" hx-post="/submit">
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">
      Email
    </label>
    <input
      type="email"
      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
    />
  </div>

  <button class="w-full bg-blue-600 text-white px-4 py-2 rounded-md font-medium hover:bg-blue-700">
    Submit
  </button>
</form>
```

#### Form with Validation Error
```html
<div>
  <label class="block text-sm font-medium text-gray-700 mb-1">
    Email
  </label>
  <input
    type="email"
    class="w-full px-3 py-2 border border-red-300 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-red-500"
  />
  <p class="text-sm text-red-600 mt-1">Please enter a valid email</p>
</div>
```

#### Checkbox with Label
```html
<label class="flex items-center gap-2 cursor-pointer">
  <input type="checkbox" class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
  <span class="text-sm text-gray-700">Remember me</span>
</label>
```

---

### Button Patterns

#### Primary Button
```html
<button class="px-4 py-2 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors">
  Save Changes
</button>
```

#### Secondary Button
```html
<button class="px-4 py-2 bg-gray-100 text-gray-700 font-medium rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors">
  Cancel
</button>
```

#### Danger Button
```html
<button class="px-4 py-2 bg-red-600 text-white font-medium rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-colors">
  Delete
</button>
```

#### Icon Button
```html
<button class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition-colors">
  <i data-feather="edit" class="w-5 h-5"></i>
</button>
```

#### Button Group
```html
<div class="flex gap-2">
  <button class="px-4 py-2 bg-blue-600 text-white rounded-md">Save</button>
  <button class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md">Cancel</button>
</div>
```

---

### Card Patterns

#### Basic Card
```html
<div class="bg-white rounded-lg shadow-md p-6">
  <h3 class="text-xl font-semibold mb-4">Card Title</h3>
  <p class="text-gray-600">Card content goes here</p>
</div>
```

#### Card with Header and Footer
```html
<div class="bg-white rounded-lg shadow-md overflow-hidden">
  <div class="px-6 py-4 border-b border-gray-200">
    <h3 class="text-xl font-semibold">Card Title</h3>
  </div>
  <div class="px-6 py-4">
    <p class="text-gray-600">Card content</p>
  </div>
  <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end gap-2">
    <button class="px-4 py-2 bg-blue-600 text-white rounded-md">Save</button>
    <button class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md">Cancel</button>
  </div>
</div>
```

#### Clickable Card
```html
<a href="/details" class="block bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
  <h3 class="text-xl font-semibold mb-2">Card Title</h3>
  <p class="text-gray-600">Click to view details</p>
</a>
```

---

### Table Patterns

#### Standard Table
```html
<div class="overflow-x-auto">
  <table class="w-full border-collapse">
    <thead>
      <tr class="bg-gray-50 border-b border-gray-200">
        <th class="px-4 py-3 text-left text-sm font-semibold text-gray-700">Name</th>
        <th class="px-4 py-3 text-left text-sm font-semibold text-gray-700">Email</th>
        <th class="px-4 py-3 text-left text-sm font-semibold text-gray-700">Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr class="border-b border-gray-200 hover:bg-gray-50">
        <td class="px-4 py-3 text-sm">John Doe</td>
        <td class="px-4 py-3 text-sm">john@example.com</td>
        <td class="px-4 py-3 text-sm">
          <button class="text-blue-600 hover:underline">Edit</button>
        </td>
      </tr>
    </tbody>
  </table>
</div>
```

---

### Alert/Message Patterns

#### Info Alert
```html
<div class="px-4 py-3 rounded-md bg-blue-50 text-blue-800 border-l-4 border-blue-500">
  <p class="text-sm">This is an informational message.</p>
</div>
```

#### Success Alert
```html
<div class="px-4 py-3 rounded-md bg-green-50 text-green-800 border-l-4 border-green-500">
  <p class="text-sm">Operation completed successfully!</p>
</div>
```

#### Warning Alert
```html
<div class="px-4 py-3 rounded-md bg-yellow-50 text-yellow-800 border-l-4 border-yellow-500">
  <p class="text-sm">Please review this warning.</p>
</div>
```

#### Error Alert
```html
<div class="px-4 py-3 rounded-md bg-red-50 text-red-800 border-l-4 border-red-500">
  <p class="text-sm">An error occurred. Please try again.</p>
</div>
```

---

### Badge Patterns

#### Status Badge
```html
<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
  Active
</span>
```

#### Count Badge
```html
<span class="inline-flex items-center justify-center w-6 h-6 text-xs font-bold text-white bg-red-500 rounded-full">
  3
</span>
```

---

## HTMX Integration

### HTMX Loading States

Use `htmx:indicator` class for loading spinners:

```html
<button
  class="px-4 py-2 bg-blue-600 text-white rounded-md relative"
  hx-post="/submit"
  hx-indicator="#spinner"
>
  <span id="spinner" class="htmx-indicator absolute inset-0 flex items-center justify-center bg-blue-600 rounded-md">
    <svg class="animate-spin h-5 w-5 text-white" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  </span>
  Submit
</button>
```

Add this CSS for indicator behavior:
```css
.htmx-indicator {
  display: none;
}

.htmx-request .htmx-indicator {
  display: flex;
}
```

### HTMX Swap Transitions

Use Alpine.js transitions with HTMX:

```html
<div
  x-data="{ show: false }"
  x-show="show"
  x-transition
  hx-get="/content"
  hx-trigger="load"
  @htmx:after-swap="show = true"
  class="bg-white p-4 rounded-md shadow"
>
  Content will fade in
</div>
```

---

## Alpine.js Integration

### Toggle Visibility
```html
<div x-data="{ open: false }">
  <button @click="open = !open" class="px-4 py-2 bg-blue-600 text-white rounded-md">
    Toggle
  </button>
  <div x-show="open" x-transition class="mt-4 p-4 bg-gray-100 rounded-md">
    Hidden content
  </div>
</div>
```

### Dropdown Menu
```html
<div x-data="{ open: false }" @click.away="open = false" class="relative">
  <button @click="open = !open" class="px-4 py-2 bg-gray-100 rounded-md">
    Menu
  </button>
  <div
    x-show="open"
    x-transition
    class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-10"
  >
    <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Option 1</a>
    <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Option 2</a>
  </div>
</div>
```

---

## Responsive Design

### Mobile-First Approach

Tailwind is mobile-first. Start with mobile styles, then add responsive modifiers:

```html
<!-- Stack on mobile, row on desktop -->
<div class="flex flex-col md:flex-row gap-4">
  <div class="w-full md:w-1/2">Column 1</div>
  <div class="w-full md:w-1/2">Column 2</div>
</div>
```

### Responsive Grid
```html
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  <!-- Items will stack on mobile, 2 columns on tablet, 3 on desktop, 4 on xl -->
</div>
```

### Hide/Show on Breakpoints
```html
<!-- Show on mobile only -->
<div class="block md:hidden">Mobile menu</div>

<!-- Show on desktop only -->
<div class="hidden md:block">Desktop menu</div>
```

---

## Dark Mode (Future)

Tailwind supports dark mode via class or media query strategy. For future dark mode support:

```html
<div class="bg-white dark:bg-gray-800 text-gray-900 dark:text-white">
  Content adapts to dark mode
</div>
```

Enable in `tailwind.config.js`:
```js
module.exports = {
  darkMode: 'class', // or 'media'
  // ...
}
```

---

## Performance Tips

### 1. Use Arbitrary Values Sparingly

❌ Avoid:
```html
<div class="p-[13px]">
```

✅ Prefer theme values:
```html
<div class="p-3">  <!-- 0.75rem = 12px -->
```

### 2. Don't Duplicate Long Class Lists

If you have a complex component used multiple times, extract to `@layer components`:

```css
@layer components {
  .btn-primary {
    @apply px-4 py-2 bg-blue-600 text-white font-medium rounded-md;
    @apply hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500;
  }
}
```

### 3. Purge Configuration

Ensure your `tailwind.config.js` `content` array covers all templates:

```js
module.exports = {
  content: [
    "./templates/**/*.html",
    "./internal/handler/**/*.go", // For string templates in Go
  ],
  // ...
}
```

---

## Common Gotchas

### 1. Go Template Syntax Conflicts

Go templates use `{{ }}` which can conflict. Be careful with spacing:

❌ Might break:
```html
<div class="{{if .Active}}bg-blue-500{{else}}bg-gray-500{{end}}">
```

✅ Better:
```html
<div class="{{ if .Active }}bg-blue-500{{ else }}bg-gray-500{{ end }}">
```

### 2. Dynamic Classes Must Be Complete

Tailwind's purge can't detect partial class names:

❌ Won't work:
```html
<div class="text-{{ .Color }}-500">  <!-- Purge will remove this -->
```

✅ Use complete class names:
```html
{{ if eq .Color "blue" }}
  <div class="text-blue-500">
{{ else if eq .Color "red" }}
  <div class="text-red-500">
{{ end }}
```

### 3. Focus Ring Offset Needs Background

For buttons with colored backgrounds, use `focus:ring-offset-2` with a background:

```html
<button class="bg-blue-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2">
  <!-- Ring offset creates space between button and ring -->
</button>
```

---

## Debugging Tips

### 1. Use Browser DevTools

Inspect elements to see which Tailwind classes are applied. Tailwind generates readable CSS class names.

### 2. Check Purge Configuration

If styles are missing in production, check that your `content` array in `tailwind.config.js` includes all template files.

### 3. Validate Output CSS

Check `static/css/output.css` to see generated CSS. In dev mode, it's ~3-5MB. In production (purged), it should be < 50KB.

---

## Resources

- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [Tailwind CSS Cheat Sheet](https://nerdcave.com/tailwind-cheat-sheet)
- [Tailwind Play (Online Playground)](https://play.tailwindcss.com/)
- [Headless UI (Components)](https://headlessui.com/)
- [Tailwind UI (Premium Components)](https://tailwindui.com/)

---

## When to Use Custom CSS

Use custom CSS in `@layer components` for:

1. **Complex animations**: Keyframe animations beyond simple transitions
2. **Pseudo-elements**: `::before`, `::after` with complex content
3. **Complex selectors**: Child combinators, sibling selectors
4. **Browser-specific hacks**: `-webkit-` prefixes, etc.
5. **Repeated complex patterns**: If a pattern has 8+ utilities and repeats 10+ times

Otherwise, prefer Tailwind utilities for maintainability and consistency.
