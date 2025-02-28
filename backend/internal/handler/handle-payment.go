package handler

import (
	"io"
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
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
	domain := config.New().BaseURL

	// For demonstration purposes, we're using the Checkout session to retrieve the customer ID.
	// Typically this is stored alongside the authenticated user in your database.
	r.ParseForm()
	sessionID := r.PostFormValue("session_id")

	portalURL, err := stripeUtil.CreateStripePortalSession(sessionID, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Printf("CreateStripePortalSession: %v", err)
		return
	}

	h.log.Trace().Str("portalURL", portalURL).Msg("CreatePortalSession")
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

	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	endpointSecret := config.New().Payment.WebhookSecret
	signatureHeader := req.Header.Get("Stripe-Signature")

	event, err := stripeUtil.VerifyStripeWebhook(payload, signatureHeader, endpointSecret)
	if err != nil {
		h.log.Error().Err(err).Msg("Webhook verification failed")
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	subscription, err := stripeUtil.ParseSubscription(event.Data.Raw)
	if err != nil {
		h.log.Error().Err(err).Msg("Error parsing subscription")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		h.log.Info().Str("subscription.ID", subscription.ID).Msg("Subscription deleted.")
		orgId := subscription.Metadata["organization_id"]
		if orgId == "" {
			h.log.Error().Msg("Organization ID not found in subscription metadata")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.log.Debug().Str("orgId", orgId).Msg("Organization ID")

		// Then define and call a func to handle the deleted subscription.
		// handleSubscriptionCanceled(subscription)
	case "customer.subscription.updated":
		h.log.Info().Str("subscription.ID", subscription.ID).Msg("Subscription updated.")
		orgId := subscription.Metadata["organization_id"]
		if orgId == "" {
			h.log.Error().Msg("Organization ID not found in subscription metadata")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.log.Debug().Str("orgId", orgId).Msg("Organization ID")

		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handleSubscriptionUpdated(subscription)
	case "customer.subscription.created":
		h.log.Info().Str("subscription.ID", subscription.ID).Msg("Subscription created.")
		orgId := subscription.Metadata["organization_id"]
		if orgId == "" {
			h.log.Error().Msg("Organization ID not found in subscription metadata")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.log.Debug().Str("orgId", orgId).Msg("Organization ID")

		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handleSubscriptionCreated(subscription)
	default:
		h.log.Debug().Str("eventType", string(event.Type)).Msg("Unhandled event type")
	}

	w.WriteHeader(http.StatusOK)
}
