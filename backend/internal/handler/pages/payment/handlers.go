package payment

import (
	"github.com/devbydaniel/release-notes-go/internal/handler/shared"
)

// Handlers provides payment handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new payment handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}
