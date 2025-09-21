package floss

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
)

// ButtonType defines the type of button
type ButtonType int

const (
	PrimaryButton ButtonType = iota
	SecondaryButton
	DangerButton
)

// Button represents a FLOSS-styled button component
type Button struct {
	Type     ButtonType
	Label    string
	Disabled bool
	Width    int
	Height   int
	onPress  func() tea.Msg
	keyMap   key.Binding
	style    lipgloss.Style
}

// NewButton creates a new FLOSS-styled button
func NewButton(buttonType ButtonType, label string, onPress func() tea.Msg) *Button {
	btn := &Button{
		Type:    buttonType,
		Label:   label,
		onPress: onPress,
		keyMap:  key.NewBinding(),
	}
	
	btn.updateStyle()
	return btn
}

// NewPrimaryButton creates a new primary button
func NewPrimaryButton(label string, onPress func() tea.Msg) *Button {
	return NewButton(PrimaryButton, label, onPress)
}

// NewSecondaryButton creates a new secondary button
func NewSecondaryButton(label string, onPress func() tea.Msg) *Button {
	return NewButton(SecondaryButton, label, onPress)
}

// NewDangerButton creates a new danger button
func NewDangerButton(label string, onPress func() tea.Msg) *Button {
	return NewButton(DangerButton, label, onPress)
}

// WithKeyBinding sets a key binding for the button
func (b *Button) WithKeyBinding(kb key.Binding) *Button {
	b.keyMap = kb
	return b
}

// WithSize sets the button size
func (b *Button) WithSize(width, height int) *Button {
	b.Width = width
	b.Height = height
	b.updateStyle()
	return b
}

// WithDisabled sets the button disabled state
func (b *Button) WithDisabled(disabled bool) *Button {
	b.Disabled = disabled
	b.updateStyle()
	return b
}

// updateStyle updates the button style based on its type and state
func (b *Button) updateStyle() {
	buttonStyles := GetButtonStyles()
	
	switch b.Type {
	case PrimaryButton:
		b.style = buttonStyles.Primary
	case SecondaryButton:
		b.style = buttonStyles.Secondary
	case DangerButton:
		b.style = buttonStyles.Danger
	}
	
	// Apply disabled styling if needed
	if b.Disabled {
		theme := styles.CurrentTheme()
		switch b.Type {
		case PrimaryButton:
			b.style = b.style.Background(theme.FgHalfMuted) // Smoke background
		case SecondaryButton:
			b.style = b.style.Foreground(theme.FgHalfMuted).BorderForeground(theme.FgHalfMuted)
		case DangerButton:
			b.style = b.style.Background(theme.FgHalfMuted).BorderForeground(theme.FgHalfMuted)
		}
		b.style = b.style.UnsetString()
	}
	
	// Apply size if set
	if b.Width > 0 {
		b.style = b.style.Width(b.Width)
	}
	if b.Height > 0 {
		b.style = b.style.Height(b.Height)
	}
}

// Init implements tea.Model
func (b *Button) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !b.Disabled && key.Matches(msg, b.keyMap) {
			if b.onPress != nil {
				return b, func() tea.Msg { return b.onPress() }
			}
		}
	}
	return b, nil
}

// View implements tea.Model
func (b *Button) View() string {
	return b.style.Render(b.Label)
}

// Bindings returns the button's key bindings
func (b *Button) Bindings() []key.Binding {
	if b.keyMap.Enabled() {
		return []key.Binding{b.keyMap}
	}
	return []key.Binding{}
}

// IsDisabled returns whether the button is disabled
func (b *Button) IsDisabled() bool {
	return b.Disabled
}

// SetDisabled sets the button's disabled state
func (b *Button) SetDisabled(disabled bool) {
	b.Disabled = disabled
	b.updateStyle()
}