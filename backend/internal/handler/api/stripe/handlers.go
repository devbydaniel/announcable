package stripe

import (
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
)

// Handlers provides Stripe API handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new Stripe API handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}
