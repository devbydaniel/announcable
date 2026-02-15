# Widget Configuration

Appearance and behavior settings for the embeddable widget.

The `widget-configs` package manages per-organisation configuration for the embeddable widget. Each organisation has one widget config controlling visual appearance and behavior.

**Key entity: `WidgetConfig`**
- Organisation-scoped (`OrganisationID`)
- Widget appearance: `Title`, `Description`, `WidgetBorderRadius/Color/Width`, `WidgetBgColor`, `WidgetTextColor`
- Widget type: `WidgetType` — `modal`, `popover`, or `sidebar`
- Release note card appearance: `ReleaseNoteBorderRadius/Color/Width`, `ReleaseNoteBgColor`, `ReleaseNoteTextColor`
- CTA: `ReleaseNoteCtaText` (default CTA label)
- Release page: `ReleasePageBaseUrl` (link to full release page)
- Likes: `EnableLikes`, `LikeButtonText`, `UnlikeButtonText`

**Key components:**
- `Service` — Get and update config for an organisation
- `Repository` — GORM queries

**Integrations:**
- Widget fetches config via `GET /api/widget-config/{orgId}`
- Widget `tasks/widget-config.ts` consumes this on the client side
- Admin UI at `/widget-config` allows editing via `pages/widget/config` handler
- Default values set via GORM tags
