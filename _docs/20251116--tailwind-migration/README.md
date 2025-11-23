# Tailwind CSS Migration Documentation

This directory contains all documentation for migrating the Announcable backend from manual CSS management to Tailwind CSS.

## Status

**Planning** - Migration not yet started

## Documents

### 1. [Migration Plan](./migration-plan.md)

**Purpose**: Complete step-by-step migration strategy

**Contents**:
- Current state analysis and problems
- Target architecture
- 7-phase migration plan with detailed steps
- Risk mitigation and rollback strategy
- Timeline estimate (2.5 days)

**Start here** to understand the overall migration approach.

---

### 2. [Component Mappings](./component-mappings.md)

**Purpose**: Reference guide mapping existing CSS to Tailwind equivalents

**Contents**:
- Side-by-side comparison of current CSS and Tailwind utilities
- Mappings for all component types (forms, buttons, cards, nav, tables, badges, alerts)
- Complex components requiring custom CSS
- Utility pattern extraction guidelines

**Use this** during template migration to quickly find Tailwind equivalents.

---

### 3. [Tailwind Development Guide](./tailwind-guide.md)

**Purpose**: Developer reference for working with Tailwind post-migration

**Contents**:
- Common layout, form, button, card, and table patterns
- HTMX and Alpine.js integration examples
- Responsive design patterns
- Performance tips and common gotchas
- When to use custom CSS vs utilities

**Use this** after migration for ongoing development.

---

## Quick Links

- **Want to start migration?** → Read [migration-plan.md](./migration-plan.md) Phase 1
- **Migrating a template?** → Reference [component-mappings.md](./component-mappings.md)
- **Adding new features?** → Follow [tailwind-guide.md](./tailwind-guide.md) patterns
- **Stuck on a component?** → Check [component-mappings.md](./component-mappings.md) complex components section

## Migration Phases Overview

1. **Setup Tailwind Build Process** - Install dependencies, configure build (2 hours)
2. **Create Component Mappings** - Document CSS → Tailwind equivalents (3 hours)
3. **Migrate Templates Page-by-Page** - Convert simple → complex pages (4 hours)
4. **Migrate Shared Components** - Nav, header, layouts (3 hours)
5. **Handle Complex Components** - Modals, popovers, animations (4 hours)
6. **Clean Up and Optimize** - Delete legacy CSS, optimize build (2 hours)
7. **Documentation** - Update CLAUDE.md, training docs (2 hours)

**Total Estimated Effort**: 20 hours (2.5 days)

## Key Benefits

After migration:
- ✅ Single CSS bundle (reduced from 15+ files)
- ✅ ~75% smaller CSS size (purging unused styles)
- ✅ Consistent design system (enforced spacing/colors)
- ✅ Faster development (no manual CSS file management)
- ✅ Better maintainability (utility-first approach)

## Questions?

Review the migration plan's FAQ section or update these docs as you discover new patterns during migration.
