# Release Page Configuration

Settings for the public-facing release notes page.

The `release-page-configs` package manages per-organisation configuration for the standalone release notes page (served at `/s/{orgSlug}`). Each organisation has one release page config controlling branding and layout.

**Key entity: `ReleasePageConfig`**
- Organisation-scoped (`OrganisationID`)
- Branding: `Title`, `Description`, `ImagePath`, `ImageUrl` (transient signed URL)
- Colors: `BgColor`, `TextColor`, `TextColorMuted`
- Layout: `BrandPosition` — `top` or `left`
- Navigation: `BackLinkLabel`, `BackLinkUrl`
- Routing: `Slug` (URL path segment for `/s/{slug}`)
- Toggle: `DisableReleasePage`

**Key components:**
- `Service` — Get and update config, with image upload handling
- `Repository` — GORM queries
- `ImageInput` — Image upload data with delete flag

**Integrations:**
- Public release page handler (`pages/public/release_page`) reads this config
- Admin UI at `/release-page-config` allows editing
- `objstore` for logo/brand image storage
- `Slug` is used in the public URL: `https://announcable.com/s/{slug}`

**Notes:**
- `ImageUrl` is transient (`gorm:"-"`) — populated with signed URL at query time
- `DisableReleasePage` allows orgs to hide their public page entirely
