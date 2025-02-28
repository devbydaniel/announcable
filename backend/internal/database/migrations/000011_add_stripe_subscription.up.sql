CREATE TABLE stripe_subscriptions (
    id UUID PRIMARY KEY,
    organisation_id UUID REFERENCES organisations(id),
    stripe_subscription_id TEXT NOT NULL,
    stripe_customer_id TEXT NOT NULL,
    status TEXT NOT NULL,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
	deleted_at TIMESTAMP,
    CONSTRAINT fk_stripe_subscription_organisation
	FOREIGN KEY (organisation_id) REFERENCES organisations(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE,
);
