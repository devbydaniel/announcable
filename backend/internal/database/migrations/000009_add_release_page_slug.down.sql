DROP INDEX IF EXISTS idx_release_page_configs_slug;

ALTER TABLE release_page_configs
DROP COLUMN slug;
