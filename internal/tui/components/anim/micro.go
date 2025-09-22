package anim

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
)

// MicroInteraction represents a subtle animation for user feedback
type MicroInteraction struct {
	id          int
	width       int
	label       string
	frames      []string
	current     int
	duration    time.Duration
	isAnimating bool
	style       lipgloss.Style
}

// MicroInteractionType defines different types of micro-interactions
type MicroInteractionType int

const (
	// ButtonPress represents a button press feedback animation
	ButtonPress MicroInteractionType = iota
	
	// Selection represents a selection feedback animation
	Selection
	
	// Loading represents a loading feedback animation
	Loading
	
	// Success represents a success feedback animation
	Success
	
	// Error represents an error feedback animation
	Error
)

// NewMicroInteraction creates a new micro-interaction with the specified type and settings
func NewMicroInteraction(interactionType MicroInteractionType, label string, duration time.Duration) *MicroInteraction {
	t := styles.CurrentTheme()
	m := &MicroInteraction{
		id:       nextID(),
		width:    20,
		label:    label,
		duration: duration,
		style:    t.S().Base,
	}
	
	// Set up frames based on interaction type
	switch interactionType {
	case ButtonPress:
		m.frames = []string{
			"▐",
			"▌",
			"█",
			"▌",
			"▐",
		}
	case Selection:
		m.frames = []string{
			"[    ]",
			"[=   ]",
			"[==  ]",
			"[=== ]",
			"[====]",
			"[ ===]",
			"[  ==]",
			"[   =]",
			"[    ]",
		}
	case Loading:
		m.frames = []string{
			"●",
			"●●",
			"●●●",
			"●●●●",
			"●●●●●",
			"●●●●",
			"●●●",
			"●●",
			"●",
		}
	case Success:
		m.frames = []string{
			"✓",
			"✓✓",
			"✓✓✓",
			"✓✓✓✓",
			"✓✓✓✓✓",
		}
	case Error:
		m.frames = []string{
			"✗",
			"✗✗",
			"✗✗✗",
			"✗✗✗✗",
			"✗✗✗✗✗",
		}
	}
	
	return m
}

// Init initializes the micro-interaction
func (m *MicroInteraction) Init() tea.Cmd {
	if len(m.frames) == 0 {
		return nil
	}
	m.isAnimating = true
	m.current = 0
	return m.tick()
}

// Update handles animation updates
func (m *MicroInteraction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StepMsg:
		if msg.id == m.id && m.isAnimating {
			m.current = (m.current + 1) % len(m.frames)
			if m.current == 0 {
				m.isAnimating = false
				return m, nil
			}
			return m, m.tick()
		}
	}
	return m, nil
}

// View renders the current frame of the micro-interaction
func (m *MicroInteraction) View() string {
	if !m.isAnimating || len(m.frames) == 0 {
		return ""
	}
	
	frame := m.frames[m.current]
	if m.label != "" {
		frame = m.label + " " + frame
	}
	
	return m.style.Render(frame)
}

// tick schedules the next animation frame
func (m *MicroInteraction) tick() tea.Cmd {
	return tea.Tick(m.duration/ time.Duration(len(m.frames)), func(t time.Time) tea.Msg {
		return StepMsg{id: m.id}
	})
}

// IsAnimating returns whether the micro-interaction is currently animating
func (m *MicroInteraction) IsAnimating() bool {
	return m.isAnimating
}

// SetStyle sets the style for the micro-interaction
func (m *MicroInteraction) SetStyle(style lipgloss.Style) {
	m.style = style
}

// ButtonPressFeedback creates a button press feedback animation
func ButtonPressFeedback() *MicroInteraction {
	return NewMicroInteraction(ButtonPress, "", 200*time.Millisecond)
}

// SelectionFeedback creates a selection feedback animation
func SelectionFeedback() *MicroInteraction {
	return NewMicroInteraction(Selection, "", 500*time.Millisecond)
}

// LoadingFeedback creates a loading feedback animation
func LoadingFeedback() *MicroInteraction {
	return NewMicroInteraction(Loading, "", 1*time.Second)
}

// SuccessFeedback creates a success feedback animation
func SuccessFeedback() *MicroInteraction {
	return NewMicroInteraction(Success, "", 500*time.Millisecond)
}

// ErrorFeedback creates an error feedback animation
func ErrorFeedback() *MicroInteraction {
	return NewMicroInteraction(Error, "", 500*time.Millisecond)
}