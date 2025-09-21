package status

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnhancedStatusCmp(t *testing.T) {
	t.Run("NewEnhancedStatusCmp", func(t *testing.T) {
		cmp := NewEnhancedStatusCmp()
		require.NotNil(t, cmp)
	})

	t.Run("EnhancedMode", func(t *testing.T) {
		cmp := NewEnhancedStatusCmp()
		enhancedCmp, ok := cmp.(EnhancedStatusCmp)
		require.True(t, ok)
		
		// Test default mode
		require.False(t, enhancedCmp.IsEnhancedMode())
		
		// Test setting enhanced mode
		enhancedCmp.SetEnhancedMode(true)
		require.True(t, enhancedCmp.IsEnhancedMode())
		
		// Test unsetting enhanced mode
		enhancedCmp.SetEnhancedMode(false)
		require.False(t, enhancedCmp.IsEnhancedMode())
	})
}