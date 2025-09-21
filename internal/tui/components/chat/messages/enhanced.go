package messages

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/google/uuid"

	"github.com/atotto/clipboard"
	"github.com/nom-nom-hub/floss/internal/config"
	"github.com/nom-nom-hub/floss/internal/message"
	"github.com/nom-nom-hub/floss/internal/tui/components/anim"
	"github.com/nom-nom-hub/floss/internal/tui/components/core/layout"
	"github.com/nom-nom-hub/floss/internal/tui/exp/list"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
	"github.com/nom-nom-hub/floss/internal/tui/util"
	"github.com/nom-nom-hub/floss/internal/tui/util/section"
)

// EnhancedMessageCmp defines the interface for enhanced message components in the chat interface.
// It extends the standard MessageCmp with additional functionality for enhanced UI/UX.
type EnhancedMessageCmp interface {
	MessageCmp
	ToggleDetails() tea.Cmd
	ToggleExpand() tea.Cmd
	IsExpanded() bool
	HasDetails() bool
}

// enhancedMessageCmp implements the EnhancedMessageCmp interface for displaying enhanced chat messages.
// It extends the standard message component with collapsible sections, detailed metadata,
// animated transitions, and interactive elements.
type enhancedMessageCmp struct {
	width   int  // Component width for text wrapping
	height  int  // Component height
	focused bool // Focus state for border styling

	// Core message data and state
	message  message.Message // The underlying message content
	spinning bool            // Whether to show loading animation
	anim     *anim.Anim      // Animation component for loading states

	// Enhanced UI features
	showDetails bool // Whether to show detailed metadata
	expanded    bool // Whether message content is expanded

	// Thinking viewport for displaying reasoning content
	thinkingViewport viewport.Model
	
	// Animation states for transitions
	detailsAnimating bool
	expandAnimating  bool
	animationStep    int
}

// NewEnhancedMessageCmp creates a new enhanced message component with the given message and options
func NewEnhancedMessageCmp(msg message.Message) EnhancedMessageCmp {
	t := styles.CurrentTheme()

	thinkingViewport := viewport.New()
	thinkingViewport.SetHeight(1)
	thinkingViewport.KeyMap = viewport.KeyMap{}

	m := &enhancedMessageCmp{
		message: msg,
		anim: anim.New(anim.Settings{
			Size:        15,
			GradColorA:  t.Primary,
			GradColorB:  t.Secondary,
			CycleColors: true,
		}),
		thinkingViewport: thinkingViewport,
		showDetails:      false,
		expanded:         false,
	}
	return m
}

// Init initializes the enhanced message component and starts animations if needed.
// Returns a command to start the animation for spinning messages.
func (m *enhancedMessageCmp) Init() tea.Cmd {
	m.spinning = m.shouldSpin()
	return m.anim.Init()
}

// Update handles incoming messages and updates the component state.
// Manages animation updates for spinning messages and stops animation when appropriate.
func (m *enhancedMessageCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case anim.StepMsg:
		m.spinning = m.shouldSpin()
		if m.spinning {
			u, cmd := m.anim.Update(msg)
			m.anim = u.(*anim.Anim)
			return m, cmd
		}
	case tea.KeyPressMsg:
		if key.Matches(msg, EnhancedCopyKey) {
			return m, tea.Sequence(
				tea.SetClipboard(m.message.Content().Text),
				func() tea.Msg {
					_ = clipboard.WriteAll(m.message.Content().Text)
					return nil
				},
				util.ReportInfo("Message copied to clipboard"),
			)
		}
	}
	return m, nil
}

// View renders the enhanced message component based on its current state.
// Returns different views for spinning, user, and assistant messages with enhanced styling.
func (m *enhancedMessageCmp) View() string {
	if m.spinning && m.message.ReasoningContent().Thinking == "" {
		return m.style().PaddingLeft(1).Render(m.anim.View())
	}
	if m.message.ID != "" {
		// this is a user or assistant message
		switch m.message.Role {
		case message.User:
			return m.renderUserMessage()
		default:
			return m.renderAssistantMessage()
		}
	}
	return m.style().Render("No message content")
}

// GetMessage returns the underlying message data
func (m *enhancedMessageCmp) GetMessage() message.Message {
	return m.message
}

func (m *enhancedMessageCmp) SetMessage(msg message.Message) {
	m.message = msg
}

// textWidth calculates the available width for text content,
// accounting for borders and padding
func (m *enhancedMessageCmp) textWidth() int {
	padding := 2
	if m.expanded {
		padding = 4 // account for all borders when expanded
	} else if m.message.Role == message.User {
		padding = 3 // account for left border and padding
	}
	return m.width - padding
}

