package shared

import (
	sharedDeps "github.com/devbydaniel/release-notes-go/internal/handler/shared"
)

// Handlers provides shared API handlers
type Handlers struct {
	*sharedDeps.Dependencies
}

// New creates a new shared API handlers instance
func New(deps *sharedDeps.Dependencies) *Handlers {
	return &Handlers{Dependencies: deps}
}
