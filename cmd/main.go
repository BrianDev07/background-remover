package main

import (
	utilities "background-remover/pkg"
	"fmt"
	"os"
)

func main() {
	const (
		imageFolder = "./input_images"  // location of original images
		outFolder   = "./output_images" // output folder for images
		threshold   = uint8(128)        // controls black and white levels
	)

	dirFiles := utilities.GetFiles(imageFolder)
	if len(dirFiles) == 0 {
		fmt.Printf("No files found in '%v'.\n", imageFolder)
		return
	}

	for _, file := range dirFiles {
		path := fmt.Sprintf("%s/%s", imageFolder, file.Name())
		baseImage, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		defer baseImage.Close()

		noBackgroundImage := utilities.Transform(baseImage, threshold, nil)
		utilities.SaveImageToFile(file, noBackgroundImage, outFolder)
	}
}
