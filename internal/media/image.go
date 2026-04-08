package media

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// ImageRenderer renders images to terminal output.
type ImageRenderer struct {
	protocol  Protocol
	maxWidth  int
	maxHeight int
}

// NewImageRenderer creates an image renderer with the given protocol.
func NewImageRenderer(protocol Protocol, maxWidth, maxHeight int) *ImageRenderer {
	return &ImageRenderer{
		protocol:  protocol,
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}

// RenderFile renders an image file to a terminal-displayable string.
func (r *ImageRenderer) RenderFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", fmt.Errorf("decoding image: %w", err)
	}

	return r.RenderImage(img)
}

// RenderImage renders a Go image to a terminal-displayable string.
func (r *ImageRenderer) RenderImage(img image.Image) (string, error) {
	// Resize to fit terminal constraints.
	img = resizeToFit(img, r.maxWidth, r.maxHeight)

	switch r.protocol {
	case ProtocolKitty:
		return renderKitty(img)
	case ProtocolSixel:
		return renderSixel(img)
	default:
		return renderBlocks(img), nil
	}
}

// resizeToFit scales an image to fit within maxWidth x maxHeight cells.
func resizeToFit(img image.Image, maxW, maxH int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Each terminal cell is ~2:1 aspect ratio (taller than wide).
	// For block rendering, each cell = 1 column, half a row.
	targetW := maxW * 2 // pixels per column
	targetH := maxH * 4 // pixels per row (2 rows per char with half-blocks)

	if w <= targetW && h <= targetH {
		return img
	}

	scaleW := float64(targetW) / float64(w)
	scaleH := float64(targetH) / float64(h)
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}

	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	return nearestNeighborResize(img, newW, newH)
}

func nearestNeighborResize(img image.Image, newW, newH int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := x * w / newW + bounds.Min.X
			srcY := y * h / newH + bounds.Min.Y
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}
