package anim

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnim(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		opts := Settings{
			Size:  10,
			Label: "Testing",
		}
		anim := New(opts)
		require.NotNil(t, anim)
	})
}