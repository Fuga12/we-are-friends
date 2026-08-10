package media

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

const (
	maxImageDimension = 1920
	jpegQuality       = 82
)

// OptimizeImage normalizes uploaded photos for fast delivery on the website.
// EXIF orientation is applied, oversized images are resized, and the result is
// encoded as a reasonably compact JPEG.
func OptimizeImage(source io.Reader) ([]byte, error) {
	img, err := imaging.Decode(source, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	if width < 1 || height < 1 {
		return nil, fmt.Errorf("image has invalid dimensions")
	}
	if width > maxImageDimension || height > maxImageDimension {
		img = imaging.Fit(img, maxImageDimension, maxImageDimension, imaging.Lanczos)
	}

	// JPEG has no alpha channel. A white background avoids black areas in
	// transparent PNG and WebP uploads.
	background := imaging.New(img.Bounds().Dx(), img.Bounds().Dy(), color.White)
	img = imaging.Paste(background, img, image.Point{})

	var output bytes.Buffer
	if err := imaging.Encode(&output, img, imaging.JPEG, imaging.JPEGQuality(jpegQuality)); err != nil {
		return nil, fmt.Errorf("encode optimized image: %w", err)
	}
	return output.Bytes(), nil
}