// style returns the lipgloss style for the enhanced message component.
// Applies different border colors and styles based on message role, focus state, and expansion.
func (msg *enhancedMessageCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	borderStyle := lipgloss.NormalBorder()
	
	if msg.expanded {
		borderStyle = lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "┌",
			TopRight:    "┐",
			BottomLeft:  "└",
			BottomRight: "┘",
		}
	} else if msg.focused {
		borderStyle = lipgloss.Border{
			Left: "▌",
		}
	}

	style := t.S().Text
	if msg.message.Role == message.User {
		// User messages: Left border with Primary color and 1 unit padding
		if msg.expanded {
			style = style.Padding(1, 1, 1, 1).Border(borderStyle).BorderStyle(borderStyle).BorderForeground(t.Primary)
		} else {
			style = style.PaddingLeft(1).BorderLeft(true).BorderStyle(borderStyle).BorderForeground(t.Primary)
		}
	} else {
		// Assistant messages: Left border with GreenDark when focused, 2 units padding when not focused
		if msg.expanded {
			style = style.Padding(1, 1, 1, 1).Border(borderStyle).BorderStyle(borderStyle).BorderForeground(t.GreenDark)
		} else if msg.focused {
			style = style.PaddingLeft(1).BorderLeft(true).BorderStyle(borderStyle).BorderForeground(t.GreenDark)
		} else {
			style = style.PaddingLeft(2)
		}
	}
	
	// Add subtle background for expanded messages
	if msg.expanded {
		style = style.Background(t.BgBaseLighter)
	}
	
	return style
}

