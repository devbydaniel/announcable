package handler

import (
	"io"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/stripeUtil"
	"github.com/google/uuid"
)

func (h *Handler) HandleCheckoutSession(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("CreateCheckoutSession")
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		http.Error(w, "Error getting organization ID", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	lookupKey := r.PostFormValue("lookup_key")
	sessionURL, err := stripeUtil.CreateStripeCheckoutSession(lookupKey, uuid.MustParse(orgId))
	if err != nil {
		h.log.Printf("CreateStripeCheckoutSession: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Debug().Str("sessionURL", sessionURL).Msg("CreateCheckoutSession")
	http.Redirect(w, r, sessionURL, http.StatusSeeOther)
}

func (h *Handler) HandlePortalSession(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("CreatePortalSession")

	// Get organization ID from context
	orgID, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		http.Error(w, "Error getting organization ID", http.StatusInternalServerError)
		return
	}

	// Get the subscription details from database
	subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))
	sub, err := subscriptionService.GetByOrgID(uuid.MustParse(orgID))
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting subscription")
		http.Error(w, "Error getting subscription details", http.StatusInternalServerError)
		return
	}

	if sub.StripeSubscriptionID == "" {
		h.log.Error().Msg("No Stripe subscription ID found")
		http.Error(w, "No active subscription found", http.StatusBadRequest)
		return
	}

	// Get the customer ID from Stripe
	customerID, err := stripeUtil.GetCustomerIDFromSubscription(sub.StripeSubscriptionID)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting customer ID from subscription")
		http.Error(w, "Error getting customer details", http.StatusInternalServerError)
		return
	}

	returnUrl := "http://" + config.New().BaseURL
	portalURL, err := stripeUtil.CreateStripePortalSession(customerID, returnUrl)
	if err != nil {
		h.log.Error().Err(err).Msg("Error creating portal session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Debug().Str("portalURL", portalURL).Msg("CreatePortalSession")
	http.Redirect(w, r, portalURL, http.StatusSeeOther)
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	h.log.Trace().Msg("HandleWebhook")

	const MaxBodyBytes = int64(65536)
	bodyReader := http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(bodyReader)
	if err != nil {
		h.log.Error().Err(err).Msg("Error reading request body")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := config.New().Payment.WebhookSecret
	signatureHeader := req.Header.Get("Stripe-Signature")

	event, err := stripeUtil.VerifyStripeWebhook(payload, signatureHeader, endpointSecret)
	if err != nil {
		h.log.Error().Err(err).Msg("Webhook verification failed")
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
		h.log.Debug().Str("eventType", string(event.Type)).Msg("Unhandled event type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only parse subscription and initialize service if we have a relevant event
	stripeSubscription, err := stripeUtil.ParseSubscription(event.Data.Raw)
	if err != nil {
		h.log.Error().Err(err).Msg("Error parsing subscription")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))

	orgId := stripeSubscription.Metadata["organization_id"]
	if orgId == "" {
		h.log.Error().Msg("Organization ID not found in subscription metadata")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		h.log.Error().Err(err).Msg("Invalid organization ID format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		h.log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription deleted.")
		if err := subscriptionService.UpdateFields(orgUUID, map[string]interface{}{
			"is_active": false,
		}); err != nil {
			h.log.Error().Err(err).Msg("Error updating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.updated":
		h.log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription updated.")
		isActive := stripeSubscription.Status == "active"
		if err := subscriptionService.UpdateFields(orgUUID, map[string]interface{}{
			"is_active": isActive,
		}); err != nil {
			h.log.Error().Err(err).Msg("Error updating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "customer.subscription.created":
		h.log.Info().Str("subscription.ID", stripeSubscription.ID).Msg("Subscription created.")
		sub := &subscription.Subscription{
			OrganisationID:       orgUUID,
			StripeSubscriptionID: stripeSubscription.ID,
			IsActive:             stripeSubscription.Status == "active",
			IsFree:               false,
		}

		if err := subscriptionService.Create(sub); err != nil {
			h.log.Error().Err(err).Msg("Error creating subscription")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
