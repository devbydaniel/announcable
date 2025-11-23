package stripe

import (
	"io"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	"github.com/devbydaniel/release-notes-go/internal/stripeUtil"
	"github.com/google/uuid"
)

// HandleWebhook handles Stripe webhook events
func (h *Handlers) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	h.Log.Trace().Msg("HandleWebhook")

	const MaxBodyBytes = int64(65536)
	bodyReader := http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(bodyReader)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error reading request body")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := config.New().Payment.WebhookSecret
	signatureHeader := req.Header.Get("Stripe-Signature")

	event, err := stripeUtil.VerifyStripeWebhook(payload, signatureHeader, endpointSecret)
	if err != nil {
		h.Log.Error().Err(err).Msg("Webhook verification failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if this is a subscription-related event we care about
	relevantEvents := map[string]bool{
		"customer.subscription.deleted": true,
		"customer.subscription.updated": true,
		"customer.subscription.created": true,
	}

	if !relevantEvents[string(event.Type)] {
		h.Log.Debug().Str("eventType", string(event.Type)).Msg("Unhandled event type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only parse subscription and initialize service if we have a relevant event
	stripeSubscription, err := stripeUtil.ParseSubscription(event.Data.Raw)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error parsing subscription")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))

	orgId := stripeSubscription.Metadata["organization_id"]
	if orgId == "" {
		h.Log.Error().Msg("Organization ID not found in subscription metadata")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		h.Log.Error().Err(err).Msg("Invalid organization ID format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		h.Log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription deleted.")
		if err := subscriptionService.UpdateFields(orgUUID, map[string]interface{}{
			"is_active": false,
		}); err != nil {
			h.Log.Error().Err(err).Msg("Error updating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.updated":
		h.Log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription updated.")
		isActive := stripeSubscription.Status == "active"
		if err := subscriptionService.UpdateFields(orgUUID, map[string]interface{}{
			"is_active": isActive,
		}); err != nil {
			h.Log.Error().Err(err).Msg("Error updating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.created":
		h.Log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription created.")
		sub := &subscription.Subscription{
			OrganisationID:       orgUUID,
			StripeSubscriptionID: stripeSubscription.ID,
			IsActive:             stripeSubscription.Status == "active",
			IsFree:               false,
		}

		if err := subscriptionService.Create(sub); err != nil {
			h.Log.Error().Err(err).Msg("Error creating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
