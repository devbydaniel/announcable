# Release Notes

Core release note management — CRUD, publishing, and image handling.

The `release-notes` package is the central domain module. Release notes are scoped to an organisation and support rich content with images, media links, CTAs, and attention mechanisms.

**Key entity: `ReleaseNote`**
- Organisation-scoped (`OrganisationID`)
- Content fields: `Title`, `DescriptionShort`, `DescriptionLong`, `ReleaseDate`
- Media: `ImagePath` (object storage), `ImageUrl` (signed URL, transient), `MediaLink`
- CTA: `CtaLabelOverride`, `CtaUrlOverride`, `HideCta`
- Publishing: `IsPublished` flag
- Visibility: `HideOnWidget`, `HideOnReleasePage`
- Attention: `AttentionMechanism` (indicator or instant open)
- Audit: `CreatedBy`, `LastUpdatedBy` (user UUIDs)

**Supporting types:**
- `PaginatedReleaseNotes` — Paginated list response
- `ReleaseNoteStatus` — Lightweight status for widget polling
- `AttentionMechanism` — Enum: `show_indicator`, `instant_open`
- `ImageInput` — Image upload data with delete flag

**Key components:**
- `Service` — CRUD operations with transactional image handling via object storage
- `Repository` — GORM queries with pagination, filtering by org, published status
- Image processing uses `imgUtil.ImgProcessConfig` (max width 1000px, quality 80)

**Integrations:**
- `organisation.Organisation` for tenant scoping
- `objstore` for image storage (Minio)
- `imgUtil` for image resizing/compression
- Referenced by `release-note-likes` and `release-note-metrics` modules
- Served to widget via `api/widget` handlers

**Notes:**
- `ImageUrl` is a transient field (`gorm:"-"`) — populated at query time with signed URLs
- Image path format: `{orgId}/{randomId}.{format}`
- Create and update operations use transactions to ensure image + record consistency
