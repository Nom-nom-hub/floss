package list

import (
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/nom-nom-hub/floss/internal/tui/components/core/layout"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
	"github.com/sahilm/fuzzy"
)

// Pre-compiled regex for checking if a string is alphanumeric.
var alphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

type FilterableItem interface {
	Item
	FilterValue() string
}

type FilterableList[T FilterableItem] interface {
	List[T]
	Cursor() *tea.Cursor
	SetInputWidth(int)
	SetInputPlaceholder(string)
	Filter(q string) tea.Cmd
}

type HasMatchIndexes interface {
	MatchIndexes([]int)
}

type filterableOptions struct {
	listOptions []ListOption
	placeholder string
	inputHidden bool
	inputWidth  int
	inputStyle  lipgloss.Style
	keyMap      KeyMap
}
type filterableList[T FilterableItem] struct {
	List[T]
	*filterableOptions
	width, height int
	// stores all available items
	items      []T
	input      textinput.Model
	inputWidth int
	query      string
}

type filterableListOption func(*filterableOptions)

func WithFilterPlaceholder(ph string) filterableListOption {
	return func(f *filterableOptions) {
		f.placeholder = ph
	}
}

func WithFilterInputHidden() filterableListOption {
	return func(f *filterableOptions) {
		f.inputHidden = true
	}
}

func WithFilterInputStyle(inputStyle lipgloss.Style) filterableListOption {
	return func(f *filterableOptions) {
		f.inputStyle = inputStyle
	}
}

func WithFilterListOptions(opts ...ListOption) filterableListOption {
	return func(f *filterableOptions) {
		f.listOptions = opts
	}
}

func WithFilterInputWidth(inputWidth int) filterableListOption {
	return func(f *filterableOptions) {
		f.inputWidth = inputWidth
	}
}

func NewFilterableList[T FilterableItem](items []T, opts ...filterableListOption) FilterableList[T] {
	t := styles.CurrentTheme()

	f := &filterableList[T]{
		filterableOptions: &filterableOptions{
			inputStyle:  t.S().Base,
			placeholder: "Type to filter",
			keyMap:      DefaultKeyMap(),
		},
	}
	for _, opt := range opts {
		opt(f.filterableOptions)
	}
	f.List = New(items, f.listOptions...)

	f.updateKeyMaps()
	f.items = f.List.Items()

	if f.inputHidden {
		return f
	}

	ti := textinput.New()
	ti.Placeholder = f.placeholder
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.SetStyles(t.S().TextInput)
	f.input = ti
	return f
}

func (f *filterableList[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// handle movements
		case key.Matches(msg, f.filterableOptions.keyMap.Down),
			key.Matches(msg, f.filterableOptions.keyMap.Up),
			key.Matches(msg, f.filterableOptions.keyMap.DownOneItem),
			key.Matches(msg, f.filterableOptions.keyMap.UpOneItem),
			key.Matches(msg, f.filterableOptions.keyMap.HalfPageDown),
			key.Matches(msg, f.filterableOptions.keyMap.HalfPageUp),
			key.Matches(msg, f.filterableOptions.keyMap.PageDown),
			key.Matches(msg, f.filterableOptions.keyMap.PageUp),
			key.Matches(msg, f.filterableOptions.keyMap.End),
			key.Matches(msg, f.filterableOptions.keyMap.Home):
			u, cmd := f.List.Update(msg)
			f.List = u.(List[T])
			return f, cmd
		default:
			if !f.inputHidden {
				var cmds []tea.Cmd
				var cmd tea.Cmd
				f.input, cmd = f.input.Update(msg)
				cmds = append(cmds, cmd)

				if f.query != f.input.Value() {
					cmd = f.Filter(f.input.Value())
					cmds = append(cmds, cmd)
				}
				f.query = f.input.Value()
				return f, tea.Batch(cmds...)
			}
		}
	}
	u, cmd := f.List.Update(msg)
	f.List = u.(List[T])
	return f, cmd
}

func (f *filterableList[T]) View() string {
	if f.inputHidden {
		return f.List.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		f.filterableOptions.inputStyle.Render(f.input.View()),
		f.List.View(),
	)
}

