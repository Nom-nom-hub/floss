package floss

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
)

// ButtonStyles defines the FLOSS-specific button styling
type ButtonStyles struct {
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Danger    lipgloss.Style
}

// GetButtonStyles returns the FLOSS-specific button styles
func GetButtonStyles() ButtonStyles {
	theme := styles.CurrentTheme()
	
	return ButtonStyles{
		Primary: lipgloss.NewStyle().
			Background(theme.Primary).          // Citron background
			Foreground(theme.BgBase).           // Pepper text
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).    // Citron border
			Padding(0, 2).                      // 8px 16px equivalent
			Bold(true).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),
		Secondary: lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).     // Transparent background
			Foreground(theme.Secondary).        // Guac text
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Secondary).  // Guac border
			Padding(0, 2).                      // 8px 16px equivalent
			Bold(true).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),
		Danger: lipgloss.NewStyle().
			Background(theme.Accent).           // Cheeky background
			Foreground(theme.FgSelected).       // Salt text
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Accent).     // Cheeky border
			Padding(0, 2).                      // 8px 16px equivalent
			Bold(true).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),
	}
}

// InputStyles defines the FLOSS-specific input field styling
type InputStyles struct {
	TextInput lipgloss.Style
	TextArea  lipgloss.Style
}

// GetInputStyles returns the FLOSS-specific input styles
func GetInputStyles() InputStyles {
	theme := styles.CurrentTheme()
	
	return InputStyles{
		TextInput: lipgloss.NewStyle().
			Background(theme.BgSubtle).         // Warmed Charcoal background
			Foreground(theme.FgBase).           // Ash text
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Border).     // Charcoal border
			Padding(0, 1).                      // 10px equivalent
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),
		TextArea: lipgloss.NewStyle().
			Background(theme.BgSubtle).         // Warmed Charcoal background
			Foreground(theme.FgBase).           // Ash text
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Border).     // Charcoal border
			Padding(0, 1).                      // 10px equivalent
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true),
	}
}

// CardStyles defines the FLOSS-specific card styling
type CardStyles struct {
	BaseCard      lipgloss.Style
	UserMessage   lipgloss.Style
	AssistantMessage lipgloss.Style
}

// GetCardStyles returns the FLOSS-specific card styles
func GetCardStyles() CardStyles {
	theme := styles.CurrentTheme()
	
	return CardStyles{
		BaseCard: lipgloss.NewStyle().
			Background(theme.BgSubtle).         // Charcoal background
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Secondary).  // Guac border
			Padding(1, 2),                      // 20px equivalent
		UserMessage: lipgloss.NewStyle().
			Background(theme.BgBase).           // Pepper background
			Border(lipgloss.ThickBorder()).
			BorderLeft(true).
			BorderForeground(theme.Primary).    // Citron border
			Padding(0, 0, 0, 1).                // 12px left padding
			MarginBottom(1),                    // 12px bottom margin
		AssistantMessage: lipgloss.NewStyle().
			Background(theme.BgBase).           // Pepper background
			Border(lipgloss.ThickBorder()).
			BorderLeft(true).
			BorderForeground(theme.Secondary).  // Guac border
			Padding(0, 0, 0, 2).                // 16px left padding
			MarginBottom(1),                    // 12px bottom margin
	}
}

// DialogStyles defines the FLOSS-specific dialog styling
type DialogStyles struct {
	ModalDialog lipgloss.Style
	Header      lipgloss.Style
	Content     lipgloss.Style
	Actions     lipgloss.Style
}

// GetDialogStyles returns the FLOSS-specific dialog styles
func GetDialogStyles() DialogStyles {
	theme := styles.CurrentTheme()
	
	return DialogStyles{
		ModalDialog: lipgloss.NewStyle().
			Background(theme.BgOverlay).        // Iron background
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).    // Citron border
			Padding(1, 2),                      // 20px equivalent
		Header: lipgloss.NewStyle().
			Foreground(theme.Primary).          // Citron text
			Bold(true).
			MarginBottom(1).
			PaddingBottom(1).
			Border(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(theme.Secondary),  // Guac border
		Content: lipgloss.NewStyle().
			Foreground(theme.FgBase).           // Ash text
			MarginBottom(1),
		Actions: lipgloss.NewStyle().
			Align(lipgloss.Right).
			MarginTop(1).
			PaddingTop(1).
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(theme.Secondary),  // Guac border
	}
}