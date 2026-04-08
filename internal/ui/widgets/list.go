package widgets

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ListItem represents a single item in a scrollable list.
type ListItem struct {
	ID       string
	Title    string
	Subtitle string
	Badge    string
	Meta     string
	Online   bool
}

// List is a generic scrollable list widget with vim-style navigation.
type List struct {
	Items        []ListItem
	Cursor       int
	Offset       int
	Width        int
	Height       int
	Focused      bool

	StyleNormal  lipgloss.Style
	StyleActive  lipgloss.Style
	StyleTitle   lipgloss.Style
	StyleSub     lipgloss.Style
	StyleMeta    lipgloss.Style
	StyleBadge   lipgloss.Style
	StyleOnline  lipgloss.Style

	itemHeight   int
}

// NewList creates a new list widget.
func NewList() List {
	return List{
		itemHeight: 2, // title + subtitle
	}
}

// SelectedItem returns the currently selected item, or nil.
func (l *List) SelectedItem() *ListItem {
	if l.Cursor >= 0 && l.Cursor < len(l.Items) {
		return &l.Items[l.Cursor]
	}
	return nil
}

// SelectedID returns the ID of the currently selected item.
func (l *List) SelectedID() string {
	if item := l.SelectedItem(); item != nil {
		return item.ID
	}
	return ""
}

// SetItems replaces all items, keeping cursor in bounds.
func (l *List) SetItems(items []ListItem) {
	l.Items = items
	if l.Cursor >= len(items) {
		l.Cursor = max(0, len(items)-1)
	}
	l.ensureVisible()
}

// Update handles key events for navigation.
func (l *List) Update(msg tea.Msg) (selected bool) {
	if !l.Focused {
		return false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if l.Cursor > 0 {
				l.Cursor--
				l.ensureVisible()
			}
		case "down", "j":
			if l.Cursor < len(l.Items)-1 {
				l.Cursor++
				l.ensureVisible()
			}
		case "home", "g":
			l.Cursor = 0
			l.Offset = 0
		case "end", "G":
			l.Cursor = max(0, len(l.Items)-1)
			l.ensureVisible()
		case "enter":
			return true
		}
	}
	return false
}

func (l *List) ensureVisible() {
	visibleItems := l.Height / l.itemHeight
	if visibleItems <= 0 {
		visibleItems = 1
	}

	if l.Cursor < l.Offset {
		l.Offset = l.Cursor
	}
	if l.Cursor >= l.Offset+visibleItems {
		l.Offset = l.Cursor - visibleItems + 1
	}
}

// View renders the list.
func (l *List) View() string {
	if len(l.Items) == 0 {
		return lipgloss.NewStyle().
			Width(l.Width).
			Height(l.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("#565F89")).
			Render("No items")
	}

	visibleItems := l.Height / l.itemHeight
	if visibleItems <= 0 {
		visibleItems = 1
	}

	var b strings.Builder

	end := min(l.Offset+visibleItems, len(l.Items))
	for i := l.Offset; i < end; i++ {
		item := l.Items[i]
		isActive := i == l.Cursor

		style := l.StyleNormal
		if isActive {
			style = l.StyleActive
		}

		// Title line with meta and badge.
		titleLine := l.StyleTitle.Render(truncate(item.Title, l.Width-10))
		if item.Online {
			titleLine = l.StyleOnline.Render("● ") + titleLine
		}
		if item.Meta != "" {
			metaW := l.Width - lipgloss.Width(titleLine) - 4
			if metaW > 0 {
				meta := l.StyleMeta.Copy().Width(metaW).Align(lipgloss.Right).Render(item.Meta)
				titleLine = titleLine + meta
			}
		}

		// Subtitle line with badge.
		subLine := l.StyleSub.Render(truncate(item.Subtitle, l.Width-8))
		if item.Badge != "" {
			badge := l.StyleBadge.Render(item.Badge)
			padW := l.Width - lipgloss.Width(subLine) - lipgloss.Width(badge) - 4
			if padW > 0 {
				subLine = subLine + strings.Repeat(" ", padW) + badge
			}
		}

		row := style.Width(l.Width).Render(
			fmt.Sprintf("%s\n%s", titleLine, subLine),
		)
		b.WriteString(row)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}
	return string(runes[:maxWidth-3]) + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