// renderAssistantMessage renders assistant messages with optional footer information and enhanced styling.
// Shows model name, response time, and finish reason when the message is complete.
func (m *enhancedMessageCmp) renderAssistantMessage() string {
	t := styles.CurrentTheme()
	parts := []string{}
	content := m.message.Content().String()
	thinking := m.message.IsThinking()
	finished := m.message.IsFinished()
	finishedData := m.message.FinishPart()
	thinkingContent := ""

	if thinking || m.message.ReasoningContent().Thinking != "" {
		m.anim.SetLabel("Thinking")
		thinkingContent = m.renderThinkingContent()
	} else if finished && content == "" && finishedData.Reason == message.FinishReasonEndTurn {
		content = ""
	} else if finished && content == "" && finishedData.Reason == message.FinishReasonCanceled {
		content = "*Canceled*"
	} else if finished && content == "" && finishedData.Reason == message.FinishReasonError {
		// Error display with appropriate styling according to UI/UX specification
		errTag := t.S().Base.Padding(0, 1).Background(t.Error).Foreground(t.White).Render("ERROR")
		truncated := ansi.Truncate(finishedData.Message, m.textWidth()-2-lipgloss.Width(errTag), "...")
		title := fmt.Sprintf("%s %s", errTag, t.S().Base.Foreground(t.FgHalfMuted).Render(truncated))
		details := t.S().Base.Foreground(t.FgSubtle).Width(m.textWidth() - 2).Render(finishedData.Details)
		errorContent := fmt.Sprintf("%s\n\n%s", title, details)
		return m.style().Render(errorContent)
	}

	if thinkingContent != "" {
		parts = append(parts, thinkingContent)
	}

	if content != "" {
		if thinkingContent != "" {
			parts = append(parts, "")
		}
		parts = append(parts, m.toMarkdown(content))
	}

	// Add metadata footer if details are shown
	if m.showDetails {
		footer := m.renderMetadataFooter()
		if footer != "" {
			parts = append(parts, "", footer)
		}
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.style().Render(joined)
}

// renderUserMessage renders user messages with file attachments and enhanced styling.
// Displays message content and any attached files with appropriate icons.
func (m *enhancedMessageCmp) renderUserMessage() string {
	t := styles.CurrentTheme()
	parts := []string{
		m.toMarkdown(m.message.Content().String()),
	}

	// Displayed as small badges with DocumentIcon according to UI/UX specification
	attachmentStyles := t.S().Text.
		MarginLeft(1).
		Background(t.FgMuted).  // Use FgMuted background for attachment badges
		Foreground(t.FgBase)    // Use FgBase text for contrast

	attachments := make([]string, len(m.message.BinaryContent()))
	for i, attachment := range m.message.BinaryContent() {
		const maxFilenameWidth = 10
		filename := filepath.Base(attachment.Path)
		attachments[i] = attachmentStyles.Render(fmt.Sprintf(
			" %s %s ",
			styles.DocumentIcon,
			ansi.Truncate(filename, maxFilenameWidth, "..."),
		))
	}

	if len(attachments) > 0 {
		parts = append(parts, "", strings.Join(attachments, ""))
	}

	// Add metadata footer if details are shown
	if m.showDetails {
		footer := m.renderMetadataFooter()
		if footer != "" {
			parts = append(parts, "", footer)
		}
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.style().Render(joined)
}

// toMarkdown converts text content to rendered markdown using the configured renderer
func (m *enhancedMessageCmp) toMarkdown(content string) string {
	r := styles.GetMarkdownRenderer(m.textWidth())
	rendered, _ := r.Render(content)
	return strings.TrimSuffix(rendered, "\n")
}

func (m *enhancedMessageCmp) renderThinkingContent() string {
	t := styles.CurrentTheme()
	reasoningContent := m.message.ReasoningContent()
	if reasoningContent.Thinking == "" {
		return ""
	}
	lines := strings.Split(reasoningContent.Thinking, "\n")
	var content strings.Builder
	// Rendered in BgBaseLighter with subtle styling according to UI/UX specification
	lineStyle := t.S().Base.Background(t.BgBaseLighter).Foreground(t.FgBase)
	for i, line := range lines {
		if line == "" {
			continue
		}
		content.WriteString(lineStyle.Width(m.textWidth() - 2).Render(line))
		if i < len(lines)-1 {
			content.WriteString("\n")
		}
	}
	fullContent := content.String()
	height := util.Clamp(lipgloss.Height(fullContent), 1, 10)
	m.thinkingViewport.SetHeight(height)
	m.thinkingViewport.SetWidth(m.textWidth())
	m.thinkingViewport.SetContent(fullContent)
	m.thinkingViewport.GotoBottom()
	finishReason := m.message.FinishPart()
	var footer string
	if reasoningContent.StartedAt > 0 {
		duration := m.message.ThinkingDuration()
		if reasoningContent.FinishedAt > 0 {
			m.anim.SetLabel("")
			opts := section.StatusOpts{
				Title:       "Thought for",
				Description: duration.String(),
			}
			if duration.String() != "0s" {
				footer = t.S().Base.PaddingLeft(1).Render(section.Status(opts, m.textWidth()-1))
			}
		} else if finishReason != nil && finishReason.Reason == message.FinishReasonCanceled {
			footer = t.S().Base.PaddingLeft(1).Render(m.toMarkdown("*Canceled*"))
		} else {
			footer = m.anim.View()
		}
	}
	return lineStyle.Width(m.textWidth()).Padding(0, 1).Render(m.thinkingViewport.View()) + "\n\n" + footer
}

// renderMetadataFooter renders detailed metadata about the message
func (m *enhancedMessageCmp) renderMetadataFooter() string {
	t := styles.CurrentTheme()
	
	// Only show metadata for assistant messages
	if m.message.Role != message.Assistant {
		return ""
	}
	
	finishData := m.message.FinishPart()
	if finishData == nil {
		return ""
	}
	
	// Format timestamp
	timestamp := time.Unix(finishData.Time, 0).Format("15:04:05")
	
	// Get model info
	model := config.Get().GetModel(m.message.Provider, m.message.Model)
	modelName := "Unknown Model"
	if model != nil {
		modelName = model.Name
	}
	
	// Format token usage
	var tokenInfo string
	if finishData.Message != "" {
		tokenInfo = fmt.Sprintf("Message: %s", finishData.Message)
	}
	
	// Create metadata lines
	metadataLines := []string{
		fmt.Sprintf("Model: %s", modelName),
		fmt.Sprintf("Timestamp: %s", timestamp),
	}
	
	if tokenInfo != "" {
		metadataLines = append(metadataLines, fmt.Sprintf("Details: %s", tokenInfo))
	}
	
	if finishData.Reason != "" {
		metadataLines = append(metadataLines, fmt.Sprintf("Finish Reason: %s", finishData.Reason))
	}
	
	// Style the metadata
	metadataStyle := t.S().Base.Foreground(t.FgSubtle).Padding(0, 1)
	metadataContent := strings.Join(metadataLines, "  •  ")
	
	return metadataStyle.Render(metadataContent)
}

// shouldSpin determines whether the message should show a loading animation.
// Only assistant messages without content that aren't finished should spin.
func (m *enhancedMessageCmp) shouldSpin() bool {
	if m.message.Role != message.Assistant {
		return false
	}

	if m.message.IsFinished() {
		return false
	}

	if m.message.Content().Text != "" {
		return false
	}
	if len(m.message.ToolCalls()) > 0 {
		return false
	}
	return true
}

// Blur removes focus from the message component
func (m *enhancedMessageCmp) Blur() tea.Cmd {
	m.focused = false
	return nil
}

// Focus sets focus on the message component
func (m *enhancedMessageCmp) Focus() tea.Cmd {
	m.focused = true
	return nil
}

// IsFocused returns whether the message component is currently focused
func (m *enhancedMessageCmp) IsFocused() bool {
	return m.focused
}

// Size management methods

// GetSize returns the current dimensions of the message component
func (m *enhancedMessageCmp) GetSize() (int, int) {
	return m.width, m.height
}

// SetSize updates the width of the message component for text wrapping
func (m *enhancedMessageCmp) SetSize(width int, height int) tea.Cmd {
	m.width = util.Clamp(width, 1, 120)
	m.height = height
	m.thinkingViewport.SetWidth(m.width - 4)
	return nil
}

// Spinning returns whether the message is currently showing a loading animation
func (m *enhancedMessageCmp) Spinning() bool {
	return m.spinning
}

// Enhanced interface methods

// ToggleDetails toggles the visibility of detailed metadata
func (m *enhancedMessageCmp) ToggleDetails() tea.Cmd {
	m.showDetails = !m.showDetails
	// Add subtle animation for toggling details
	m.detailsAnimating = true
	m.animationStep = 0
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		m.animationStep++
		if m.animationStep >= 5 {
			m.detailsAnimating = false
		}
		return nil
	})
}

