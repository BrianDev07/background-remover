package main

import (
	utilities "background-remover/pkg"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Not enough arguments: <input_folder> <output_folder> <mode>")
		os.Exit(1)
	}

	const threshold = uint8(128)

	imageFolder := os.Args[1] // location of the original images
	outFolder := os.Args[2]   // output folder for images

	// controls whether the output image keeps original color or not
	mode, err := strconv.ParseBool(os.Args[3])
	if err != nil {
		panic(err)
	}

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

		noBackgroundImage := utilities.Transform(baseImage, threshold, mode)
		utilities.SaveImageToFile(file, noBackgroundImage, outFolder)
	}
}
