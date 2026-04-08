package widgets

import (
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// RenderQRCode generates an ASCII QR code using Unicode block characters.
// Uses U+2580 (upper half block) and U+2584 (lower half block) to pack
// two rows of QR modules into one terminal row.
func RenderQRCode(content string, size int) string {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "Error generating QR code"
	}

	qr.DisableBorder = false
	bitmap := qr.Bitmap()

	rows := len(bitmap)
	cols := 0
	if rows > 0 {
		cols = len(bitmap[0])
	}

	var b strings.Builder

	// Process two rows at a time using half-block characters.
	for y := 0; y < rows; y += 2 {
		for x := 0; x < cols; x++ {
			top := bitmap[y][x]
			bottom := false
			if y+1 < rows {
				bottom = bitmap[y+1][x]
			}

			switch {
			case top && bottom:
				b.WriteString("█")
			case top && !bottom:
				b.WriteString("▀")
			case !top && bottom:
				b.WriteString("▄")
			default:
				b.WriteString(" ")
			}
		}
		if y+2 < rows {
			b.WriteString("\n")
		}
	}

	return b.String()
}