// ToggleExpand toggles the expanded state of the message
func (m *enhancedMessageCmp) ToggleExpand() tea.Cmd {
	m.expanded = !m.expanded
	// Add subtle animation for expanding/collapsing
	m.expandAnimating = true
	m.animationStep = 0
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		m.animationStep++
		if m.animationStep >= 5 {
			m.expandAnimating = false
		}
		return nil
	})
}

// IsExpanded returns whether the message is currently expanded
func (m *enhancedMessageCmp) IsExpanded() bool {
	return m.expanded
}

// HasDetails returns whether the message has details to show
func (m *enhancedMessageCmp) HasDetails() bool {
	// Assistant messages with finish data have details
	if m.message.Role == message.Assistant {
		return m.message.FinishPart() != nil
	}
	return false
}

func (m *enhancedMessageCmp) ID() string {
	return m.message.ID
}

// EnhancedAssistantSection represents an enhanced section header for assistant messages
type EnhancedAssistantSection interface {
	list.Item
	layout.Sizeable
	IsSectionHeader() bool
}

type enhancedAssistantSectionModel struct {
	width               int
	id                  string
	message             message.Message
	lastUserMessageTime time.Time
}

// ID implements EnhancedAssistantSection.
func (m *enhancedAssistantSectionModel) ID() string {
	return m.id
}

func NewEnhancedAssistantSection(message message.Message, lastUserMessageTime time.Time) EnhancedAssistantSection {
	return &enhancedAssistantSectionModel{
		width:               0,
		id:                  uuid.NewString(),
		message:             message,
		lastUserMessageTime: lastUserMessageTime,
	}
}

func (m *enhancedAssistantSectionModel) Init() tea.Cmd {
	return nil
}

func (m *enhancedAssistantSectionModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *enhancedAssistantSectionModel) View() string {
	t := styles.CurrentTheme()
	finishData := m.message.FinishPart()
	finishTime := time.Unix(finishData.Time, 0)
	duration := finishTime.Sub(m.lastUserMessageTime)
	infoMsg := t.S().Subtle.Render(duration.String())
	icon := t.S().Subtle.Render(styles.ModelIcon)
	model := config.Get().GetModel(m.message.Provider, m.message.Model)
	if model == nil {
		// This means the model is not configured anymore
		model = &catwalk.Model{
			Name: "Unknown Model",
		}
	}
	modelFormatted := t.S().Muted.Render(model.Name)
	assistant := fmt.Sprintf("%s %s %s", icon, modelFormatted, infoMsg)
	return t.S().Base.PaddingLeft(2).Render(
		section.Section(assistant, m.width-2),
	)
}

func (m *enhancedAssistantSectionModel) GetSize() (int, int) {
	return m.width, 1
}

func (m *enhancedAssistantSectionModel) SetSize(width int, height int) tea.Cmd {
	m.width = width
	return nil
}

func (m *enhancedAssistantSectionModel) IsSectionHeader() bool {
	return true
}