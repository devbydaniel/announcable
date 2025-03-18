-- Remove like feature configuration from widget_configs
ALTER TABLE widget_configs
DROP COLUMN enable_likes,
DROP COLUMN like_button_text,
DROP COLUMN unlike_button_text;

-- Drop likes table
DROP TABLE IF EXISTS release_note_likes; 