CREATE TABLE IF NOT EXISTS widget_configs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organisation_id UUID NOT NULL,
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	widget_border_radius INT DEFAULT 0,
	widget_border_color VARCHAR(255) DEFAULT '#000000',
	widget_border_width INT DEFAULT 0,
	widget_bg_color VARCHAR(255) DEFAULT '#FFFFFF',
	widget_text_color VARCHAR(255) DEFAULT '#000000',
	widget_type VARCHAR(255) NOT NULL,
	release_note_border_color VARCHAR(255) DEFAULT '#000000',
	release_note_border_radius INT DEFAULT 0,
	release_note_border_width INT DEFAULT 0,
	release_note_bg_color VARCHAR(255) DEFAULT '#FFFFFF',
	release_note_text_color VARCHAR(255) DEFAULT '#000000',
	release_note_cta_text VARCHAR(255) DEFAULT 'Learn More',
	release_page_base_url VARCHAR(255),
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,

	CONSTRAINT fk_org_widget_configs
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_widget_configs_org_id on widget_configs(organisation_id);
