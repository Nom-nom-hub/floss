package editor

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/nom-nom-hub/floss/internal/app"
	"github.com/nom-nom-hub/floss/internal/fsext"
	"github.com/nom-nom-hub/floss/internal/message"
	"github.com/nom-nom-hub/floss/internal/session"
	"github.com/nom-nom-hub/floss/internal/tui/components/chat"
	"github.com/nom-nom-hub/floss/internal/tui/components/completions"
	"github.com/nom-nom-hub/floss/internal/tui/components/dialogs"
	"github.com/nom-nom-hub/floss/internal/tui/components/dialogs/commands"
	"github.com/nom-nom-hub/floss/internal/tui/components/dialogs/filepicker"
	"github.com/nom-nom-hub/floss/internal/tui/components/dialogs/quit"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
	"github.com/nom-nom-hub/floss/internal/tui/util"
)

type enhancedEditorCmp struct {
	width              int
	height             int
	x, y               int
	app                *app.App
	session            session.Session
	textarea           *textarea.Model
	attachments        []message.Attachment
	deleteMode         bool
	readyPlaceholder   string
	workingPlaceholder string
	enhancedMode       bool // New: Enhanced mode flag

	keyMap EditorKeyMap

	// File path completions
	currentQuery          string
	completionsStartIndex int
	isCompletionsOpen     bool
	
	// Animation states
	attachmentAnimationStep int
	isAnimating             bool
}

var EnhancedDeleteKeyMaps = DeleteAttachmentKeyMaps{
	AttachmentDeleteMode: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r+{i}", "delete attachment at index i"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel delete mode"),
	),
	DeleteAllAttachments: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("ctrl+r+r", "delete all attachments"),
	),
}

func (m *enhancedEditorCmp) openEditor(value string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Use platform-appropriate default editor
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "nvim"
		}
	}

	tmpfile, err := os.CreateTemp("", "msg_*.md")
	if err != nil {
		return util.ReportError(err)
	}
	defer tmpfile.Close() //nolint:errcheck
	if _, err := tmpfile.WriteString(value); err != nil {
		return util.ReportError(err)
	}
	c := exec.CommandContext(context.TODO(), editor, tmpfile.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return util.ReportError(err)
		}
		content, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return util.ReportError(err)
		}
		if len(content) == 0 {
			return util.ReportWarn("Message is empty")
		}
		os.Remove(tmpfile.Name())
		return OpenEditorMsg{
			Text: strings.TrimSpace(string(content)),
		}
	})
}

func (m *enhancedEditorCmp) Init() tea.Cmd {
	return nil
}

func (m *enhancedEditorCmp) send() tea.Cmd {
	value := m.textarea.Value()
	value = strings.TrimSpace(value)

	switch value {
	case "exit", "quit":
		m.textarea.Reset()
		return util.CmdHandler(dialogs.OpenDialogMsg{Model: quit.NewQuitDialog()})
	}

	m.textarea.Reset()
	attachments := m.attachments

	m.attachments = nil
	if value == "" {
		return nil
	}

	// Change the placeholder when sending a new message.
	m.randomizePlaceholders()

	return tea.Batch(
		util.CmdHandler(chat.SendMsg{
			Text:        value,
			Attachments: attachments,
		}),
	)
}

func (m *enhancedEditorCmp) repositionCompletions() tea.Msg {
	x, y := m.completionsPosition()
	return completions.RepositionCompletionsMsg{X: x, Y: y}
}

