-- Recreate tos_confirms table
CREATE TABLE IF NOT EXISTS tos_confirms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    version VARCHAR(255) NOT NULL,
    confirmed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_user_tos
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tos_confirm_user_id ON tos_confirms(user_id);
CREATE INDEX IF NOT EXISTS idx_tos_confirm_version ON tos_confirms(version);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tos_confirm_user_version ON tos_confirms(user_id, version);

-- Recreate privacy_policy_confirms table
CREATE TABLE IF NOT EXISTS privacy_policy_confirms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    version VARCHAR(255) NOT NULL,
    confirmed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_user_privacy_policy
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_privacy_policy_confirm_user_id ON privacy_policy_confirms(user_id);
CREATE INDEX IF NOT EXISTS idx_privacy_policy_confirm_version ON privacy_policy_confirms(version);
CREATE UNIQUE INDEX IF NOT EXISTS idx_privacy_policy_confirm_user_version ON privacy_policy_confirms(user_id, version);
