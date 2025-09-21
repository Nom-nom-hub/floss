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

type flossStatusCmp struct {
	info       util.InfoMsg
	width      int
	messageTTL time.Duration
	help       help.Model
	keyMap     help.KeyMap
}

// clearMessageCmd is a command that clears status messages after a timeout
func (m *flossStatusCmp) clearMessageCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return util.ClearStatusMsg{}
	})
}

func (m *flossStatusCmp) Init() tea.Cmd {
	return nil
}

func (m *flossStatusCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, m.clearMessageCmd(ttl)
	case util.ClearStatusMsg:
		m.info = util.InfoMsg{}
	}
	return m, nil
}

func (m *flossStatusCmp) View() string {
	// Apply FLOSS-specific styling to the status bar
	status := getFlossStatusStyle().Padding(0, 1, 1, 1).Render(m.help.View(m.keyMap))
	if m.info.Msg != "" {
		status = m.infoMsg()
	}
	return status
}

// getFlossStatusStyle returns the FLOSS-specific status bar style
func getFlossStatusStyle() lipgloss.Style {
	theme := styles.CurrentTheme()
	
	return theme.S().Base.
		Background(theme.BgBase).           // Pepper background
		BorderTop(true).
		BorderForeground(theme.Secondary).  // Guac border
		Height(2)                           // Slightly taller
}

func (m *flossStatusCmp) infoMsg() string {
	t := styles.CurrentTheme()
	message := ""
	infoType := ""
	switch m.info.Type {
	case util.InfoTypeError:
		// Error: Cheeky background with Salt text according to FLOSS specification
		infoType = t.S().Base.Background(t.Error).Foreground(t.FgSelected).Padding(0, 1).Render("ERROR")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = t.S().Base.Background(t.Error).Width(widthLeft+2).Foreground(t.FgSelected).Padding(0, 1).Render(info)
	case util.InfoTypeWarn:
		// Warning: Zest background with BgOverlay text according to FLOSS specification
		infoType = t.S().Base.Foreground(t.BgOverlay).Background(t.Warning).Padding(0, 1).Render("WARNING")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = t.S().Base.Foreground(t.BgOverlay).Width(widthLeft+2).Background(t.Warning).Padding(0, 1).Render(info)
	default:
		// Info: Guac background with Salt text according to FLOSS specification
		infoType = t.S().Base.Background(t.Secondary).Foreground(t.FgSelected).Padding(0, 1).Render("OKAY!")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = t.S().Base.Background(t.Secondary).Width(widthLeft+2).Foreground(t.FgSelected).Padding(0, 1).Render(info)
	}
	return ansi.Truncate(infoType+message, m.width, "…")
}

func (m *flossStatusCmp) ToggleFullHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

func (m *flossStatusCmp) SetKeyMap(keyMap help.KeyMap) {
	m.keyMap = keyMap
}

// NewFlossStatusCmp creates a new FLOSS-styled status component
func NewFlossStatusCmp() StatusCmp {
	t := styles.CurrentTheme()
	help := help.New()
	help.Styles = t.S().Help
	return &flossStatusCmp{
		messageTTL: 5 * time.Second,
		help:       help,
	}
}