package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/devbydaniel/release-notes-go/templates"
)

var subscriptionCancelTmpl = templates.Construct(
	"subscription-cancel",
	"layouts/root.html",
	"layouts/fullscreenmessage.html",
	"pages/subscription-confirm.html",
)

func (h *Handler) HandleSubscriptionCancel(w http.ResponseWriter, r *http.Request) {
	h.log.Trace().Str("path", r.URL.Path).Str("method", r.Method).Msg("HandleSubscriptionCancel")
	
	data := map[string]interface{}{
		"Title": "Subscription Cancelled",
		"Message": "Your subscription was cancelled. You can still use the app until your current period ends.",
	}
	
	if err := subscriptionCancelTmpl.ExecuteTemplate(w, "root", data); err != nil {
		fmt.Fprintf(os.Stderr, "HandleSubscriptionCancel: Error executing template: %v\n", err)
		h.log.Error().Err(err).Msg("Error rendering page")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
} 