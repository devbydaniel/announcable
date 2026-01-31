package mw

import (
	"net/http"

	"github.com/devbydaniel/announcable/internal/ratelimit"
)

const (
	refillIntervallSeconds = 20
	maxValue               = 10
	costPerRequest         = 2
)

var tbr = ratelimit.New(refillIntervallSeconds, maxValue)

// RateLimit is middleware that enforces per-user request rate limiting.
func (h *Handler) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw RateLimit")
		userID := r.Context().Value(UserIDKey).(string)
		if err := tbr.Deduct(userID, costPerRequest); err != nil {
			h.log.Warn().Str("userID", userID).Err(err).Msg("")
			http.Error(w, "Rate limit reached", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
