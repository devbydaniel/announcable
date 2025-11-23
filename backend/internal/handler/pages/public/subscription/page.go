package subscription

import (
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
)

// Handlers provides public subscription handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new public subscription handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}
