package stripeUtil

import (
	"encoding/json"
	"fmt"

	"github.com/devbydaniel/release-notes-go/config"
	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	portalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/webhook"
)

var log = logger.Get()

func Setup() {
	stripe.Key = config.New().Payment.StripeKey
}

func CreateStripeCheckoutSession(lookupKey string, organizationID uuid.UUID) (string, error) {
	baseUrl := config.New().BaseURL
	params := &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{
			lookupKey,
		}),
	}
	i := price.List(params)

	var price *stripe.Price
	for i.Next() {
		p := i.Price()
		price = p
	}

	// controls what the customer sees on the checkout page
	checkoutParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(price.ID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://" + baseUrl + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://" + baseUrl + "/cancel.html"),
		// Add metadata to link this session to your organization
		Metadata: map[string]string{
			"organization_id": organizationID.String(),
		},
		// Also add metadata to the subscription
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"organization_id": fmt.Sprintf("%d", organizationID),
			},
		},
	}

	s, err := session.New(checkoutParams)
	if err != nil {
		return "", fmt.Errorf("session.New: %w", err)
	}
	log.Debug().Str("session_id", s.ID).Str("orgId", organizationID.String()).Msg("Checkout session created")

	return s.URL, nil
}

// CreateStripePortalSession creates a Stripe billing portal session for a customer
// based on the provided checkout session ID and return URL.
func CreateStripePortalSession(checkoutSessionID, returnURL string) (string, error) {
	// Get the checkout session to retrieve the customer ID
	s, err := session.Get(checkoutSessionID, nil)
	if err != nil {
		return "", fmt.Errorf("session.Get: %w", err)
	}

	// Create the billing portal session
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(s.Customer.ID),
		ReturnURL: stripe.String(returnURL),
	}

	ps, err := portalsession.New(params)
	if err != nil {
		return "", fmt.Errorf("portalsession.New: %w", err)
	}

	return ps.URL, nil
}

// VerifyStripeWebhook verifies the Stripe webhook signature
func VerifyStripeWebhook(payload []byte, signatureHeader, endpointSecret string) (stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		return stripe.Event{}, fmt.Errorf("webhook signature verification failed: %w", err)
	}
	return event, nil
}

func ParseSubscription(payload []byte) (stripe.Subscription, error) {
	var subscription stripe.Subscription
	if err := json.Unmarshal(payload, &subscription); err != nil {
		return stripe.Subscription{}, err
	}
	return subscription, nil
}
