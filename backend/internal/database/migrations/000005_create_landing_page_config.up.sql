CREATE TABLE IF NOT EXISTS lp_configs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organisation_id UUID NOT NULL,
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	bg_color VARCHAR(255) DEFAULT '#FFFFFF',
	text_color VARCHAR(255) DEFAULT '#000000',
	text_color_muted VARCHAR(255) DEFAULT '#000000',
	brand_position VARCHAR(255) DEFAULT 'top',
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,

	CONSTRAINT fk_org_lp_configs
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE

);

CREATE INDEX IF NOT EXISTS idx_lp_configs_org_id on lp_configs(organisation_id);
