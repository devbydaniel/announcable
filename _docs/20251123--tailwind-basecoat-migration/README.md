# Tailwind + Basecoat Migration

**Date:** November 23, 2025
**Status:** Planning
**Based on:** mono2 repository setup

## Overview

This migration adds Tailwind CSS v4 and Basecoat UI to the mono backend, replacing the current custom CSS system. The approach is based on the proven setup in the mono2 repository.

## Key Decisions

- **Migration Strategy:** Incremental (page-by-page)
- **Design Philosophy:** Embrace Basecoat defaults with minor brand adjustments
- **Features:** Dark mode, Lucide icons, Basecoat toast system, bundled dependencies
- **Dev Workflow:** Air + npm run dev in separate terminals, Docker Compose only for services

## Quick Links

- [Migration Plan](./migration-plan.md) - Detailed step-by-step implementation plan
- [Technical Comparison](./technical-comparison.md) - Differences between mono and mono2 setups
- [Implementation Notes](./implementation-notes.md) - Technical details and code examples

## Quick Start (After Migration)

### Development Workflow

```bash
# Terminal 1: Start service containers (postgres, mail, minio, etc.)
cd backend
make dev-services

# Terminal 2: Start Air (Go hot-reload)
cd backend
air

# Terminal 3: Start Vite (CSS/JS watch)
cd backend
npm run dev
```

## Estimated Timeline

- **Phase 1-2** (Infrastructure): 2-3 hours
- **Phase 3-4** (Template migration): 8-12 hours
- **Phase 5-6** (Toast/Icon systems): 2-3 hours
- **Phase 7** (Dark mode): 2-3 hours
- **Phase 8** (Cleanup): 1-2 hours

**Total:** ~15-23 hours

## Current Status

- [ ] Phase 1: Infrastructure Setup
- [ ] Phase 2: Template Foundation
- [ ] Phase 3: Incremental Page Migration
- [ ] Phase 4: Component Migration
- [ ] Phase 5: Toast System Migration
- [ ] Phase 6: Icon System Migration
- [ ] Phase 7: Dark Mode Implementation
- [ ] Phase 8: Cleanup & Documentation
