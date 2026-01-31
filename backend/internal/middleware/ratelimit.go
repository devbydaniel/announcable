package mw

import (
	"fmt"
	"net/http"

	"github.com/devbydaniel/announcable/config"
	"github.com/devbydaniel/announcable/internal/ratelimit"
	"github.com/devbydaniel/announcable/internal/util"
	"github.com/go-chi/chi/v5"
)

const (
	refillIntervallSeconds = 20
	maxValue               = 10
	costPerRequest         = 2
)

var (
	tbr                   = ratelimit.New(refillIntervallSeconds, maxValue)
	publicIPRateLimiter   ratelimit.RaterLimiter
	publicOrgRateLimiter  ratelimit.RaterLimiter
	rateLimitConfig       *config.Config
)

func init() {
	cfg := config.New()
	rateLimitConfig = cfg
	publicIPRateLimiter = ratelimit.New(
		int64(cfg.RateLimit.PublicRefillIntervalSeconds),
		float64(cfg.RateLimit.PublicMaxTokens),
	)
	publicOrgRateLimiter = ratelimit.New(
		int64(cfg.RateLimit.PublicRefillIntervalSeconds),
		float64(cfg.RateLimit.PublicMaxTokens),
	)
}

func (h *Handler) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw RateLimit")
		userId := r.Context().Value(UserIDKey).(string)
		remaining, err := tbr.Deduct(userId, costPerRequest)
		if err != nil {
			h.log.Warn().Str("userId", userId).Err(err).Msg("")
			http.Error(w, "Rate limit reached", http.StatusTooManyRequests)
			return
		}
		h.log.Trace().Str("userId", userId).Float64("remaining", remaining).Msg("Rate limit check passed")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RateLimitPublicByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw RateLimitPublicByIP")
		clientIP := util.GetClientIP(r)

		limit := rateLimitConfig.RateLimit.PublicMaxTokens
		cost := float64(rateLimitConfig.RateLimit.PublicMaxTokens) / float64(rateLimitConfig.RateLimit.PublicRequestsPerInterval)

		remaining, err := publicIPRateLimiter.Deduct(clientIP, cost)

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", int(remaining)))

		if err != nil {
			h.log.Warn().Str("clientIP", clientIP).Err(err).Msg("Rate limit exceeded")
			w.Header().Set("X-RateLimit-Retry-After", fmt.Sprintf("%d", rateLimitConfig.RateLimit.PublicRefillIntervalSeconds))
			http.Error(w, "Rate limit reached", http.StatusTooManyRequests)
			return
		}

		h.log.Trace().Str("clientIP", clientIP).Float64("remaining", remaining).Msg("Rate limit check passed")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RateLimitPublicByOrg(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Trace().Msg("mw RateLimitPublicByOrg")
		orgId := chi.URLParam(r, "orgId")
		if orgId == "" {
			h.log.Warn().Msg("orgId not found in URL parameters")
			next.ServeHTTP(w, r)
			return
		}

		limit := rateLimitConfig.RateLimit.PublicMaxTokens
		cost := float64(rateLimitConfig.RateLimit.PublicMaxTokens) / float64(rateLimitConfig.RateLimit.PublicRequestsPerInterval)

		remaining, err := publicOrgRateLimiter.Deduct(orgId, cost)

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", int(remaining)))

		if err != nil {
			h.log.Warn().Str("orgId", orgId).Err(err).Msg("Rate limit exceeded")
			w.Header().Set("X-RateLimit-Retry-After", fmt.Sprintf("%d", rateLimitConfig.RateLimit.PublicRefillIntervalSeconds))
			http.Error(w, "Rate limit reached", http.StatusTooManyRequests)
			return
		}

		h.log.Trace().Str("orgId", orgId).Float64("remaining", remaining).Msg("Rate limit check passed")
		next.ServeHTTP(w, r)
	})
}
