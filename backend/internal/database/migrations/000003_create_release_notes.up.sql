CREATE TABLE IF NOT EXISTS release_notes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organisation_id UUID NOT NULL,
	title VARCHAR(255) NOT NULL,
	description_short TEXT NOT NULL,
	description_long TEXT,
	release_date DATE,
	is_published BOOLEAN DEFAULT FALSE,
	cta_label_override VARCHAR(255),
	cta_url_override VARCHAR(255),
	hide_cta BOOLEAN DEFAULT FALSE,
	attention_mechanism VARCHAR(255),
	created_at TIMESTAMPTZ,
	created_by UUID,
	updated_at TIMESTAMPTZ,
	last_updated_by UUID,
	deleted_at TIMESTAMPTZ,

	CONSTRAINT fk_org_release_notes
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,

	CONSTRAINT fk_user_created_by
	FOREIGN KEY (created_by) REFERENCES users(id)
	ON DELETE SET NULL
	ON UPDATE CASCADE,

	CONSTRAINT fk_user_last_updated_by
	FOREIGN KEY (last_updated_by) REFERENCES users(id)
	ON DELETE SET NULL
	ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_release_notes_org_id on release_notes(organisation_id);