func (m *enhancedEditorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, m.repositionCompletions
	case filepicker.FilePickedMsg:
		if len(m.attachments) >= maxAttachments {
			return m, util.ReportError(fmt.Errorf("cannot add more than %d images", maxAttachments))
		}
		m.attachments = append(m.attachments, msg.Attachment)
		
		// Animate attachment addition
		m.isAnimating = true
		m.attachmentAnimationStep = 0
		cmds = append(cmds, tea.Tick(time.Millisecond*30, func(t time.Time) tea.Msg {
			m.attachmentAnimationStep++
			if m.attachmentAnimationStep >= 10 {
				m.isAnimating = false
			}
			return nil
		}))
		
		return m, tea.Batch(cmds...)
	case completions.CompletionsOpenedMsg:
		m.isCompletionsOpen = true
	case completions.CompletionsClosedMsg:
		m.isCompletionsOpen = false
		m.currentQuery = ""
		m.completionsStartIndex = 0
	case completions.SelectCompletionMsg:
		if !m.isCompletionsOpen {
			return m, nil
		}
		if item, ok := msg.Value.(FileCompletionItem); ok {
			word := m.textarea.Word()
			// If the selected item is a file, insert its path into the textarea
			value := m.textarea.Value()
			value = value[:m.completionsStartIndex] + // Remove the current query
				item.Path + // Insert the file path
				value[m.completionsStartIndex+len(word):] // Append the rest of the value
			// XXX: This will always move the cursor to the end of the textarea.
			m.textarea.SetValue(value)
			m.textarea.MoveToEnd()
			if !msg.Insert {
				m.isCompletionsOpen = false
				m.currentQuery = ""
				m.completionsStartIndex = 0
			}
		}

	case commands.OpenExternalEditorMsg:
		if m.app.CoderAgent.IsSessionBusy(m.session.ID) {
			return m, util.ReportWarn("Agent is working, please wait...")
		}
		return m, m.openEditor(m.textarea.Value())
	case OpenEditorMsg:
		m.textarea.SetValue(msg.Text)
		m.textarea.MoveToEnd()
	case tea.PasteMsg:
		path := strings.ReplaceAll(string(msg), "\\ ", " ")
		// try to get an image
		path, err := filepath.Abs(strings.TrimSpace(path))
		if err != nil {
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
		isAllowedType := false
		for _, ext := range filepicker.AllowedTypes {
			if strings.HasSuffix(path, ext) {
				isAllowedType = true
				break
			}
		}
		if !isAllowedType {
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
		tooBig, _ := filepicker.IsFileTooBig(path, filepicker.MaxAttachmentSize)
		if tooBig {
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

		content, err := os.ReadFile(path)
		if err != nil {
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
		mimeBufferSize := min(512, len(content))
		mimeType := http.DetectContentType(content[:mimeBufferSize])
		fileName := filepath.Base(path)
		attachment := message.Attachment{FilePath: path, FileName: fileName, MimeType: mimeType, Content: content}
		return m, util.CmdHandler(filepicker.FilePickedMsg{
			Attachment: attachment,
		})

	case commands.ToggleYoloModeMsg:
		m.setEditorPrompt()
		return m, nil
	case tea.KeyPressMsg:
		cur := m.textarea.Cursor()
		curIdx := m.textarea.Width()*cur.Y + cur.X
		switch {
		// Completions
		case msg.String() == "/" && !m.isCompletionsOpen &&
			// only show if beginning of prompt, or if previous char is a space or newline:
			(len(m.textarea.Value()) == 0 || unicode.IsSpace(rune(m.textarea.Value()[len(m.textarea.Value())-1]))):
			m.isCompletionsOpen = true
			m.currentQuery = ""
			m.completionsStartIndex = curIdx
			cmds = append(cmds, m.startCompletions)
		case m.isCompletionsOpen && curIdx <= m.completionsStartIndex:
			cmds = append(cmds, util.CmdHandler(completions.CloseCompletionsMsg{}))
		}
		if key.Matches(msg, EnhancedDeleteKeyMaps.AttachmentDeleteMode) {
			m.deleteMode = true
			return m, nil
		}
		if key.Matches(msg, EnhancedDeleteKeyMaps.DeleteAllAttachments) && m.deleteMode {
			m.deleteMode = false
			m.attachments = nil
			return m, nil
		}
		rune := msg.Code
		if m.deleteMode && unicode.IsDigit(rune) {
			num := int(rune - '0')
			m.deleteMode = false
			if num < 10 && len(m.attachments) > num {
				if num == 0 {
					m.attachments = m.attachments[num+1:]
				} else {
					m.attachments = slices.Delete(m.attachments, num, num+1)
				}
				return m, nil
			}
		}
		if key.Matches(msg, m.keyMap.OpenEditor) {
			if m.app.CoderAgent.IsSessionBusy(m.session.ID) {
				return m, util.ReportWarn("Agent is working, please wait...")
			}
			return m, m.openEditor(m.textarea.Value())
		}
		if key.Matches(msg, EnhancedDeleteKeyMaps.Escape) {
			m.deleteMode = false
			return m, nil
		}
		if key.Matches(msg, m.keyMap.Newline) {
			m.textarea.InsertRune('\n')
			cmds = append(cmds, util.CmdHandler(completions.CloseCompletionsMsg{}))
		}
		// Handle Enter key
		if m.textarea.Focused() && key.Matches(msg, m.keyMap.SendMessage) {
			value := m.textarea.Value()
			if strings.HasSuffix(value, "\\") {
				// If the last character is a backslash, remove it and add a newline.
				m.textarea.SetValue(strings.TrimSuffix(value, "\\"))
			} else {
				// Otherwise, send the message
				return m, m.send()
			}
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	if m.textarea.Focused() {
		kp, ok := msg.(tea.KeyPressMsg)
		if ok {
			if kp.String() == "space" || m.textarea.Value() == "" {
				m.isCompletionsOpen = false
				m.currentQuery = ""
				m.completionsStartIndex = 0
				cmds = append(cmds, util.CmdHandler(completions.CloseCompletionsMsg{}))
			} else {
				word := m.textarea.Word()
				if strings.HasPrefix(word, "/") {
					// XXX: wont' work if editing in the middle of the field.
					m.completionsStartIndex = strings.LastIndex(m.textarea.Value(), word)
					m.currentQuery = word[1:]
					x, y := m.completionsPosition()
					x -= len(m.currentQuery)
					m.isCompletionsOpen = true
					cmds = append(cmds,
						util.CmdHandler(completions.FilterCompletionsMsg{
							Query:  m.currentQuery,
							Reopen: m.isCompletionsOpen,
							X:      x,
							Y:      y,
						}),
					)
				} else if m.isCompletionsOpen {
					m.isCompletionsOpen = false
					m.currentQuery = ""
					m.completionsStartIndex = 0
					cmds = append(cmds, util.CmdHandler(completions.CloseCompletionsMsg{}))
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *enhancedEditorCmp) setEditorPrompt() {
	if m.app.Permissions.SkipRequests() {
		m.textarea.SetPromptFunc(4, enhancedYoloPromptFunc)
		return
	}
	m.textarea.SetPromptFunc(4, enhancedNormalPromptFunc)
}

func (m *enhancedEditorCmp) completionsPosition() (int, int) {
	cur := m.textarea.Cursor()
	if cur == nil {
		return m.x, m.y + 1 // adjust for padding
	}
	x := cur.X + m.x
	y := cur.Y + m.y + 1 // adjust for padding
	return x, y
}

func (m *enhancedEditorCmp) Cursor() *tea.Cursor {
	cursor := m.textarea.Cursor()
	if cursor != nil {
		cursor.X = cursor.X + m.x + 1
		cursor.Y = cursor.Y + m.y + 1 // adjust for padding
	}
	return cursor
}

var enhancedReadyPlaceholders = [...]string{
	"Ready for your input!",
	"Ready to assist...",
	"Type your message...",
	"What can I help with?",
}

var enhancedWorkingPlaceholders = [...]string{
	"Processing your request...",
	"Thinking...",
	"Working on it...",
	"Brrrrr... computing...",
	"Analyzing...",
}

func (m *enhancedEditorCmp) randomizePlaceholders() {
	m.workingPlaceholder = enhancedWorkingPlaceholders[rand.Intn(len(enhancedWorkingPlaceholders))]
	m.readyPlaceholder = enhancedReadyPlaceholders[rand.Intn(len(enhancedReadyPlaceholders))]
}

func (m *enhancedEditorCmp) View() string {
	t := styles.CurrentTheme()
	
	// Update placeholder with more descriptive text
	if m.app.CoderAgent != nil && m.app.CoderAgent.IsBusy() {
		m.textarea.Placeholder = m.workingPlaceholder
	} else {
		m.textarea.Placeholder = m.readyPlaceholder
	}
	if m.app.Permissions.SkipRequests() {
		m.textarea.Placeholder = "Yolo mode enabled!"
	}
	
	// Enhanced styling with border and background
	baseStyle := t.S().Base
	
	// Add subtle border to editor
	editorStyle := baseStyle.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Background(t.BgBaseLighter)
	
	if m.enhancedMode {
		// Enhanced mode: Add gradient border
		editorStyle = editorStyle.
			BorderForeground(t.Primary).
			Background(t.BgBase)
	}
	
	if len(m.attachments) == 0 {
		// No attachments: 1 unit padding all around
		content := editorStyle.Padding(1).Render(m.textarea.View())
		return content
	}
	
	// With attachments: 1 unit top/bottom padding, 1 unit left/right padding
	content := editorStyle.Padding(0, 1, 1, 1).Render(
		lipgloss.JoinVertical(lipgloss.Top,
			m.attachmentsContent(),
			m.textarea.View(),
		),
	)
	return content
}

func (m *enhancedEditorCmp) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	// Ensure minimum size constraints
	textWidth := max(10, width-4)   // adjust for padding and borders
	textHeight := max(3, height-4)  // adjust for padding and borders
	m.textarea.SetWidth(textWidth)
	m.textarea.SetHeight(textHeight)
	return nil
}

func (m *enhancedEditorCmp) GetSize() (int, int) {
	return m.textarea.Width(), m.textarea.Height()
}

func (m *enhancedEditorCmp) attachmentsContent() string {
	var styledAttachments []string
	t := styles.CurrentTheme()
	
	// Enhanced attachment styling with better visual design
	attachmentStyles := t.S().Base.
		MarginLeft(1).
		Padding(0, 1).
		Background(t.FgMuted).
		Foreground(t.FgBase).
		Bold(true)
		
	// Add rounded corners for enhanced look
	attachmentStyles = attachmentStyles.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.FgSubtle)
	
	// Animate if adding attachments
	if m.isAnimating {
		// Add pulsing effect during animation
		animationIntensity := float64(m.attachmentAnimationStep) / 10.0
		attachmentStyles = attachmentStyles.
			Background(styles.Lighten(t.FgMuted, animationIntensity*30))
	}
	
	for i, attachment := range m.attachments {
		var filename string
		if len(attachment.FileName) > 15 {
			filename = fmt.Sprintf(" %s %s...", styles.DocumentIcon, attachment.FileName[0:12])
		} else {
			filename = fmt.Sprintf(" %s %s", styles.DocumentIcon, attachment.FileName)
		}
		
		// Delete Mode: Number prefixes for quick deletion according to UI/UX specification
		if m.deleteMode {
			filename = fmt.Sprintf("%d%s", i, filename)
		}
		styledAttachments = append(styledAttachments, attachmentStyles.Render(filename))
	}
	
	// Add subtle shadow effect for attachments
	attachmentContainer := t.S().Base.
		Padding(0, 0, 1, 0).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, styledAttachments...))
		
	return attachmentContainer
}

func (m *enhancedEditorCmp) SetPosition(x, y int) tea.Cmd {
	m.x = x
	m.y = y
	return nil
}

func (m *enhancedEditorCmp) startCompletions() tea.Msg {
	files, _, _ := fsext.ListDirectory(".", nil, 0)
	slices.Sort(files)
	completionItems := make([]completions.Completion, 0, len(files))
	for _, file := range files {
		file = strings.TrimPrefix(file, "./")
		completionItems = append(completionItems, completions.Completion{
			Title: file,
			Value: FileCompletionItem{
				Path: file,
			},
		})
	}

	x, y := m.completionsPosition()
	return completions.OpenCompletionsMsg{
		Completions: completionItems,
		X:           x,
		Y:           y,
	}
}

// Blur implements Container.
func (c *enhancedEditorCmp) Blur() tea.Cmd {
	c.textarea.Blur()
	return nil
}

// Focus implements Container.
func (c *enhancedEditorCmp) Focus() tea.Cmd {
	return c.textarea.Focus()
}

// IsFocused implements Container.
func (c *enhancedEditorCmp) IsFocused() bool {
	return c.textarea.Focused()
}

// Bindings implements Container.
func (c *enhancedEditorCmp) Bindings() []key.Binding {
	return c.keyMap.KeyBindings()
}

// TODO: most likely we do not need to have the session here
// we need to move some functionality to the page level
func (c *enhancedEditorCmp) SetSession(session session.Session) tea.Cmd {
	c.session = session
	return nil
}

func (c *enhancedEditorCmp) IsCompletionsOpen() bool {
	return c.isCompletionsOpen
}

func (c *enhancedEditorCmp) HasAttachments() bool {
	return len(c.attachments) > 0
}

// Enhanced interface methods

func (c *enhancedEditorCmp) ToggleEnhancedMode() tea.Cmd {
	c.enhancedMode = !c.enhancedMode
	return nil
}

func (c *enhancedEditorCmp) IsEnhancedMode() bool {
	return c.enhancedMode
}

func enhancedNormalPromptFunc(info textarea.PromptInfo) string {
	t := styles.CurrentTheme()
	if info.LineNumber == 0 {
		return "  > "
	}
	// Enhanced styling with better visual design
	if info.Focused {
		return t.S().Base.Foreground(t.GreenDark).Bold(true).Render("::: ")
	}
	return t.S().Muted.Render("::: ")
}

func enhancedYoloPromptFunc(info textarea.PromptInfo) string {
	t := styles.CurrentTheme()
	if info.LineNumber == 0 {
		// Special YoloIcon with Mustard background according to UI/UX specification
		// Enhanced with better visual styling
		if info.Focused {
			return fmt.Sprintf("%s ", t.YoloIconFocused.Bold(true))
		} else {
			return fmt.Sprintf("%s ", t.YoloIconBlurred)
		}
	}
	// Animated dots in Zest color according to UI/UX specification
	// Enhanced with better animation
	if info.Focused {
		return fmt.Sprintf("%s ", t.YoloDotsFocused.Bold(true))
	}
	return fmt.Sprintf("%s ", t.YoloDotsBlurred)
}

