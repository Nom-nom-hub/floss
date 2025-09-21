package section

import (
	"image/color"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/nom-nom-hub/floss/internal/tui/exp/diffview"
	"github.com/nom-nom-hub/floss/internal/tui/styles"
)

// Section creates a section header with a title and a line extending to the specified width
func Section(text string, width int) string {
	t := styles.CurrentTheme()
	char := "─"
	length := lipgloss.Width(text) + 1
	remainingWidth := width - length
	lineStyle := t.S().Base.Foreground(t.Border)
	if remainingWidth > 0 {
		text = text + " " + lineStyle.Render(strings.Repeat(char, remainingWidth))
	}
	return text
}

// SectionWithInfo creates a section header with a title, info text, and a line extending to the specified width
func SectionWithInfo(text string, width int, info string) string {
	t := styles.CurrentTheme()
	char := "─"
	length := lipgloss.Width(text) + 1
	remainingWidth := width - length

	if info != "" {
		remainingWidth -= lipgloss.Width(info) + 1 // 1 for the space before info
	}
	lineStyle := t.S().Base.Foreground(t.Border)
	if remainingWidth > 0 {
		text = text + " " + lineStyle.Render(strings.Repeat(char, remainingWidth)) + " " + info
	}
	return text
}

// Title creates a title with a gradient line extending to the specified width
func Title(title string, width int) string {
	t := styles.CurrentTheme()
	char := "╱"
	length := lipgloss.Width(title) + 1
	remainingWidth := width - length
	titleStyle := t.S().Base.Foreground(t.Primary)
	if remainingWidth > 0 {
		lines := strings.Repeat(char, remainingWidth)
		lines = styles.ApplyForegroundGrad(lines, t.Primary, t.Secondary)
		title = titleStyle.Render(title) + " " + lines
	}
	return title
}

type StatusOpts struct {
	Icon             string // if empty no icon will be shown
	Title            string
	TitleColor       color.Color
	Description      string
	DescriptionColor color.Color
	ExtraContent     string // additional content to append after the description
}

func Status(opts StatusOpts, width int) string {
	t := styles.CurrentTheme()
	icon := opts.Icon
	title := opts.Title
	titleColor := t.FgMuted
	if opts.TitleColor != nil {
		titleColor = opts.TitleColor
	}
	description := opts.Description
	descriptionColor := t.FgSubtle
	if opts.DescriptionColor != nil {
		descriptionColor = opts.DescriptionColor
	}
	title = t.S().Base.Foreground(titleColor).Render(title)
	if description != "" {
		extraContentWidth := lipgloss.Width(opts.ExtraContent)
		if extraContentWidth > 0 {
			extraContentWidth += 1
		}
		description = ansi.Truncate(description, width-lipgloss.Width(icon)-lipgloss.Width(title)-2-extraContentWidth, "…")
	}
	description = t.S().Base.Foreground(descriptionColor).Render(description)

	content := []string{}
	if icon != "" {
		content = append(content, icon)
	}
	content = append(content, title, description)
	if opts.ExtraContent != "" {
		content = append(content, opts.ExtraContent)
	}

	return strings.Join(content, " ")
}

func DiffFormatter() *diffview.DiffView {
	t := styles.CurrentTheme()
	formatDiff := diffview.New()
	style := chroma.MustNewStyle("floss", styles.GetChromaTheme())
	diff := formatDiff.ChromaStyle(style).Style(t.S().Diff).TabWidth(4)
	return diff
}