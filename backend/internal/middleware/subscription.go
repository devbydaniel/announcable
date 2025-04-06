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
		var hasActiveSubscription bool
		subscriptionService := subscription.NewService(*subscription.NewRepository(h.DB))
		sub, err := subscriptionService.Get(uuid.MustParse(orgId))
		if err != nil {
			hasActiveSubscription = false
		} else {
			hasActiveSubscription = sub.IsActive || sub.IsFree
		}

		// Add subscription status to context
		ctx = context.WithValue(ctx, HasActiveSubscription, hasActiveSubscription)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
