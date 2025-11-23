# Tailwind CSS + Basecoat UI Usage Guide

This guide covers how to use Tailwind CSS and Basecoat UI in your Go html/template files.

## Table of Contents

- [Overview](#overview)
- [Development Workflow](#development-workflow)
- [Tailwind CSS Basics](#tailwind-css-basics)
- [Basecoat UI Components](#basecoat-ui-components)
- [Common Patterns](#common-patterns)
- [Dark Mode](#dark-mode)
- [Icons with Lucide](#icons-with-lucide)
- [Toast Notifications](#toast-notifications)
- [Customization](#customization)

---

## Overview

**Tailwind CSS** is a utility-first CSS framework. Instead of writing custom CSS, you compose designs using pre-built utility classes.

**Basecoat UI** is a component library built on top of Tailwind. It provides pre-styled components (buttons, cards, forms, etc.) using Tailwind utilities under the hood.

### When to Use What

- **Tailwind utilities**: For layout, spacing, colors, sizing, etc.
- **Basecoat components**: For common UI elements (buttons, inputs, cards, etc.)
- **Mix both**: Use Basecoat for the base component, Tailwind to customize

Example:

```html
<button class="btn w-full">
  <!-- 'btn' is Basecoat, 'w-full' is Tailwind -->
  Sign In
</button>
```

---

## Development Workflow

### Starting Development

Open **three terminals**:

```bash
# Terminal 1: Start service containers
cd backend
make dev-services

# Terminal 2: Start Go app with Air (hot-reload)
cd backend
air

# Terminal 3: Start Vite (CSS/JS watch mode)
cd backend
npm run dev
```

### What Happens

- **Vite** watches `assets/css/main.css` and `assets/js/main.js`
- Changes automatically rebuild to `static/dist/styles.css` and `static/dist/main.js`
- **Air** watches Go files and templates, reloads the server
- Changes to templates with new Tailwind classes trigger Vite rebuild

### Stopping Development

```bash
# Ctrl+C in Terminal 2 (Air)
# Ctrl+C in Terminal 3 (Vite)
# Terminal 1:
make dev-stop
```

---

## Tailwind CSS Basics

Tailwind provides utility classes for almost everything. Here are the most common categories:

### Layout

```html
<!-- Flexbox -->
<div class="flex items-center justify-between gap-4">
  <!-- flex: display flex -->
  <!-- items-center: align-items center -->
  <!-- justify-between: justify-content space-between -->
  <!-- gap-4: gap 1rem -->
</div>

<!-- Grid -->
<div class="grid grid-cols-3 gap-6">
  <!-- 3 column grid with 1.5rem gap -->
</div>

<!-- Container -->
<div class="container mx-auto px-4">
  <!-- Centered container with horizontal padding -->
</div>
```

### Spacing

Pattern: `{property}{side}-{size}`

- **Properties**: `m` (margin), `p` (padding)
- **Sides**: `t` (top), `r` (right), `b` (bottom), `l` (left), `x` (horizontal), `y` (vertical), or none for all sides
- **Sizes**: `0` to `96` (0 to 24rem), plus `px` for 1px

```html
<div class="p-4"><!-- padding: 1rem all sides --></div>
<div class="mt-8"><!-- margin-top: 2rem --></div>
<div class="mx-auto"><!-- margin: 0 auto (center) --></div>
<div class="space-y-4"><!-- vertical spacing between children --></div>
```

### Sizing

```html
<!-- Width -->
<div class="w-full"><!-- width: 100% --></div>
<div class="w-64"><!-- width: 16rem --></div>
<div class="w-1/2"><!-- width: 50% --></div>

<!-- Height -->
<div class="h-screen"><!-- height: 100vh --></div>
<div class="h-64"><!-- height: 16rem --></div>

<!-- Max/Min -->
<div class="max-w-4xl"><!-- max-width: 56rem --></div>
<div class="min-h-screen"><!-- min-height: 100vh --></div>
```

### Typography

```html
<!-- Text size -->
<p class="text-sm">Small text</p>
<p class="text-base">Base text (16px)</p>
<p class="text-lg">Large text</p>
<h1 class="text-3xl">Large heading</h1>

<!-- Font weight -->
<span class="font-normal">Normal</span>
<span class="font-medium">Medium</span>
<span class="font-semibold">Semi-bold</span>
<span class="font-bold">Bold</span>

<!-- Text alignment -->
<p class="text-left">Left</p>
<p class="text-center">Center</p>
<p class="text-right">Right</p>

<!-- Text color -->
<p class="text-gray-700">Gray text</p>
<p class="text-red-600">Red text</p>
<p class="text-blue-500">Blue text</p>
```

### Colors

Basecoat defines semantic color names that work with dark mode:

```html
<!-- Background -->
<div class="bg-background"><!-- Base background --></div>
<div class="bg-card"><!-- Card background --></div>
<div class="bg-primary"><!-- Primary brand color --></div>
<div class="bg-secondary"><!-- Secondary color --></div>

<!-- Text -->
<p class="text-foreground"><!-- Base text --></p>
<p class="text-muted-foreground"><!-- Muted text --></p>
<p class="text-primary"><!-- Primary text --></p>

<!-- Borders -->
<div class="border border-border"><!-- Standard border --></div>
```

### Borders & Rounded Corners

```html
<!-- Borders -->
<div class="border"><!-- border: 1px solid --></div>
<div class="border-2"><!-- border: 2px solid --></div>
<div class="border-t"><!-- border-top only --></div>

<!-- Rounded corners -->
<div class="rounded"><!-- border-radius: 0.25rem --></div>
<div class="rounded-md"><!-- border-radius: 0.375rem --></div>
<div class="rounded-lg"><!-- border-radius: 0.5rem --></div>
<div class="rounded-full"><!-- border-radius: 9999px (circle) --></div>
```

### Shadows

```html
<div class="shadow-sm">Small shadow</div>
<div class="shadow">Default shadow</div>
<div class="shadow-md">Medium shadow</div>
<div class="shadow-lg">Large shadow</div>
```

### Responsive Design

Add breakpoint prefixes to make utilities responsive:

- `sm:` - 640px and up
- `md:` - 768px and up
- `lg:` - 1024px and up
- `xl:` - 1280px and up

```html
<div class="w-full md:w-1/2 lg:w-1/3">
  <!-- Full width on mobile, half on tablet, third on desktop -->
</div>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
  <!-- Responsive grid -->
</div>
```

### Hover, Focus, and Other States

```html
<button class="bg-blue-500 hover:bg-blue-600 focus:ring-2 focus:ring-blue-300">
  Hover and focus states
</button>

<div class="opacity-50 hover:opacity-100 transition-opacity">
  Smooth hover transition
</div>
```

---

## Basecoat UI Components

Basecoat provides ready-to-use component classes. Here are the most common:

### Buttons

```html
<!-- Primary button (default) -->
<button class="btn">Primary Button</button>
<button class="btn-primary">Primary Button</button>

<!-- Secondary button -->
<button class="btn-secondary">Secondary</button>

<!-- Outline button -->
<button class="btn-outline">Outline</button>

<!-- Ghost button -->
<button class="btn-ghost">Ghost</button>

<!-- Destructive button -->
<button class="btn-destructive">Delete</button>

<!-- Sizes -->
<button class="btn btn-sm">Small</button>
<button class="btn">Default</button>
<button class="btn btn-lg">Large</button>

<!-- Full width with Tailwind -->
<button class="btn w-full">Full Width</button>

<!-- With icon -->
<button class="btn">
  <svg>...</svg>
  Button with Icon
</button>
```

### Forms

#### Input Fields

```html
<!-- Text input -->
<div class="space-y-2">
  <label class="label" for="email">Email</label>
  <input class="input" type="email" id="email" placeholder="you@example.com" />
</div>

<!-- Input with error -->
<div class="space-y-2">
  <label class="label" for="password">Password</label>
  <input class="input input-error" type="password" id="password" />
  <p class="text-sm text-destructive">Password is required</p>
</div>

<!-- Disabled input -->
<input class="input" type="text" disabled value="Disabled" />
```

#### Textarea

```html
<div class="space-y-2">
  <label class="label" for="message">Message</label>
  <textarea class="textarea" id="message" rows="4"></textarea>
</div>
```

#### Select

```html
<div class="space-y-2">
  <label class="label" for="country">Country</label>
  <select class="select" id="country">
    <option>United States</option>
    <option>Canada</option>
    <option>United Kingdom</option>
  </select>
</div>
```

#### Checkbox

```html
<div class="flex items-center gap-2">
  <input class="checkbox" type="checkbox" id="terms" />
  <label class="label" for="terms">I agree to the terms</label>
</div>
```

#### Radio

```html
<div class="space-y-2">
  <div class="flex items-center gap-2">
    <input class="radio" type="radio" id="option1" name="option" />
    <label class="label" for="option1">Option 1</label>
  </div>
  <div class="flex items-center gap-2">
    <input class="radio" type="radio" id="option2" name="option" />
    <label class="label" for="option2">Option 2</label>
  </div>
</div>
```

#### Switch

```html
<div class="flex items-center gap-2">
  <input class="switch" type="checkbox" id="notifications" />
  <label class="label" for="notifications">Enable notifications</label>
</div>
```

### Cards

```html
<!-- Basic card -->
<div class="card">
  <header>
    <h3>Card Title</h3>
    <p class="text-sm text-muted-foreground">Card description</p>
  </header>
  <section>
    <p>Card content goes here.</p>
  </section>
  <footer>
    <button class="btn">Action</button>
  </footer>
</div>

<!-- Card with no sections -->
<div class="card p-6">
  <h3 class="text-lg font-semibold mb-2">Simple Card</h3>
  <p>Just add padding and structure with Tailwind.</p>
</div>
```

### Alerts

```html
<!-- Info alert (default) -->
<div class="alert">
  <svg class="alert-icon">...</svg>
  <div>
    <p class="alert-title">Information</p>
    <p class="alert-description">This is an informational message.</p>
  </div>
</div>

<!-- Success alert -->
<div class="alert alert-success">
  <svg class="alert-icon">...</svg>
  <div>
    <p class="alert-title">Success</p>
    <p class="alert-description">Your changes have been saved.</p>
  </div>
</div>

<!-- Warning alert -->
<div class="alert alert-warning">...</div>

<!-- Error alert -->
<div class="alert alert-destructive">...</div>
```

### Badges

```html
<span class="badge">Default</span>
<span class="badge badge-secondary">Secondary</span>
<span class="badge badge-outline">Outline</span>
<span class="badge badge-destructive">Error</span>
```

### Tables

```html
<table class="table">
  <thead>
    <tr>
      <th>Name</th>
      <th>Email</th>
      <th>Role</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>John Doe</td>
      <td>john@example.com</td>
      <td>Admin</td>
    </tr>
    <tr>
      <td>Jane Smith</td>
      <td>jane@example.com</td>
      <td>User</td>
    </tr>
  </tbody>
</table>
```

### Tabs

```html
<div x-data="{ tab: 'overview' }">
  <div class="tabs">
    <button
      class="tab"
      :class="{ 'active': tab === 'overview' }"
      @click="tab = 'overview'"
    >
      Overview
    </button>
    <button
      class="tab"
      :class="{ 'active': tab === 'details' }"
      @click="tab = 'details'"
    >
      Details
    </button>
  </div>

  <div x-show="tab === 'overview'" class="mt-4">Overview content</div>
  <div x-show="tab === 'details'" class="mt-4">Details content</div>
</div>
```

---

## Common Patterns

### Form Layout

```html
<form class="space-y-6" hx-post="/submit">
  <div class="grid gap-4">
    <div class="space-y-2">
      <label class="label" for="name">Name</label>
      <input class="input" type="text" id="name" name="name" />
    </div>

    <div class="space-y-2">
      <label class="label" for="email">Email</label>
      <input class="input" type="email" id="email" name="email" />
    </div>
  </div>

  <div class="flex justify-end gap-3">
    <button type="button" class="btn-secondary">Cancel</button>
    <button type="submit" class="btn">Submit</button>
  </div>
</form>
```

### Two-Column Layout

```html
<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <div class="card">
    <header>
      <h3>Left Column</h3>
    </header>
    <section>Content</section>
  </div>

  <div class="card">
    <header>
      <h3>Right Column</h3>
    </header>
    <section>Content</section>
  </div>
</div>
```

### List with Actions

```html
<div class="divide-y divide-border">
  <div class="flex items-center justify-between py-4">
    <div>
      <h4 class="font-medium">Item Title</h4>
      <p class="text-sm text-muted-foreground">Item description</p>
    </div>
    <div class="flex gap-2">
      <button class="btn btn-sm btn-secondary">Edit</button>
      <button class="btn btn-sm btn-destructive">Delete</button>
    </div>
  </div>
  <!-- More items... -->
</div>
```

### Modal/Dialog Layout

```html
<div class="card max-w-lg mx-auto">
  <header>
    <h2 class="text-xl font-semibold">Confirm Action</h2>
    <p class="text-sm text-muted-foreground">
      Are you sure you want to continue?
    </p>
  </header>
  <section>
    <p>This action cannot be undone.</p>
  </section>
  <footer class="flex justify-end gap-3">
    <button class="btn-secondary" onclick="closeModal()">Cancel</button>
    <button class="btn-destructive">Confirm</button>
  </footer>
</div>
```

### Empty State

```html
<div class="flex flex-col items-center justify-center py-12 text-center">
  <svg class="w-16 h-16 text-muted-foreground mb-4">
    <!-- Icon -->
  </svg>
  <h3 class="text-lg font-semibold mb-2">No items found</h3>
  <p class="text-muted-foreground mb-6">
    Get started by creating your first item.
  </p>
  <button class="btn">Create Item</button>
</div>
```

---

## Dark Mode

Basecoat includes built-in dark mode support via the `.dark` class on the `<html>` element.

### Enabling Dark Mode

```html
<!-- Add to root.html or base layout -->
<html class="dark">
  <!-- Dark mode is now active -->
</html>
```

### Toggle Implementation

```html
<div
  x-data="{ darkMode: false }"
  x-init="darkMode = localStorage.getItem('darkMode') === 'true'"
>
  <button
    @click="
      darkMode = !darkMode;
      localStorage.setItem('darkMode', darkMode);
      document.documentElement.classList.toggle('dark', darkMode);
    "
    class="btn-secondary"
  >
    <span x-show="!darkMode">🌙 Dark</span>
    <span x-show="darkMode">☀️ Light</span>
  </button>
</div>
```

### Dark Mode Utilities

Use Tailwind's `dark:` prefix for custom dark mode styles:

```html
<div class="bg-white dark:bg-gray-900 text-black dark:text-white">
  This adapts to dark mode
</div>
```

---

## Icons with Lucide

Lucide icons are included via the npm package. You can use them in two ways:

### Inline SVG (Recommended)

Copy the SVG code from [lucide.dev](https://lucide.dev):

```html
<button class="btn">
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <path d="M5 12h14" />
    <path d="m12 5 7 7-7 7" />
  </svg>
  Next
</button>
```

### With Tailwind Size Classes

```html
<!-- Small icon -->
<svg class="w-4 h-4" ...>...</svg>

<!-- Default icon -->
<svg class="w-5 h-5" ...>...</svg>

<!-- Large icon -->
<svg class="w-6 h-6" ...>...</svg>
```

---

## Toast Notifications

Basecoat automatically integrates with HTMX for toast notifications.

### Backend (Go)

Use the `HX-Trigger` header with the `basecoat:toast` event:

```go
// Success toast
w.Header().Set("HX-Trigger", `{"basecoat:toast": {"text": "Settings saved successfully", "type": "success"}}`)

// Error toast
w.Header().Set("HX-Trigger", `{"basecoat:toast": {"text": "An error occurred", "type": "error"}}`)

// Warning toast
w.Header().Set("HX-Trigger", `{"basecoat:toast": {"text": "Please verify your email", "type": "warning"}}`)

// Info toast
w.Header().Set("HX-Trigger", `{"basecoat:toast": {"text": "New update available", "type": "info"}}`)
```

### Template Requirement

Add the toaster container to your base layout (already in templates/layouts/root.html):

```html
<body>
  <!-- Your content -->

  <div id="toaster" class="toaster"></div>

  <script type="module" src="/static/dist/main.js"></script>
</body>
```

### No JavaScript Required

Toasts automatically appear when the HTMX response includes the `HX-Trigger` header. No client-side JavaScript needed!

---

## Customization

### Custom Colors

Create `assets/css/brand.css` and import it in `main.css`:

```css
/* assets/css/brand.css */
@layer base {
  :root {
    --primary: 220 70% 50%; /* HSL values */
    --primary-foreground: 0 0% 100%;
    /* Override other colors as needed */
  }
}
```

```css
/* assets/css/main.css */
@import "tailwindcss";
@import "basecoat-css";
@import "./brand.css";
```

### Custom Components

Add custom component classes in a new CSS file:

```css
/* assets/css/custom-components.css */
@layer components {
  .my-custom-button {
    @apply btn bg-gradient-to-r from-purple-500 to-pink-500;
  }
}
```

Import in `main.css`:

```css
@import "tailwindcss";
@import "basecoat-css";
@import "./custom-components.css";
```

### Arbitrary Values

Use square brackets for one-off custom values:

```html
<div class="w-[347px] bg-[#1da1f2] top-[117px]">Custom values</div>
```

---

## Resources

- **Tailwind Docs**: https://tailwindcss.com/docs
- **Basecoat Docs**: https://basecoatui.com
- **Lucide Icons**: https://lucide.dev
- **Tailwind Cheat Sheet**: https://nerdcave.com/tailwind-cheat-sheet

---

## Tips

1. **Use semantic color names** (`bg-primary`, `text-foreground`) instead of hardcoded colors for automatic dark mode support
2. **Compose utilities** - Combine multiple small utilities instead of writing custom CSS
3. **Use `space-y-*`** for vertical spacing between children (easier than individual margins)
4. **Leverage responsive prefixes** - Design mobile-first, add breakpoints as needed
5. **Keep it simple** - Don't over-engineer, Basecoat + Tailwind handle most cases out of the box
