package widget

import (
	"github.com/devbydaniel/announcable/internal/handler/shared"
)

// Handlers provides widget API handlers
type Handlers struct {
	*shared.Dependencies
}

// New creates a new widget API handlers instance
func New(deps *shared.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}
