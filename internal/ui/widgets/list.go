package widgets

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

// ListItem represents a single item in a scrollable list.
type ListItem struct {
	ID       string
	Title    string
	Subtitle string
	Badge    string
	Meta     string
	Online   bool
	Avatar   string // 2-line rendered avatar (half-block image or initials)
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
	avatarW := 5 // avatar column width (4 chars + 1 space)
	textW := l.Width - avatarW

	for i := l.Offset; i < end; i++ {
		item := l.Items[i]
		isActive := i == l.Cursor

		style := l.StyleNormal
		if isActive {
			style = l.StyleActive
		}

		// Avatar: use rendered image or colored initials
		avatar := item.Avatar
		if avatar == "" {
			avatar = renderInitials(item.Title, isActive)
		}

		// Title line with meta
		titleText := item.Title
		if item.Online {
			titleText = "● " + titleText
		}
		titleLine := l.StyleTitle.Render(truncate(titleText, textW-8))
		if item.Meta != "" {
			metaW := textW - lipgloss.Width(titleLine) - 2
			if metaW > 0 {
				meta := l.StyleMeta.Copy().Width(metaW).Align(lipgloss.Right).Render(item.Meta)
				titleLine = titleLine + meta
			}
		}

		// Subtitle line with badge
		subLine := l.StyleSub.Render(truncate(item.Subtitle, textW-6))
		if item.Badge != "" {
			badge := l.StyleBadge.Render(item.Badge)
			padW := textW - lipgloss.Width(subLine) - lipgloss.Width(badge) - 2
			if padW > 0 {
				subLine = subLine + strings.Repeat(" ", padW) + badge
			}
		}

		// Join avatar + text side by side
		textContent := fmt.Sprintf("%s\n%s", titleLine, subLine)

		// Split avatar into lines (should be 2 lines for half-block)
		avatarLines := strings.Split(avatar, "\n")
		textLines := strings.Split(textContent, "\n")

		// Pad to same height
		for len(avatarLines) < 2 {
			avatarLines = append(avatarLines, strings.Repeat(" ", 4))
		}
		for len(textLines) < 2 {
			textLines = append(textLines, "")
		}

		var rowLines []string
		for ri := 0; ri < 2; ri++ {
			av := avatarLines[ri]
			tx := ""
			if ri < len(textLines) {
				tx = textLines[ri]
			}
			rowLines = append(rowLines, av+" "+tx)
		}

		row := style.Width(l.Width).Render(strings.Join(rowLines, "\n"))
		b.WriteString(row)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderInitials creates a 2-line colored box with initials from the title.
func renderInitials(title string, active bool) string {
	// Extract up to 2 initials
	initials := ""
	words := strings.Fields(title)
	for _, w := range words {
		r := []rune(w)
		if len(r) > 0 && r[0] > 32 {
			// Skip emoji-like chars
			if r[0] < 127 || r[0] > 0x2000 {
				initials += string(r[0])
			}
			if len(initials) >= 2 {
				break
			}
		}
	}
	if initials == "" {
		initials = "?"
	}

	// Pick a color based on hash of title
	colors := []string{"196", "208", "220", "34", "39", "129", "170", "214", "49", "201"}
	hash := 0
	for _, r := range title {
		hash = hash*31 + int(r)
	}
	if hash < 0 {
		hash = -hash
	}
	bg := colors[hash%len(colors)]

	style := lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color("231")).
		Bold(true).
		Width(4).
		Align(lipgloss.Center)

	line1 := style.Render(initials)
	line2 := style.Render("  ")

	return line1 + "\n" + line2
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
