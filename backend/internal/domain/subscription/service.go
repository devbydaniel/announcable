package subscription

import (
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/google/uuid"
)

var log = logger.Get()

type service struct {
	repo repository
}

func NewService(r repository) *service {
	log.Trace().Msg("NewService")
	return &service{repo: r}
}

// Get retrieves a subscription by organisation ID
func (s *service) Get(orgID uuid.UUID) (*Subscription, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("Get")
	return s.repo.Get(orgID)
}

// GetByStripeSubscriptionID retrieves a subscription by Stripe subscription ID
func (s *service) GetByStripeSubscriptionID(stripeSubID string) (*Subscription, error) {
	log.Trace().Str("stripeSubID", stripeSubID).Msg("GetByStripeSubscriptionID")
	return s.repo.GetByStripeSubscriptionID(stripeSubID)
}

// Create creates a new subscription
func (s *service) Create(sub *Subscription) error {
	log.Trace().Msg("Create")
	return s.repo.Create(sub, nil)
}

// Update updates an existing subscription
func (s *service) Update(orgID uuid.UUID, sub *Subscription) error {
	log.Trace().Str("orgID", orgID.String()).Msg("Update")
	return s.repo.Update(orgID, sub, nil)
}

// UpdateFields updates specific fields of a subscription
func (s *service) UpdateFields(orgID uuid.UUID, fields map[string]interface{}) error {
	log.Trace().Str("orgID", orgID.String()).Msg("UpdateFields")
	return s.repo.UpdateWithNil(orgID, fields, nil)
}

func (s *service) GetByOrgID(orgID uuid.UUID) (*Subscription, error) {
	return s.repo.Get(orgID)
}

// IsFreeSubscription checks if the subscription is a free subscription
func (s *service) IsFreeSubscription(orgID uuid.UUID) (bool, error) {
	log.Trace().Str("orgID", orgID.String()).Msg("IsFreeSubscription")
	sub, err := s.repo.Get(orgID)
	if err != nil {
		return false, err
	}
	return sub.IsFree, nil
}
