package anim

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlossEnhancedAnim(t *testing.T) {
	t.Run("NewFlossEnhanced", func(t *testing.T) {
		opts := FlossAnimationSettings{
			Size:  10,
			Label: "Testing",
		}
		anim := NewFlossEnhanced(opts)
		require.NotNil(t, anim)
	})
}