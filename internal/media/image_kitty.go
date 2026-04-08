package media

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"strings"
)

// renderKitty renders an image using the Kitty graphics protocol.
// See: https://sw.kovidgoyal.net/kitty/graphics-protocol/
func renderKitty(img image.Image) (string, error) {
	bounds := img.Bounds()

	// Encode image as PNG.
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("encoding PNG for kitty: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	var b strings.Builder

	// Send image data in chunks of 4096 bytes.
	chunkSize := 4096
	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}

		chunk := encoded[i:end]
		more := 1
		if end >= len(encoded) {
			more = 0
		}

		if i == 0 {
			// First chunk: include image metadata.
			b.WriteString(fmt.Sprintf("\033_Ga=T,f=100,s=%d,v=%d,m=%d;%s\033\\",
				bounds.Dx(), bounds.Dy(), more, chunk))
		} else {
			b.WriteString(fmt.Sprintf("\033_Gm=%d;%s\033\\", more, chunk))
		}
	}

	return b.String(), nil
}