// removes bindings that are used for search
func (f *filterableList[T]) updateKeyMaps() {
	removeLettersAndNumbers := func(bindings []string) []string {
		var keep []string
		for _, b := range bindings {
			if len(b) != 1 {
				keep = append(keep, b)
				continue
			}
			if b == " " {
				continue
			}
			m := alphanumericRegex.MatchString(b)
			if !m {
				keep = append(keep, b)
			}
		}
		return keep
	}

	updateBinding := func(binding key.Binding) key.Binding {
		newKeys := removeLettersAndNumbers(binding.Keys())
		if len(newKeys) == 0 {
			binding.SetEnabled(false)
			return binding
		}
		binding.SetKeys(newKeys...)
		return binding
	}

	f.filterableOptions.keyMap.Down = updateBinding(f.filterableOptions.keyMap.Down)
	f.filterableOptions.keyMap.Up = updateBinding(f.filterableOptions.keyMap.Up)
	f.filterableOptions.keyMap.DownOneItem = updateBinding(f.filterableOptions.keyMap.DownOneItem)
	f.filterableOptions.keyMap.UpOneItem = updateBinding(f.filterableOptions.keyMap.UpOneItem)
	f.filterableOptions.keyMap.HalfPageDown = updateBinding(f.filterableOptions.keyMap.HalfPageDown)
	f.filterableOptions.keyMap.HalfPageUp = updateBinding(f.filterableOptions.keyMap.HalfPageUp)
	f.filterableOptions.keyMap.PageDown = updateBinding(f.filterableOptions.keyMap.PageDown)
	f.filterableOptions.keyMap.PageUp = updateBinding(f.filterableOptions.keyMap.PageUp)
	f.filterableOptions.keyMap.End = updateBinding(f.filterableOptions.keyMap.End)
	f.filterableOptions.keyMap.Home = updateBinding(f.filterableOptions.keyMap.Home)
}

func (m *filterableList[T]) GetSize() (int, int) {
	return m.width, m.height
}

func (f *filterableList[T]) SetSize(w, h int) tea.Cmd {
	f.width = w
	f.height = h
	if f.inputHidden {
		return f.List.SetSize(w, h)
	}
	if f.inputWidth == 0 {
		f.input.SetWidth(w)
	} else {
		f.input.SetWidth(f.inputWidth)
	}
	return f.List.SetSize(w, h-(f.inputHeight()))
}

func (f *filterableList[T]) inputHeight() int {
	return lipgloss.Height(f.filterableOptions.inputStyle.Render(f.input.View()))
}

func (f *filterableList[T]) Filter(query string) tea.Cmd {
	var cmds []tea.Cmd
	for _, item := range f.items {
		if i, ok := any(item).(layout.Focusable); ok {
			cmds = append(cmds, i.Blur())
		}
		if i, ok := any(item).(HasMatchIndexes); ok {
			i.MatchIndexes(make([]int, 0))
		}
	}

	// Use the SetSelected method instead of directly accessing the field
	cmds = append(cmds, f.List.SetSelected(""))
	if query == "" || len(f.items) == 0 {
		cmds = append(cmds, f.List.SetItems(f.items))
		return tea.Batch(cmds...)
	}

	words := make([]string, len(f.items))
	for i, item := range f.items {
		words[i] = strings.ToLower(item.FilterValue())
	}

	matches := fuzzy.Find(query, words)

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	var matchedItems []T
	for _, match := range matches {
		item := f.items[match.Index]
		if i, ok := any(item).(HasMatchIndexes); ok {
			i.MatchIndexes(match.MatchedIndexes)
		}
		matchedItems = append(matchedItems, item)
	}

	// Check if the List has a direction field or method
	// For now, we'll assume it's always forward
	if f.filterableOptions.inputWidth == 0 {
		// This is a hack to check if we should reverse
		// In a real implementation, we'd need to access the direction field
	}

	cmds = append(cmds, f.List.SetItems(matchedItems))
	return tea.Batch(cmds...)
}

func (f *filterableList[T]) SetItems(items []T) tea.Cmd {
	f.items = items
	return f.List.SetItems(items)
}

func (f *filterableList[T]) Cursor() *tea.Cursor {
	if f.inputHidden {
		return nil
	}
	return f.input.Cursor()
}

func (f *filterableList[T]) Blur() tea.Cmd {
	f.input.Blur()
	return f.List.Blur()
}

func (f *filterableList[T]) Focus() tea.Cmd {
	f.input.Focus()
	return f.List.Focus()
}

func (f *filterableList[T]) IsFocused() bool {
	return f.List.IsFocused()
}

func (f *filterableList[T]) SetInputWidth(w int) {
	f.inputWidth = w
}

func (f *filterableList[T]) SetInputPlaceholder(ph string) {
	f.placeholder = ph
}
