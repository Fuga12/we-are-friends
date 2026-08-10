package media

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestOptimizeImageResizesAndEncodesJPEG(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 2400, 1200))
	for y := 0; y < source.Bounds().Dy(); y++ {
		for x := 0; x < source.Bounds().Dx(); x++ {
			source.SetRGBA(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 180, A: 255})
		}
	}
	var input bytes.Buffer
	if err := png.Encode(&input, source); err != nil {
		t.Fatal(err)
	}

	output, err := OptimizeImage(&input)
	if err != nil {
		t.Fatal(err)
	}
	optimized, err := jpeg.Decode(bytes.NewReader(output))
	if err != nil {
		t.Fatalf("result is not JPEG: %v", err)
	}
	if got, want := optimized.Bounds().Dx(), 1920; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}
	if got, want := optimized.Bounds().Dy(), 960; got != want {
		t.Fatalf("height = %d, want %d", got, want)
	}
}

func TestOptimizeImageRejectsInvalidData(t *testing.T) {
	if _, err := OptimizeImage(bytes.NewBufferString("not an image")); err == nil {
		t.Fatal("OptimizeImage() accepted invalid data")
	}
}
