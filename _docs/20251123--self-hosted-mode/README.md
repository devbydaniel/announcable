# Self-Hosted Mode Implementation

## Overview
This document describes the self-hosted mode feature that allows the application to run without Stripe integration and subscription limitations.

## Configuration

Set `APP_ENVIRONMENT=self-hosted` in your `.env` file, or leave it unset (defaults to self-hosted mode).

## What Changes in Self-Hosted Mode

### Backend
1. **Subscription Checks**: All bypassed - `HasActiveSubscription` always returns `true`
2. **Release Note Limits**: No 5-note limit for free tier
3. **Stripe Integration**: Not initialized, API keys not required
4. **Payment Routes**: Return 404 - not accessible
5. **Subscription Middleware**: Short-circuits to always allow access

### Frontend/Templates
1. **Subscribe Button**: Hidden from navigation
2. **Subscription Management**: Hidden from settings page
3. **Subscription Indicator**: Not shown in UI

## Modified Files

### Configuration
- `backend/config/config.go` - Added AppEnvironment field and helper methods

### Middleware
- `backend/internal/middleware/subscription.go` - Bypass logic for self-hosted
- `backend/internal/middleware/environment.go` - NEW: CloudOnly middleware

### Handlers
- `backend/internal/handler/pages/release_notes/create/release_note_create_action.go` - Conditional limit check
- All page handlers - Added ShowSubscriptionUI flag

### Templates
- `backend/templates/partials/nav.html` - Conditional subscribe button
- `backend/templates/pages/settings-page.html` - Conditional subscription section

### Main
- `backend/main.go` - Conditional Stripe initialization and route protection

## Testing Self-Hosted Mode

1. Set `APP_ENVIRONMENT=self-hosted` in `.env`
2. Remove or leave empty Stripe environment variables
3. Start application
4. Verify:
   - No "Subscribe" button in navigation
   - Can create unlimited release notes
   - Settings page has no subscription section
   - `/payment/*` routes return 404
   - `/subscription/*` routes return 404

## Testing Cloud Mode

1. Set `APP_ENVIRONMENT=cloud` in `.env`
2. Configure Stripe API keys
3. Start application
4. Verify:
   - "Subscribe" button shows when no active subscription
   - Free tier limited to 5 release notes
   - Settings page shows subscription management
   - Payment routes accessible
   - Subscription webhook processing works
