package utilities

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"strings"
)

// Creates a new png image with data from the canvas and stores it in outFolder.
func SaveImageToFile(file fs.FileInfo, canvas *image.RGBA, outFolder string) {
	nameNoExtension := strings.Split(file.Name(), ".")[0]
	path := fmt.Sprintf("%s/%s-no-bg.%s", outFolder, nameNoExtension, "png")

	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	if encError := png.Encode(outFile, canvas); encError != nil {
		panic(encError)
	}
}

// Transforms an image based on the mode parameter. Pixels with a computed
// luminance value above the given threshold are turned transparent.
func Transform(baseImage *os.File, threshold uint8, mode bool) *image.RGBA {
	imageData, _, err := image.Decode(baseImage)
	if err != nil {
		panic(err)
	}

	// creating a new image so the original one is not modified
	rectangle := imageData.Bounds()    // rectangle with the same size as the original image
	canvas := image.NewRGBA(rectangle) // foundation for the new image

	// applies changes pixel by pixel
	for x := 0; x < rectangle.Dx(); x++ {
		for y := 0; y < rectangle.Dy(); y++ {
			oldPixel := imageData.At(x, y)
			r, g, b, _ := oldPixel.RGBA() // 16 bit RGBA values (0-65535)

			// optimal grayscale conversion weights
			grayscale := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			luminance := uint8(grayscale / 256) // 8 bit luminance values for ease of use

			setAlpha(canvas, threshold, luminance, x, y, oldPixel, mode)
		}
	}

	return canvas
}

// Compares the luminance value of the pixel with the threshold. This can have two outcomes:
//   - luminance < threshold: alpha channel is set to 255, making said pixel fully opaque;
//   - luminance > threshold: alpha is set to 0, thus turning the pixel transparent.
//
// Modes:
//   - "keep" maintains original color of the image.
//   - nil converts to black and white.
func setAlpha(canvas *image.RGBA, threshold uint8, luminance uint8, x int, y int, oldPixel color.Color, mode bool) {
	if luminance > threshold {
		canvas.SetRGBA(x, y, color.RGBA{255, 255, 255, 0})
		return
	}

	if mode {
		canvas.Set(x, y, oldPixel) // original pixel color
		return
	}

	canvas.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
}
