ALTER TABLE release_notes
ADD COLUMN image_path VARCHAR(1024);

ALTER TABLE release_page_configs
ADD COLUMN image_path VARCHAR(1024);
