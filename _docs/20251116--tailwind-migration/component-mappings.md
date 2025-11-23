# Tailwind Component Mappings

This document maps existing CSS components to their Tailwind equivalents. Use this as a reference during template migration.

## Mapping Strategy

For each component, we document:
1. **Original CSS**: The custom CSS class structure
2. **Tailwind Equivalent**: Utility classes that replicate the styling
3. **Custom CSS Needed**: Any remaining custom CSS required (animations, complex selectors, etc.)

---

## Form Components

### Form Container

**Original**: `.form`
```css
.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
```

**Tailwind Equivalent**:
```html
<form class="space-y-4">
  <!-- form fields -->
</form>
```

**Custom CSS Needed**: None

---

### Form Group

**Original**: `.form__group`
```css
.form__group {
  display: flex;
  flex-direction: column;
  margin-bottom: 1rem;
}
```

**Tailwind Equivalent**:
```html
<div class="mb-4">
  <label>...</label>
  <input>
</div>
```

Or use `space-y-4` on parent form and no margin on children.

**Custom CSS Needed**: None

---

### Form Label

**Original**: `.form__label`
```css
.form__label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.25rem;
}
```

**Tailwind Equivalent**:
```html
<label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
```

**Custom CSS Needed**: None

---

### Form Input

**Original**: `.form__input`
```css
.form__input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.form__input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}
```

**Tailwind Equivalent**:
```html
<input class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
```

**With @tailwindcss/forms plugin**:
```html
<input class="rounded-md border-gray-300 focus:border-blue-500 focus:ring-blue-500">
```

**Custom CSS Needed**: None (use @tailwindcss/forms plugin)

---

### Form Subtext

**Original**: `.form__subtext`
```css
.form__subtext {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.25rem;
}
```

**Tailwind Equivalent**:
```html
<p class="text-xs text-gray-500 mt-1">Forgot password?</p>
```

**Custom CSS Needed**: None

---

## Button Components

### Base Button

**Original**: `.button`
```css
.button {
  padding: 0.5rem 1rem;
  font-weight: 600;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}
```

**Tailwind Equivalent**:
```html
<button class="px-4 py-2 font-semibold rounded-md cursor-pointer transition-all duration-200">
```

**Custom CSS Needed**: None

---

### Primary Button

**Original**: `.button--primary`
```css
.button--primary {
  background-color: #3b82f6;
  color: white;
}

.button--primary:hover {
  background-color: #2563eb;
}
```

**Tailwind Equivalent**:
```html
<button class="px-4 py-2 font-semibold rounded-md bg-blue-600 text-white hover:bg-blue-700 transition-colors">
```

**Custom CSS Needed**: None

---

### Secondary Button

**Original**: `.button--secondary`
```css
.button--secondary {
  background-color: #f3f4f6;
  color: #374151;
}

.button--secondary:hover {
  background-color: #e5e7eb;
}
```

**Tailwind Equivalent**:
```html
<button class="px-4 py-2 font-semibold rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors">
```

**Custom CSS Needed**: None

---

### Danger Button

**Original**: `.button--danger`
```css
.button--danger {
  background-color: #ef4444;
  color: white;
}

.button--danger:hover {
  background-color: #dc2626;
}
```

**Tailwind Equivalent**:
```html
<button class="px-4 py-2 font-semibold rounded-md bg-red-500 text-white hover:bg-red-600 transition-colors">
```

**Custom CSS Needed**: None

---

### Block Button

**Original**: `.button--block`
```css
.button--block {
  width: 100%;
}
```

**Tailwind Equivalent**:
```html
<button class="w-full px-4 py-2 ...">
```

**Custom CSS Needed**: None

---

### Small Button

**Original**: `.button--sm`
```css
.button--sm {
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem;
}
```

**Tailwind Equivalent**:
```html
<button class="px-3 py-1 text-sm font-semibold rounded-md ...">
```

**Custom CSS Needed**: None

---

## Card Components

### Card

**Original**: `.card`
```css
.card {
  background-color: white;
  border-radius: 0.5rem;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
```

