package subscription

import (
	"github.com/devbydaniel/release-notes-go/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *repository {
	log.Trace().Msg("NewRepository")
	return &repository{db: db}
}

// Create creates a new subscription
func (r *repository) Create(s *Subscription, tx *gorm.DB) error {
	log.Trace().Msg("Create")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Create(s).Error; err != nil {
		log.Error().Err(err).Msg("Error creating subscription")
		return err
	}
	return nil
}

// Get retrieves a subscription by organisation ID
func (r *repository) Get(orgID uuid.UUID) (*Subscription, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Get")
	var subscription Subscription

	if err := r.db.Client.Where("organisation_id = ?", orgID).First(&subscription).Error; err != nil {
		log.Error().Err(err).Msg("Error finding subscription")
		return nil, err
	}
	return &subscription, nil
}

// GetByStripeSubscriptionID retrieves a subscription by Stripe subscription ID
func (r *repository) GetByStripeSubscriptionID(stripeSubID string) (*Subscription, error) {
	log.Trace().Str("stripeSubID", stripeSubID).Msg("GetByStripeSubscriptionID")
	var subscription Subscription

	if err := r.db.Client.Where("stripe_subscription_id = ?", stripeSubID).First(&subscription).Error; err != nil {
		log.Error().Err(err).Msg("Error finding subscription")
		return nil, err
	}
	return &subscription, nil
}

// Update updates an existing subscription by organisation ID
func (r *repository) Update(orgID uuid.UUID, s *Subscription, tx *gorm.DB) error {
	log.Trace().Str("orgID", orgID.String()).Msg("Update")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&Subscription{}).Where("organisation_id = ?", orgID).Updates(s).Error; err != nil {
		log.Error().Err(err).Msg("Error updating subscription")
		return err
	}
	return nil
}

// UpdateWithNil updates specific fields of a subscription by organisation ID, allowing nil values
func (r *repository) UpdateWithNil(orgID uuid.UUID, data map[string]interface{}, tx *gorm.DB) error {
	log.Trace().Str("orgID", orgID.String()).Msg("UpdateWithNil")
	var client *gorm.DB
	if tx != nil {
		client = tx
	} else {
		client = r.db.Client
	}

	if err := client.Model(&Subscription{}).Where("organisation_id = ?", orgID).Updates(data).Error; err != nil {
		log.Error().Err(err).Msg("Error updating subscription")
		return err
	}
	return nil
}

func (r *repository) GetByOrgID(orgID uuid.UUID) (*Subscription, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("GetByOrgID")
	var subscription Subscription

	if err := r.db.Client.Where("organisation_id = ?", orgID).First(&subscription).Error; err != nil {
		log.Error().Err(err).Msg("Error finding subscription")
		return nil, err
	}
	return &subscription, nil
}
