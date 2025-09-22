package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
)

// EnhancedStyle provides consistent styling for UI components based on the FLOSS design system
type EnhancedStyle struct {
	// Base styles
	Base        lipgloss.Style
	Primary     lipgloss.Style
	Secondary   lipgloss.Style
	Background  lipgloss.Style
	Border      lipgloss.Style
	
	// Text styles
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Body        lipgloss.Style
	Caption     lipgloss.Style
	Muted       lipgloss.Style
	Subtle      lipgloss.Style
	
	// Status styles
	Success     lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Info        lipgloss.Style
	
	// Interactive styles
	Focused     lipgloss.Style
	Selected    lipgloss.Style
	Disabled    lipgloss.Style
	
	// Component specific styles
	Button      lipgloss.Style
	Input       lipgloss.Style
	Card        lipgloss.Style
	Panel       lipgloss.Style
}

// NewEnhancedStyle creates a new enhanced style configuration based on the current theme
func NewEnhancedStyle() *EnhancedStyle {
	t := CurrentTheme()
	
	// Create base styles with proper spacing
	base := lipgloss.NewStyle().
		Padding(S).
		Foreground(t.FgBase).
		Background(t.BgBase)
		
	primary := base.
		Foreground(t.Primary)
		
	secondary := base.
		Foreground(t.Secondary)
		
	background := base.
		Background(t.BgSubtle)
		
	border := base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
		
	// Text hierarchy styles
	title := base.
		Foreground(t.Accent).
		Bold(true).
		PaddingBottom(S)
		
	subtitle := base.
		Foreground(t.Secondary).
		Bold(true)
		
	body := base.
		Foreground(t.FgBase)
		
	caption := base.
		Foreground(t.FgSubtle).
		Italic(true)
		
	muted := base.
		Foreground(t.FgMuted)
		
	subtle := base.
		Foreground(t.FgSubtle)
		
	// Status styles
	success := base.
		Foreground(t.Success)
		
	errorStyle := base.
		Foreground(t.Error)
		
	warning := base.
		Foreground(t.Warning)
		
	info := base.
		Foreground(t.Info)
		
	// Interactive styles
	focused := base.
		BorderForeground(t.BorderFocus).
		BorderStyle(lipgloss.ThickBorder())
		
	selected := base.
		Background(t.Primary).
		Foreground(t.FgSelected)
		
	disabled := base.
		Foreground(t.FgMuted).
		Faint(true)
		
	// Component specific styles
	button := base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(XS, S).
		Background(t.BgOverlay)
		
	input := base.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(t.Border).
		Padding(XS, S)
		
	card := base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(M).
		Background(t.BgBaseLighter)
		
	panel := base.
		BorderStyle(lipgloss.HiddenBorder()).
		Padding(M).
		Background(t.BgBase)

	return &EnhancedStyle{
		Base:        base,
		Primary:     primary,
		Secondary:   secondary,
		Background:  background,
		Border:      border,
		Title:       title,
		Subtitle:    subtitle,
		Body:        body,
		Caption:     caption,
		Muted:       muted,
		Subtle:      subtle,
		Success:     success,
		Error:       errorStyle,
		Warning:     warning,
		Info:        info,
		Focused:     focused,
		Selected:    selected,
		Disabled:    disabled,
		Button:      button,
		Input:       input,
		Card:        card,
		Panel:       panel,
	}
}

// ApplyFocus applies focus styling to a component
func (e *EnhancedStyle) ApplyFocus(style lipgloss.Style) lipgloss.Style {
	return style.
		BorderForeground(CurrentTheme().BorderFocus).
		BorderStyle(lipgloss.DoubleBorder())
}

// ApplySelection applies selection styling to a component
func (e *EnhancedStyle) ApplySelection(style lipgloss.Style) lipgloss.Style {
	return style.
		Background(CurrentTheme().Primary).
		Foreground(CurrentTheme().FgSelected)
}

// ApplyDisabled applies disabled styling to a component
func (e *EnhancedStyle) ApplyDisabled(style lipgloss.Style) lipgloss.Style {
	return style.
		Foreground(CurrentTheme().FgMuted).
		Faint(true)
}

// ApplySuccess applies success styling to a component
func (e *EnhancedStyle) ApplySuccess(style lipgloss.Style) lipgloss.Style {
	return style.
		Foreground(CurrentTheme().Success)
}

// ApplyError applies error styling to a component
func (e *EnhancedStyle) ApplyError(style lipgloss.Style) lipgloss.Style {
	return style.
		Foreground(CurrentTheme().Error)
}

// ApplyWarning applies warning styling to a component
func (e *EnhancedStyle) ApplyWarning(style lipgloss.Style) lipgloss.Style {
	return style.
		Foreground(CurrentTheme().Warning)
}