**Tailwind Equivalent**:
```html
<div class="bg-white rounded-lg p-6 shadow-md">
```

**Custom CSS Needed**: None

---

### Card Title

**Original**: `.card__title`
```css
.card__title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #111827;
}
```

**Tailwind Equivalent**:
```html
<h2 class="text-2xl font-semibold mb-4 text-gray-900">Login</h2>
```

**Custom CSS Needed**: None

---

### Card Section

**Original**: `.card__section`
```css
.card__section {
  padding: 1rem 0;
  border-top: 1px solid #e5e7eb;
}
```

**Tailwind Equivalent**:
```html
<div class="py-4 border-t border-gray-200">
```

**Custom CSS Needed**: None

---

## Navigation Components

### Nav Container

**Original**: `.nav`
```css
.nav {
  display: flex;
  flex-direction: column;
  background-color: #1f2937;
  height: 100vh;
  width: 240px;
  padding: 1.5rem 1rem;
}
```

**Tailwind Equivalent**:
```html
<nav class="flex flex-col bg-gray-800 h-screen w-60 px-4 py-6">
```

**Custom CSS Needed**: May need sticky positioning or overflow handling

---

### Nav Brand

**Original**: `.nav__brand`
```css
.nav__brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  margin-bottom: 2rem;
}
```

**Tailwind Equivalent**:
```html
<div class="flex items-center gap-2 text-xl font-bold text-white mb-8">
  <img src="..." width="16" height="16" alt="" />
  <span>Announcable</span>
</div>
```

**Custom CSS Needed**: None

---

### Nav List

**Original**: `.nav__list`
```css
.nav__list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
```

**Tailwind Equivalent**:
```html
<ul class="space-y-1">
```

**Custom CSS Needed**: None

---

### Nav List Item

**Original**: `.nav__list__item a`
```css
.nav__list__item a {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.625rem 1rem;
  color: #d1d5db;
  border-radius: 0.375rem;
  text-decoration: none;
  transition: all 0.2s;
}

.nav__list__item a:hover {
  background-color: #374151;
  color: white;
}

.nav__list__item--active a {
  background-color: #3b82f6;
  color: white;
}
```

**Tailwind Equivalent**:
```html
<li>
  <a href="..." class="flex items-center gap-3 px-4 py-2.5 text-gray-300 rounded-md hover:bg-gray-700 hover:text-white transition-colors">
    <i data-feather="list" class="w-4 h-4"></i>
    <span>Release Notes</span>
  </a>
</li>

<!-- Active state (conditionally applied) -->
<li>
  <a href="..." class="flex items-center gap-3 px-4 py-2.5 bg-blue-600 text-white rounded-md">
    ...
  </a>
</li>
```

**Custom CSS Needed**: Active state detection logic (may need server-side template logic)

---

### Nav Divider

**Original**: `.nav__divider`
```css
.nav__divider {
  height: 1px;
  background-color: #4b5563;
  margin: 1rem 0;
}
```

**Tailwind Equivalent**:
```html
<div class="h-px bg-gray-600 my-4"></div>
```

**Custom CSS Needed**: None

---

## Layout Components

### App Frame Layout

**Original**: `.app`
```css
.app {
  display: grid;
  grid-template-columns: 240px 1fr;
  height: 100vh;
}
```

**Tailwind Equivalent**:
```html
<div class="grid grid-cols-[240px_1fr] h-screen">
```

Or using flexbox:
```html
<div class="flex h-screen">
  <nav class="w-60 flex-shrink-0">...</nav>
  <main class="flex-1">...</main>
</div>
```

**Custom CSS Needed**: None

---

### App Main

**Original**: `.app__main`
```css
.app__main {
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
```

**Tailwind Equivalent**:
```html
<main class="flex flex-col overflow-y-auto">
```

**Custom CSS Needed**: None

---

### App Content

**Original**: `.app__content`
```css
.app__content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}
```

**Tailwind Equivalent**:
```html
<div class="flex-1 p-8 max-w-7xl mx-auto w-full">
```

