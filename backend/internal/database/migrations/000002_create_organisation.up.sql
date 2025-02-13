CREATE TABLE IF NOT EXISTS organisations (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL UNIQUE,
	external_id UUID NOT NULL,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS organisation_users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organisation_id UUID NOT NULL,
	user_id UUID NOT NULL,
	role VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_organisation_user
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,
	CONSTRAINT fk_user_organisation
	FOREIGN KEY (user_id) REFERENCES users(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS organisation_invites (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organisation_id UUID NOT NULL,
	email VARCHAR(255) NOT NULL,
	role VARCHAR(255) NOT NULL,
	external_id VARCHAR(255) NOT NULL,
	expires_at BIGINT NOT NULL,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_organisation_invite
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_organisation_users_organisation_id on organisation_users(organisation_id);
CREATE INDEX IF NOT EXISTS idx_organisation_users_user_id on organisation_users(user_id);
CREATE INDEX IF NOT EXISTS idx_organisation_invites_external_id on organisation_invites(external_id);
CREATE INDEX IF NOT EXISTS idx_organisation_invites_organiation_id on organisation_invites(organisation_id);
