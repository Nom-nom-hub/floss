package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

func NewFlossTheme() *Theme {
	t := &Theme{
		Name:   "floss",
		IsDark: true,

		// New color palette for Floss - updated to differentiate from Crush
		Primary:   charmtone.Citron,  // Pink-Red accent (was Mauve/Purple)
		Secondary: charmtone.Guac,    // Deep teal-blue (was Squid/Teal)
		Tertiary:  charmtone.Mustard, // Yellow (unchanged)
		Accent:    charmtone.Cheeky,  // Pink (reserved for warnings/errors)

		// Backgrounds - adjusted for better Citron integration
		BgBase:        charmtone.Pepper,   // Dark background (unchanged)
		BgBaseLighter: charmtone.BBQ,      // Slightly lighter (increased contrast)
		BgSubtle:      ParseHex("#221f30"), // Warmed Charcoal for better Citron integration
		BgOverlay:     charmtone.Iron,     // Overlay (unchanged)

		// Foregrounds - maintained core text hierarchy
		FgBase:      charmtone.Ash,    // Base text (unchanged)
		FgMuted:     charmtone.Squid,  // Muted text (unchanged)
		FgHalfMuted: charmtone.Smoke,  // Half-muted (unchanged)
		FgSubtle:    charmtone.Oyster, // Subtle text (unchanged)
		FgSelected:  charmtone.Salt,   // Selected text (unchanged)

		// Borders - adjusted for FLOSS styling
		Border:      charmtone.Charcoal,
		BorderFocus: charmtone.Citron, // Citron for focus (was Mauve)

		// Status - maintained but with adjusted applications
		Success: charmtone.Julep,    // Green (unchanged)
		Error:   charmtone.Sriracha, // Red (unchanged)
		Warning: charmtone.Zest,     // Orange (unchanged)
		Info:    charmtone.Guac,     // Blue-Green (was Malibu, now matches Secondary)

		// Colors
		White: charmtone.Butter,

		BlueLight: charmtone.Sardine,
		Blue:      charmtone.Malibu,

		Yellow: charmtone.Mustard,
		Citron: charmtone.Citron,

		Green:      charmtone.Julep,
		GreenDark:  charmtone.Guac,
		GreenLight: charmtone.Bok,

		Red:      charmtone.Coral,
		RedDark:  charmtone.Sriracha,
		RedLight: charmtone.Salmon,
		Cherry:   charmtone.Cherry,
	}

	// Text selection - updated to use Citron instead of Mauve
	t.TextSelection = lipgloss.NewStyle().Foreground(charmtone.Salt).Background(charmtone.Citron)

	// LSP and MCP status - updated to use FLOSS color scheme
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(charmtone.Squid).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(charmtone.Mustard)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(charmtone.Coral)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(charmtone.Julep)

	// Yolo mode styling - updated to use FLOSS colors
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(charmtone.Oyster).Background(charmtone.Citron).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(charmtone.Pepper).Background(charmtone.Guac)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(charmtone.Zest).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(charmtone.Guac)

	return t
}
