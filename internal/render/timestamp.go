package render

import (
	"fmt"
	"time"
)

// FormatTimestamp formats a Unix timestamp for display.
func FormatTimestamp(ts int32) string {
	t := time.Unix(int64(ts), 0)
	now := time.Now()

	if sameDay(t, now) {
		return t.Format("15:04")
	}

	if sameDay(t, now.AddDate(0, 0, -1)) {
		return "Yesterday " + t.Format("15:04")
	}

	if now.Sub(t) < 7*24*time.Hour {
		return t.Format("Mon 15:04")
	}

	if t.Year() == now.Year() {
		return t.Format("Jan 02")
	}

	return t.Format("2006-01-02")
}

// FormatRelativeTime returns a relative time string (e.g., "2m ago").
func FormatRelativeTime(ts int32) string {
	t := time.Unix(int64(ts), 0)
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return FormatTimestamp(ts)
	}
}

// FormatLastSeen formats a user's last seen timestamp.
func FormatLastSeen(ts int32) string {
	if ts == 0 {
		return "unknown"
	}
	return "last seen " + FormatRelativeTime(ts)
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
