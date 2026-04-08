package widgets

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderBadge renders an unread count badge.
func RenderBadge(count int32, style lipgloss.Style) string {
	if count <= 0 {
		return ""
	}

	text := fmt.Sprintf("%d", count)
	if count > 999 {
		text = "999+"
	}

	return style.Render(text)
}

// RenderOnlineDot renders a small online status indicator.
func RenderOnlineDot(online bool, style lipgloss.Style) string {
	if online {
		return style.Render("●")
	}
	return " "
}
