package status

import (
	"time"

	"github.com/charmbracelet/bubbles/v2/help"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
	"github.com/nom-nom-hub/floss/internal/tui/util"
)

type EnhancedStatusCmp interface {
	StatusCmp
	SetEnhancedMode(bool) tea.Cmd
	IsEnhancedMode() bool
}

type enhancedStatusCmp struct {
	info         util.InfoMsg
	width        int
	messageTTL   time.Duration
	help         help.Model
	keyMap       help.KeyMap
	enhancedMode bool
	animationStep int
	isAnimating  bool
}

// clearMessageCmd is a command that clears status messages after a timeout
func (m *enhancedStatusCmp) clearMessageCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return util.ClearStatusMsg{}
	})
}

func (m *enhancedStatusCmp) Init() tea.Cmd {
	return nil
}

func (m *enhancedStatusCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width - 2
		return m, nil

	// Handle status info
	case util.InfoMsg:
		m.info = msg
		ttl := msg.TTL
		if ttl == 0 {
			ttl = m.messageTTL
		}
		
		// Add animation for new messages
		m.isAnimating = true
		m.animationStep = 0
		animationCmd := tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			m.animationStep++
			if m.animationStep >= 10 {
				m.isAnimating = false
			}
			return nil
		})
		
		return m, tea.Sequence(animationCmd, m.clearMessageCmd(ttl))
	case util.ClearStatusMsg:
		m.info = util.InfoMsg{}
	}
	return m, nil
}

func (m *enhancedStatusCmp) View() string {
	t := styles.CurrentTheme()
	
	// Base styling with proper height and padding according to UI/UX specification
	// Height: 2 lines, Padding: 0 top/bottom, 1 left/right
	statusStyle := t.S().Base.Padding(0, 1, 1, 1)
	
	// Enhanced mode styling
	if m.enhancedMode {
		statusStyle = statusStyle.
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(t.Border)
	}
	
	status := statusStyle.Render(m.help.View(m.keyMap))
	if m.info.Msg != "" {
		status = m.infoMsg()
	}
	return status
}

func (m *enhancedStatusCmp) infoMsg() string {
	t := styles.CurrentTheme()
	message := ""
	infoType := ""
	
	switch m.info.Type {
	case util.InfoTypeError:
		// Error: Red background with white text according to UI/UX specification
		// Enhanced with better visual design
		errorStyle := t.S().Base.Background(t.Error).Foreground(t.White).Padding(0, 1)
		if m.isAnimating {
			errorStyle = errorStyle.Bold(true)
		}
		infoType = errorStyle.Render("ERROR")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		messageStyle := t.S().Base.Background(t.Error).Width(widthLeft+2).Foreground(t.White).Padding(0, 1)
		if m.isAnimating {
			messageStyle = messageStyle.Bold(true)
		}
		message = messageStyle.Render(info)
	case util.InfoTypeWarn:
		// Warning: Yellow background with BgOverlay text according to UI/UX specification
		// Enhanced with better visual design
		warnStyle := t.S().Base.Foreground(t.BgOverlay).Background(t.Warning).Padding(0, 1)
		if m.isAnimating {
			warnStyle = warnStyle.Bold(true)
		}
		infoType = warnStyle.Render("WARNING")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		messageStyle := t.S().Base.Foreground(t.BgOverlay).Width(widthLeft+2).Background(t.Warning).Padding(0, 1)
		if m.isAnimating {
			messageStyle = messageStyle.Bold(true)
		}
		message = messageStyle.Render(info)
	default:
		// Info: Green background with white text according to UI/UX specification
		// Enhanced with better visual design
		infoStyle := t.S().Base.Background(t.Success).Foreground(t.White).Padding(0, 1)
		if m.isAnimating {
			infoStyle = infoStyle.Bold(true)
		}
		infoType = infoStyle.Render("OKAY!")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		messageStyle := t.S().Base.Background(t.Success).Width(widthLeft+2).Foreground(t.White).Padding(0, 1)
		if m.isAnimating {
			messageStyle = messageStyle.Bold(true)
		}
		message = messageStyle.Render(info)
	}
	
	// Apply animation effect
	if m.isAnimating {
		// Add subtle glow effect
		glowColor := styles.Lighten(t.Primary, 0.3*float64(m.animationStep)/10.0)
		infoType = lipgloss.NewStyle().Foreground(glowColor).Render(infoType)
	}
	
	return ansi.Truncate(infoType+message, m.width, "…")
}

func (m *enhancedStatusCmp) ToggleFullHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

func (m *enhancedStatusCmp) SetKeyMap(keyMap help.KeyMap) {
	m.keyMap = keyMap
}

// Enhanced interface methods

func (m *enhancedStatusCmp) SetEnhancedMode(enabled bool) tea.Cmd {
	m.enhancedMode = enabled
	return nil
}

func (m *enhancedStatusCmp) IsEnhancedMode() bool {
	return m.enhancedMode
}

func NewEnhancedStatusCmp() EnhancedStatusCmp {
	t := styles.CurrentTheme()
	help := help.New()
	help.Styles = t.S().Help
	return &enhancedStatusCmp{
		messageTTL: 5 * time.Second,
		help:       help,
	}
}