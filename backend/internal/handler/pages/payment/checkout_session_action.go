package payment

import (
	"net/http"

	mw "github.com/devbydaniel/release-notes-go/internal/middleware"
	"github.com/devbydaniel/release-notes-go/internal/stripeUtil"
	"github.com/google/uuid"
)

// HandleCheckoutSession creates a Stripe checkout session
func (h *Handlers) HandleCheckoutSession(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("CreateCheckoutSession")
	orgId, ok := r.Context().Value(mw.OrgIDKey).(string)
	if !ok {
		http.Error(w, "Error getting organization ID", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	lookupKey := r.PostFormValue("lookup_key")
	sessionURL, err := stripeUtil.CreateStripeCheckoutSession(lookupKey, uuid.MustParse(orgId))
	if err != nil {
		h.Log.Printf("CreateStripeCheckoutSession: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Log.Debug().Str("sessionURL", sessionURL).Msg("CreateCheckoutSession")
	http.Redirect(w, r, sessionURL, http.StatusSeeOther)
}
