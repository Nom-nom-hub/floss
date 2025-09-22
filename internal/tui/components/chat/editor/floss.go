package editor

import (
	"github.com/nom-nom-hub/floss/internal/app"
)

// NewFloss creates a new FLOSS-styled editor
func NewFloss(app *app.App) Editor {
	// For now, return the standard editor with FLOSS styling
	// In the future, this could return a fully customized FLOSS editor
	return New(app)
}