**Custom CSS Needed**: None

---

### Onboard Layout (Centered)

**Original**: `.onboard`
```css
.onboard {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #f9fafb;
  padding: 1rem;
}
```

**Tailwind Equivalent**:
```html
<div class="flex items-center justify-center min-h-screen bg-gray-50 p-4">
```

**Custom CSS Needed**: None

---

## Alert Components

### Alert

**Original**: `.alert`
```css
.alert {
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}
```

**Tailwind Equivalent**:
```html
<div class="px-4 py-3 rounded-md mb-4 text-sm">
```

**Custom CSS Needed**: None

---

### Alert Info

**Original**: `.alert--info`
```css
.alert--info {
  background-color: #dbeafe;
  color: #1e40af;
  border-left: 4px solid #3b82f6;
}
```

**Tailwind Equivalent**:
```html
<div class="px-4 py-3 rounded-md mb-4 text-sm bg-blue-50 text-blue-800 border-l-4 border-blue-500">
```

**Custom CSS Needed**: None

---

### Alert Success

**Original**: `.alert--success`
```css
.alert--success {
  background-color: #d1fae5;
  color: #065f46;
  border-left: 4px solid #10b981;
}
```

**Tailwind Equivalent**:
```html
<div class="px-4 py-3 rounded-md mb-4 text-sm bg-green-50 text-green-800 border-l-4 border-green-500">
```

**Custom CSS Needed**: None

---

### Alert Warning

**Original**: `.alert--warning`
```css
.alert--warning {
  background-color: #fef3c7;
  color: #92400e;
  border-left: 4px solid #f59e0b;
}
```

**Tailwind Equivalent**:
```html
<div class="px-4 py-3 rounded-md mb-4 text-sm bg-yellow-50 text-yellow-800 border-l-4 border-yellow-500">
```

**Custom CSS Needed**: None

---

### Alert Danger

**Original**: `.alert--danger`
```css
.alert--danger {
  background-color: #fee2e2;
  color: #991b1b;
  border-left: 4px solid #ef4444;
}
```

**Tailwind Equivalent**:
```html
<div class="px-4 py-3 rounded-md mb-4 text-sm bg-red-50 text-red-800 border-l-4 border-red-500">
```

**Custom CSS Needed**: None

---

## Table Components

### Table

**Original**: `.table`
```css
.table {
  width: 100%;
  border-collapse: collapse;
}
```

**Tailwind Equivalent**:
```html
<table class="w-full border-collapse">
```

**Custom CSS Needed**: None

---

### Table Header

**Original**: `.table thead th`
```css
.table thead th {
  background-color: #f9fafb;
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.875rem;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}
```

**Tailwind Equivalent**:
```html
<th class="bg-gray-50 px-4 py-3 text-left font-semibold text-sm text-gray-700 border-b border-gray-200">
```

**Custom CSS Needed**: None

---

### Table Cell

**Original**: `.table tbody td`
```css
.table tbody td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #e5e7eb;
  font-size: 0.875rem;
}
```

**Tailwind Equivalent**:
```html
<td class="px-4 py-3 border-b border-gray-200 text-sm">
```

**Custom CSS Needed**: None

---

### Table Row Hover

**Original**: `.table tbody tr:hover`
```css
.table tbody tr:hover {
  background-color: #f9fafb;
}
```

**Tailwind Equivalent**:
```html
<tr class="hover:bg-gray-50">
```

**Custom CSS Needed**: None

---

## Badge Components

### Badge

**Original**: `.badge`
```css
.badge {
  display: inline-flex;
  padding: 0.25rem 0.625rem;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 9999px;
}
```

**Tailwind Equivalent**:
```html
<span class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full">
```

**Custom CSS Needed**: None

---

### Badge Success

**Original**: `.badge--success`
```css
.badge--success {
  background-color: #d1fae5;
  color: #065f46;
}
```

**Tailwind Equivalent**:
```html
<span class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
```

**Custom CSS Needed**: None

---

### Badge Danger

