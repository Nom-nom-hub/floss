package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

func NewFlossTheme() *Theme {
	t := &Theme{
		Name:   "floss",
		IsDark: true,

		// New color palette for Floss
		Primary:   charmtone.Mauve,     // Purple accent
		Secondary: charmtone.Squid,     // Teal
		Tertiary:  charmtone.Mustard,   // Yellow
		Accent:    charmtone.Cheeky,    // Pink

		// Backgrounds
		BgBase:        charmtone.Pepper,   // Dark background
		BgBaseLighter: charmtone.BBQ,      // Slightly lighter
		BgSubtle:      charmtone.Charcoal, // Subtle background
		BgOverlay:     charmtone.Iron,     // Overlay

		// Foregrounds
		FgBase:      charmtone.Ash,     // Base text
		FgMuted:     charmtone.Squid,   // Muted text
		FgHalfMuted: charmtone.Smoke,   // Half-muted
		FgSubtle:    charmtone.Oyster,  // Subtle text
		FgSelected:  charmtone.Salt,    // Selected text

		// Borders
		Border:      charmtone.Charcoal,
		BorderFocus: charmtone.Mauve,

		// Status
		Success: charmtone.Julep,    // Green
		Error:   charmtone.Sriracha, // Red
		Warning: charmtone.Zest,     // Orange
		Info:    charmtone.Malibu,   // Blue

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

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(charmtone.Salt).Background(charmtone.Mauve)

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(charmtone.Squid).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(charmtone.Mustard)
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(charmtone.Coral)
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(charmtone.Julep)

	t.YoloIconFocused = lipgloss.NewStyle().Foreground(charmtone.Oyster).Background(charmtone.Mustard).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(charmtone.Pepper).Background(charmtone.Squid)
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(charmtone.Zest).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(charmtone.Squid)

	return t
}