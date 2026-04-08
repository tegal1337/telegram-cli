package media

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

type ImageRenderer struct {
	protocol  Protocol
	maxWidth  int
	maxHeight int
}

func NewImageRenderer(protocol Protocol, maxWidth, maxHeight int) *ImageRenderer {
	return &ImageRenderer{
		protocol:  protocol,
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}

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

func (r *ImageRenderer) RenderImage(img image.Image) (string, error) {
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

// resizeToFit scales image to fit terminal dimensions.
// For blocks: each column = 1 pixel wide, each row = 2 pixels tall (half-blocks).
func resizeToFit(img image.Image, maxCols, maxRows int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Target pixel dimensions
	targetW := maxCols
	targetH := maxRows * 2 // 2 pixels per row with half-blocks

	if srcW <= targetW && srcH <= targetH {
		return img
	}

	scaleW := float64(targetW) / float64(srcW)
	scaleH := float64(targetH) / float64(srcH)
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}

	newW := int(float64(srcW) * scale)
	newH := int(float64(srcH) * scale)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	// Use CatmullRom for high-quality scaling
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}