**Original**: `.badge--danger`
```css
.badge--danger {
  background-color: #fee2e2;
  color: #991b1b;
}
```

**Tailwind Equivalent**:
```html
<span class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800">
```

**Custom CSS Needed**: None

---

### Badge Warning

**Original**: `.badge--warning`
```css
.badge--warning {
  background-color: #fef3c7;
  color: #92400e;
}
```

**Tailwind Equivalent**:
```html
<span class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800">
```

**Custom CSS Needed**: None

---

## Complex Components Requiring Custom CSS

These components have animations, complex selectors, or positioning that may need custom CSS in `@layer components` or separate files.

### Modal

**Custom CSS Needed**: Yes
- Backdrop overlay
- Enter/exit animations
- Focus trap styles

**Recommendation**: Use `@layer components` in `input.css`:

```css
@layer components {
  .modal-backdrop {
    @apply fixed inset-0 bg-black bg-opacity-50 z-40;
  }

  .modal-container {
    @apply fixed inset-0 z-50 flex items-center justify-center p-4;
  }

  .modal-content {
    @apply bg-white rounded-lg shadow-xl max-w-md w-full;
  }

  /* Animations */
  @keyframes modal-enter {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  .modal-enter {
    animation: modal-enter 0.2s ease-out;
  }
}
```

---

### Popover

**Custom CSS Needed**: Yes
- Absolute positioning with arrow
- z-index layering
- Show/hide transitions

**Recommendation**: Consider using Alpine.js `x-transition` directives with Tailwind classes, or add custom CSS for arrow positioning.

---

### Skeleton

**Custom CSS Needed**: Yes
- Loading animation shimmer effect

**Recommendation**: Use `@layer components`:

```css
@layer components {
  .skeleton {
    @apply bg-gray-200 rounded animate-pulse;
  }

  @keyframes shimmer {
    0% {
      background-position: -1000px 0;
    }
    100% {
      background-position: 1000px 0;
    }
  }

  .skeleton-shimmer {
    background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
    background-size: 1000px 100%;
    animation: shimmer 2s infinite;
  }
}
```

---

### Menu/Dropdown

**Custom CSS Needed**: Maybe
- Positioning (can use Tailwind `absolute`, `top-full`, etc.)
- Show/hide transitions (can use Alpine.js `x-transition`)

**Recommendation**: Try Tailwind + Alpine.js first. If too complex, add custom CSS.

---

## Utility Pattern Extraction

For frequently repeated utility combinations, consider extracting to `@layer components`:

### Example: Primary Input

If this pattern appears 20+ times:
```html
<input class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
```

Extract to component class:
```css
@layer components {
  .input-primary {
    @apply w-full px-3 py-2 border border-gray-300 rounded-md;
    @apply focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
  }
}
```

Then use:
```html
<input class="input-primary">
```

**When to extract**:
- Pattern repeats 10+ times
- Pattern is complex (5+ utility classes)
- Pattern has a semantic name

**When NOT to extract**:
- Pattern varies slightly per use
- Only used 1-3 times
- Utilities are simple (1-2 classes)

---

## Color Mapping

Extract colors from `variables.css` to Tailwind config:

```js
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#3b82f6', // blue-600
          light: '#60a5fa',    // blue-400
          dark: '#2563eb',     // blue-700
        },
        secondary: {
          DEFAULT: '#64748b', // slate-500
          light: '#94a3b8',   // slate-400
          dark: '#475569',    // slate-600
        },
        // Add other semantic colors from variables.css
      },
    },
  },
}
```

Then use as:
```html
<button class="bg-primary hover:bg-primary-dark text-white">
```

---

## Summary

- **90% of components** can be fully replaced with Tailwind utilities
- **Complex components** (modal, popover, skeleton) need custom CSS via `@layer components`
- **Repeated patterns** (10+ uses) should be extracted to component classes
- **Use @tailwindcss/forms plugin** to automatically style form inputs

**Next Steps**:
1. Start with simple pages (login, register) using this mapping
2. Identify any unmapped components during migration
3. Add custom CSS to `input.css` for complex components as needed
