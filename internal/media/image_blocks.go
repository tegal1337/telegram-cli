package media

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// renderBlocks renders an image using Unicode half-block characters.
// Uses U+2580 (▀ upper half) and U+2584 (▄ lower half) to pack
// two pixel rows into one terminal row.
func renderBlocks(img image.Image) string {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var b strings.Builder

	for y := bounds.Min.Y; y < bounds.Min.Y+h; y += 2 {
		for x := bounds.Min.X; x < bounds.Min.X+w; x++ {
			topColor := img.At(x, y)

			// Bottom pixel (may not exist if odd height).
			var bottomColor color.Color
			if y+1 < bounds.Min.Y+h {
				bottomColor = img.At(x, y+1)
			} else {
				bottomColor = color.Black
			}

			tr, tg, tb, _ := topColor.RGBA()
			br, bg, bb, _ := bottomColor.RGBA()

			// Use upper half block with top color as foreground,
			// bottom color as background.
			b.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀\033[0m",
				tr>>8, tg>>8, tb>>8,
				br>>8, bg>>8, bb>>8,
			))
		}
		if y+2 < bounds.Min.Y+h {
			b.WriteString("\n")
		}
	}

	return b.String()
}
