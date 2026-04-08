package media

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// renderSixel renders an image using the Sixel graphics protocol.
// Sixel encodes 6 vertical pixels per character row.
func renderSixel(img image.Image) (string, error) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Build a color palette (up to 256 colors).
	palette := buildPalette(img, 256)

	var b strings.Builder

	// Sixel header: DCS q
	b.WriteString("\033Pq\n")

	// Define palette.
	for i, c := range palette {
		r, g, bl, _ := c.RGBA()
		// Sixel uses 0-100 range for RGB.
		b.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d\n",
			i,
			int(r>>8)*100/255,
			int(g>>8)*100/255,
			int(bl>>8)*100/255,
		))
	}

	// Encode image in 6-pixel-high bands.
	for band := 0; band < h; band += 6 {
		for colorIdx, _ := range palette {
			b.WriteString(fmt.Sprintf("#%d", colorIdx))

			for x := 0; x < w; x++ {
				sixelByte := byte(0)
				for bit := 0; bit < 6; bit++ {
					y := band + bit
					if y >= h {
						break
					}
					pixel := img.At(bounds.Min.X+x, bounds.Min.Y+y)
					if closestColor(pixel, palette) == colorIdx {
						sixelByte |= 1 << uint(bit)
					}
				}
				b.WriteByte(sixelByte + 63) // Sixel data starts at ASCII 63
			}
			b.WriteString("$") // carriage return within band
		}
		b.WriteString("-") // newline between bands
	}

	// Sixel terminator
	b.WriteString("\033\\")

	return b.String(), nil
}

func buildPalette(img image.Image, maxColors int) []color.Color {
	colorSet := make(map[uint32]color.Color)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, bl, _ := c.RGBA()
			// Quantize to reduce color count.
			key := ((r >> 11) << 10) | ((g >> 11) << 5) | (bl >> 11)
			colorSet[key] = c

			if len(colorSet) >= maxColors {
				goto done
			}
		}
	}
done:

	palette := make([]color.Color, 0, len(colorSet))
	for _, c := range colorSet {
		palette = append(palette, c)
	}
	return palette
}

func closestColor(target color.Color, palette []color.Color) int {
	tr, tg, tb, _ := target.RGBA()
	best := 0
	bestDist := uint32(^uint32(0))

	for i, c := range palette {
		cr, cg, cb, _ := c.RGBA()
		dr := int32(tr) - int32(cr)
		dg := int32(tg) - int32(cg)
		db := int32(tb) - int32(cb)
		dist := uint32(dr*dr + dg*dg + db*db)
		if dist < bestDist {
			bestDist = dist
			best = i
		}
	}
	return best
}
