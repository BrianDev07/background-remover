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

// Obtains a list of files inside of the given path.
func GetFiles(path string) []fs.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	return entries
}

// Creates a new png image with data from the canvas and stores it in outFolder.
func SaveImageToFile(file fs.DirEntry, canvas *image.RGBA, outFolder string) {
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

// Transforms an image based on the mode parameter. In every case, all pixels with a computed
// luminance value above the given threshold are turned transparent.
func Transform(baseImage *os.File, threshold uint8, mode interface{}) *image.RGBA {
	imageData, _, err := image.Decode(baseImage)
	if err != nil {
		panic(err)
	}

	// creating a new image so the original one is not modified
	rectangle := imageData.Bounds()    // rectangle with the same size as the original image
	canvas := image.NewRGBA(rectangle) // foundation for the new image

	// applies changes pixel by pixel
	for x := 0; x < rectangle.Max.X; x++ {
		for y := 0; y < rectangle.Max.Y; y++ {
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
// when luminance < threshold, the alpha channel is set to 255, making said pixel fully opaque;
// in contrast, when luminance > threshold, alpha is then set to 0, thus turning the pixel transparent.
func setAlpha(canvas *image.RGBA, threshold uint8, luminance uint8, x int, y int, oldPixel color.Color, mode interface{}) {

	if luminance > threshold {
		canvas.SetRGBA(x, y, color.RGBA{uint8(255), uint8(255), uint8(255), uint8(0)})
		return
	}

	if mode == "keep" {
		canvas.Set(x, y, oldPixel) // keeps original color
		return
	}

	r, g, b, _ := uint32(0), uint32(0), uint32(0), uint32(0)

	canvas.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(255)})
}
