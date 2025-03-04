CREATE TYPE release_note_metric_type AS ENUM ('view', 'cta_click');

CREATE TABLE release_note_metrics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  release_note_id UUID NOT NULL REFERENCES release_notes(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  metric_type release_note_metric_type NOT NULL,
  organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX release_note_metrics_release_note_id_idx ON release_note_metrics(release_note_id);
CREATE INDEX release_note_metrics_organisation_id_idx ON release_note_metrics(organisation_id);