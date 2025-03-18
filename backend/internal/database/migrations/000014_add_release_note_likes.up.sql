-- Add like feature configuration to widget_configs
ALTER TABLE widget_configs
ADD COLUMN enable_likes BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN like_button_text VARCHAR(255) NOT NULL DEFAULT 'Like',
ADD COLUMN unlike_button_text VARCHAR(255) NOT NULL DEFAULT 'Unlike';

-- Create likes table
CREATE TABLE release_note_likes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  release_note_id UUID NOT NULL REFERENCES release_notes(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL,
  organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX release_note_likes_release_note_id_idx ON release_note_likes(release_note_id);
CREATE INDEX release_note_likes_organisation_id_idx ON release_note_likes(organisation_id);
CREATE UNIQUE INDEX release_note_likes_unique_client_idx ON release_note_likes(release_note_id, client_id) WHERE deleted_at IS NULL; 