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

func (h *Handler) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw RateLimit")
		userId := r.Context().Value(UserIDKey).(string)
		if err := tbr.Deduct(userId, costPerRequest); err != nil {
			h.log.Warn().Str("userId", userId).Err(err).Msg("")
			http.Error(w, "Rate limit reached", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
