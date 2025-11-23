package payment

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/stripeUtil"
	"github.com/google/uuid"
)

// HandlePortalSession creates a Stripe portal session
func (h *Handlers) HandlePortalSession(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("CreatePortalSession")

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
		h.Log.Error().Err(err).Msg("Error getting subscription")
		http.Error(w, "Error getting subscription details", http.StatusInternalServerError)
		return
	}

	if sub.StripeSubscriptionID == "" {
		h.Log.Error().Msg("No Stripe subscription ID found")
		http.Error(w, "No active subscription found", http.StatusBadRequest)
		return
	}

	// Get the customer ID from Stripe
	customerID, err := stripeUtil.GetCustomerIDFromSubscription(sub.StripeSubscriptionID)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error getting customer ID from subscription")
		http.Error(w, "Error getting customer details", http.StatusInternalServerError)
		return
	}

	returnUrl := "http://" + config.New().BaseURL
	portalURL, err := stripeUtil.CreateStripePortalSession(customerID, returnUrl)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error creating portal session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Log.Debug().Str("portalURL", portalURL).Msg("CreatePortalSession")
	http.Redirect(w, r, portalURL, http.StatusSeeOther)
}
