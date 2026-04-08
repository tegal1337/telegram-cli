package utils

import (
	"strings"
	"unicode"
)

// SanitizeForTerminal removes control characters that could interfere
// with terminal rendering, preserving newlines and tabs.
func SanitizeForTerminal(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// StripANSI removes ANSI escape sequences from a string.
func StripANSI(s string) string {
	var b strings.Builder
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		b.WriteByte(s[i])
	}

	return b.String()
}

// EscapeMarkdown escapes special markdown characters.
func EscapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"*", "\\*",
		"_", "\\_",
		"`", "\\`",
		"~", "\\~",
		"[", "\\[",
		"]", "\\]",
	)
	return replacer.Replace(s)
}
