ALTER TABLE release_page_configs
ADD COLUMN slug VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_release_page_configs_slug on release_page_configs(slug);
