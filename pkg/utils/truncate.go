package utils

// Truncate truncates a string to maxWidth, adding ellipsis if needed.
func Truncate(s string, maxWidth int) string {
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

// TruncateMiddle truncates a string in the middle, preserving start and end.
func TruncateMiddle(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 5 {
		return string(runes[:maxWidth])
	}
	half := (maxWidth - 3) / 2
	return string(runes[:half]) + "..." + string(runes[len(runes)-half:])
}

// PadRight pads a string with spaces to reach the target width.
func PadRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return s
	}
	padding := width - len(runes)
	result := make([]rune, 0, width)
	result = append(result, runes...)
	for i := 0; i < padding; i++ {
		result = append(result, ' ')
	}
	return string(result)
}
