package handler

import (
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/templates"
)

type subscriptionSuccessPageData struct {
	BaseTemplateData
}

var subscriptionSuccessPageTmpl = templates.Construct(
	"subscription-success",
	"layouts/root.html",
	"layouts/appframe.html",
	"pages/subscription-success.html",
)

func (h *Handler) HandleSubscriptionSuccess(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Msg("HandleSubscriptionSuccess")

	// Get subscription status from context
	hasActiveSubscription, ok := r.Context().Value(mw.HasActiveSubscription).(bool)
	if !ok {
		h.log.Error().Msg("Subscription status not found in context")
		http.Error(w, "Error checking subscription status", http.StatusInternalServerError)
		return
	}

	// Create empty subscription struct to check status below
	sub := &subscription.Subscription{
		IsActive: hasActiveSubscription,
		IsFree:   hasActiveSubscription, // If active from context, could be either paid or free
	}
	if !sub.IsActive && !sub.IsFree {
		h.log.Warn().Msg("Subscription not active")
		http.Error(w, "Subscription not active", http.StatusBadRequest)
		return
	}

	// Render the success page
	pageData := subscriptionSuccessPageData{
		BaseTemplateData: BaseTemplateData{
			Title:                 "Subscription Activated",
			HasActiveSubscription: hasActiveSubscription,
		},
	}

	if err := subscriptionSuccessPageTmpl.ExecuteTemplate(w, "root", pageData); err != nil {
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}
