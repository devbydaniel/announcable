package organisation

import (
	"net/http"
	"net/url"

	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type subscriptionCreateForm struct {
	StripeSubscriptionID string `schema:"stripe_subscription_id"`
	IsActive             bool   `schema:"is_active"`
	IsFree               bool   `schema:"is_free"`
}

// HandleSubscriptionCreate creates a subscription for an organisation
func (h *Handlers) HandleSubscriptionCreate(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleSubscriptionCreate")

	// Get organisation ID from URL params
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		h.Log.Error().Str("orgId", orgIDStr).Err(err).Msg("Error parsing organisation ID")
		http.Error(w, "Invalid organisation ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.Log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Decode form
	var createDTO subscriptionCreateForm
	if err := h.Decoder.Decode(&createDTO, r.PostForm); err != nil {
		h.Log.Error().Err(err).Msg("Error decoding form")
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	// Validate form
	validate := validator.New()
	if err := validate.Struct(createDTO); err != nil {
		h.Log.Error().Err(err).Msg("Validation error")
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	// Create subscription service
	subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))

	// Create subscription
	sub := &subscription.Subscription{
		OrganisationID: orgID,
		IsActive:       createDTO.IsActive,
		IsFree:         createDTO.IsFree,
	}
	if createDTO.StripeSubscriptionID != "" {
		sub.StripeSubscriptionID = createDTO.StripeSubscriptionID
	}

	err = subscriptionService.Create(sub)
	if err != nil {
		h.Log.Error().Err(err).Msg("Error creating subscription")
		http.Error(w, "Error creating subscription", http.StatusInternalServerError)
		return
	}

	// Redirect back to the organisation details page
	successMsg := "subscription created"
	escapedMsg := url.QueryEscape(successMsg)
	redirectURL := "/admin/organisations/" + orgID.String() + "?success=" + escapedMsg
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusCreated)
}
