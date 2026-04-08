package media

import (
	"os"
	"strings"
)

// Protocol represents a terminal image rendering protocol.
type Protocol int

const (
	ProtocolBlocks Protocol = iota // Unicode half-block fallback
	ProtocolSixel                  // Sixel graphics
	ProtocolKitty                  // Kitty graphics protocol
)

// DetectProtocol probes the terminal for image rendering capabilities.
func DetectProtocol() Protocol {
	termProgram := os.Getenv("TERM_PROGRAM")

	// Kitty detection
	if strings.Contains(strings.ToLower(termProgram), "kitty") {
		return ProtocolKitty
	}

	// WezTerm supports kitty graphics
	if strings.Contains(strings.ToLower(termProgram), "wezterm") {
		return ProtocolKitty
	}

	// iTerm2 supports sixel
	if strings.Contains(strings.ToLower(termProgram), "iterm") {
		return ProtocolSixel
	}

	term := os.Getenv("TERM")
	// xterm with sixel support
	if strings.Contains(term, "xterm") {
		// Check for sixel support via TERM features
		if os.Getenv("SIXEL_SUPPORT") == "1" {
			return ProtocolSixel
		}
	}

	// mlterm supports sixel
	if strings.Contains(strings.ToLower(termProgram), "mlterm") {
		return ProtocolSixel
	}

	// foot terminal supports sixel
	if strings.Contains(strings.ToLower(termProgram), "foot") {
		return ProtocolSixel
	}

	// Fallback to block characters
	return ProtocolBlocks
}

// ProtocolName returns a human-readable name for the protocol.
func ProtocolName(p Protocol) string {
	switch p {
	case ProtocolKitty:
		return "kitty"
	case ProtocolSixel:
		return "sixel"
	default:
		return "blocks"
	}
}

// ResolveProtocol returns the protocol to use based on config and detection.
func ResolveProtocol(configValue string) Protocol {
	switch strings.ToLower(configValue) {
	case "kitty":
		return ProtocolKitty
	case "sixel":
		return ProtocolSixel
	case "blocks":
		return ProtocolBlocks
	default:
		return DetectProtocol()
	}
}
