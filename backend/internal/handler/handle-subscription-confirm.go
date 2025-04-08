package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/devbydaniel/release-notes-go/templates"
)

var subscriptionConfirmTmpl = templates.Construct(
	"subscription-confirm",
	"layouts/root.html",
	"layouts/fullscreenmessage.html",
	"pages/subscription-confirm.html",
)

func (h *Handler) HandleSubscriptionConfirm(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Str("path", r.URL.Path).Str("method", r.Method).Msg("HandleSubscriptionConfirm")
	
	data := map[string]interface{}{
		"Title": "You're all set! 🚀",
		"Message": "You've successfully subscribed to our service. Welcome aboard!",
	}
	
	if err := subscriptionConfirmTmpl.ExecuteTemplate(w, "root", data); err != nil {
		fmt.Fprintf(os.Stderr, "HandleSubscriptionConfirm: Error executing template: %v\n", err)
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
} 