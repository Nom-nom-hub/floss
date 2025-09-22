package list

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/nom-nom-hub/floss/internal/tui/components/core/layout"
	"github.com/nom-nom-hub/floss/internal/tui/util"
)

type Group[T Item] struct {
	Section ItemSection
	Items   []T
}
type GroupedList[T Item] interface {
	util.Model
	layout.Sizeable
	layout.Focusable

	// Just change state
	MoveUp(int) tea.Cmd
	MoveDown(int) tea.Cmd
	GoToTop() tea.Cmd
	GoToBottom() tea.Cmd
	SelectItemAbove() tea.Cmd
	SelectItemBelow() tea.Cmd
	SetItems([]T) tea.Cmd
	SetSelected(string) tea.Cmd
	SelectedItem() *T
	Items() []T
	UpdateItem(string, T) tea.Cmd
	DeleteItem(string) tea.Cmd
	PrependItem(T) tea.Cmd
	AppendItem(T) tea.Cmd
	StartSelection(col, line int)
	EndSelection(col, line int)
	SelectionStop()
	SelectionClear()
	SelectWord(col, line int)
	SelectParagraph(col, line int)
	GetSelectedText(paddingLeft int) string
	HasSelection() bool

	Groups() []Group[T]
	SetGroups([]Group[T]) tea.Cmd
}
type groupedList[T Item] struct {
	List[T]
	groups []Group[T]
}

func NewGroupedList[T Item](groups []Group[T], opts ...ListOption) GroupedList[T] {
	// Create a temporary list to get the items
	var items []T
	for _, group := range groups {
		items = append(items, group.Section.(T))
		for _, item := range group.Items {
			items = append(items, item)
		}
	}
	
	list := New(items, opts...)
	
	return &groupedList[T]{
		List: list,
	}
}

func (g *groupedList[T]) Init() tea.Cmd {
	g.convertItems()
	return g.List.Init()
}

func (g *groupedList[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	u, cmd := g.List.Update(msg)
	g.List = u.(List[T])
	return g, cmd
}

func (g *groupedList[T]) SelectedItem() *T {
	item := g.List.SelectedItem()
	if item == nil {
		return nil
	}
	dRef := *item
	c, ok := any(dRef).(T)
	if !ok {
		return nil
	}
	return &c
}

func (g *groupedList[T]) convertItems() {
	var items []T
	for _, group := range g.groups {
		// Only add the items, not the section header
		items = append(items, group.Items...)
	}
	g.List.SetItems(items)
}

func (g *groupedList[T]) SetGroups(groups []Group[T]) tea.Cmd {
	g.groups = groups
	g.convertItems()
	return g.List.Init() // Use Init instead of render since render is not in the interface
}

func (g *groupedList[T]) Groups() []Group[T] {
	return g.groups
}

func (g *groupedList[T]) Items() []T {
	return g.List.Items()
}

func (g *groupedList[T]) View() string {
	return g.List.View()
}

func (g *groupedList[T]) Focus() tea.Cmd {
	return g.List.Focus()
}

func (g *groupedList[T]) Blur() tea.Cmd {
	return g.List.Blur()
}

func (g *groupedList[T]) IsFocused() bool {
	return g.List.IsFocused()
}

func (g *groupedList[T]) MoveUp(n int) tea.Cmd {
	return g.List.MoveUp(n)
}

func (g *groupedList[T]) MoveDown(n int) tea.Cmd {
	return g.List.MoveDown(n)
}

func (g *groupedList[T]) GoToTop() tea.Cmd {
	return g.List.GoToTop()
}

func (g *groupedList[T]) GoToBottom() tea.Cmd {
	return g.List.GoToBottom()
}

func (g *groupedList[T]) SelectItemAbove() tea.Cmd {
	return g.List.SelectItemAbove()
}

func (g *groupedList[T]) SelectItemBelow() tea.Cmd {
	return g.List.SelectItemBelow()
}

func (g *groupedList[T]) SetItems(items []T) tea.Cmd {
	return g.List.SetItems(items)
}

func (g *groupedList[T]) SetSelected(id string) tea.Cmd {
	return g.List.SetSelected(id)
}

func (g *groupedList[T]) UpdateItem(id string, item T) tea.Cmd {
	return g.List.UpdateItem(id, item)
}

func (g *groupedList[T]) DeleteItem(id string) tea.Cmd {
	return g.List.DeleteItem(id)
}

func (g *groupedList[T]) PrependItem(item T) tea.Cmd {
	return g.List.PrependItem(item)
}

func (g *groupedList[T]) AppendItem(item T) tea.Cmd {
	return g.List.AppendItem(item)
}

func (g *groupedList[T]) StartSelection(col, line int) {
	g.List.StartSelection(col, line)
}

func (g *groupedList[T]) EndSelection(col, line int) {
	g.List.EndSelection(col, line)
}

func (g *groupedList[T]) SelectionStop() {
	g.List.SelectionStop()
}

func (g *groupedList[T]) SelectionClear() {
	g.List.SelectionClear()
}

func (g *groupedList[T]) SelectWord(col, line int) {
	g.List.SelectWord(col, line)
}

func (g *groupedList[T]) SelectParagraph(col, line int) {
	g.List.SelectParagraph(col, line)
}

func (g *groupedList[T]) GetSelectedText(paddingLeft int) string {
	return g.List.GetSelectedText(paddingLeft)
}

func (g *groupedList[T]) HasSelection() bool {
	return g.List.HasSelection()
}
