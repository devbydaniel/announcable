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
	"github.com/stripe/stripe-go/v81/subscription"
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
			lookupKey, // Currently 'pro_monthly', change if additional products are added
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
		SuccessURL: stripe.String("http://" + baseUrl + "/payment/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://" + baseUrl + "/payment/cancel"),
		// Add metadata to link this session to your organization
		Metadata: map[string]string{
			"organization_id": organizationID.String(),
		},
		// Also add metadata to the subscription
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"organization_id": organizationID.String(),
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

// GetCustomerIDFromSubscription retrieves the customer ID associated with a Stripe subscription
func GetCustomerIDFromSubscription(subscriptionID string) (string, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return "", fmt.Errorf("error getting subscription: %w", err)
	}
	return sub.Customer.ID, nil
}

// CreateStripePortalSession creates a Stripe billing portal session for a customer
// based on the provided customer ID and return URL.
func CreateStripePortalSession(customerID, returnURL string) (string, error) {
	// Create the billing portal session
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
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
