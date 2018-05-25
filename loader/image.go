package loader

import (
	"image"
	"image/draw"
	"io"

	// Supported image formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Image tries to load a file as an image
func readImage(read io.Reader) (*image.RGBA, error) {
	// Decode the image
	img, _, err := image.Decode(read)
	if err != nil {
		return nil, err
	}

	// Check if it's already the correct format
	rgba, ok := img.(*image.RGBA)
	if !ok {
		// Convert to RGBA
		b := img.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	}

	return rgba, nil
}
