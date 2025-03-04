package mw

import (
	"context"
	"net/http"

	"github.com/devbydaniel/release-notes-go/internal/domain/subscription"
	"github.com/google/uuid"
)

func (h *Handler) WithSubscriptionStatus(next http.Handler) http.Handler {
	h.log.Trace().Msg("WithSubscriptionStatus")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get organization ID from context
		orgId, ok := ctx.Value(OrgIDKey).(string)
		if !ok {
			h.log.Error().Msg("Organisation ID not found in context")
			http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
			return
		}

		// Get subscription status
		subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))
		sub, err := subscriptionService.Get(uuid.MustParse(orgId))
		h.log.Debug().Interface("sub", sub).Msg("Subscription status")
		hasActiveSubscription := false
		if err == nil {
			hasActiveSubscription = sub.IsActive || sub.IsFree
		} else {
			h.log.Error().Err(err).Msg("Error getting subscription status")
		}

		// Add subscription status to context
		ctx = context.WithValue(ctx, HasActiveSubscription, hasActiveSubscription)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
