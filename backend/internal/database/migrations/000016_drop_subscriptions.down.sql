CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID REFERENCES organisations(id),
    stripe_subscription_id TEXT NOT NULL,
	is_active BOOLEAN NOT NULL,
	is_free BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_stripe_subscription_organisation
    FOREIGN KEY (organisation_id) REFERENCES organisations(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE INDEX idx_subscriptions_organisation_id ON subscriptions(organisation_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_stripe_id ON subscriptions(stripe_subscription_id);
