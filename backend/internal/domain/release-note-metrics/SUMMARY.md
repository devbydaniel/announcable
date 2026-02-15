# Release Note Metrics

View and engagement tracking for release notes.

The `release-note-metrics` package tracks widget user interactions with release notes. Like the likes module, metrics use anonymous `ClientID` for identification.

**Key entity: `ReleaseNoteMetric`**
- `ReleaseNoteID` — The tracked release note
- `OrganisationID` — Tenant scoping
- `ClientID` — Anonymous client identifier
- `MetricType` — Enum: `view`, `cta_click`

**Key components:**
- `Service` — Record metrics, aggregate counts
- `Repository` — GORM queries for metric operations

**Integrations:**
- References `release-notes.ReleaseNote` and `organisation.Organisation`
- Widget API (`api/widget`) exposes metric creation endpoint
- Widget `tasks/release-note-metrics.ts` sends metrics on view/click
- `MetricType` is a Postgres enum type (`release_note_metric_type`)
