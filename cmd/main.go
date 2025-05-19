package main

import (
	utilities "background-remover/pkg"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/ncruces/zenity"
)

type CustomFile struct {
	cFile fs.FileInfo
	path  string
}

var (
	inputFiles []CustomFile
	mode       bool
)

const threshold = uint8(128)

func init() {

	// Obtain file paths with a selection window
	paths, err := zenity.SelectFileMultiple()
	checkError(err)

	for _, path := range paths {
		fileAsFileInfo, err := os.Stat(path)
		checkError(err)

		customFile := CustomFile{
			cFile: fileAsFileInfo,
			path:  path,
		}

		inputFiles = append(inputFiles, customFile)
	}

	// Question window to change the mode of the output image
	ans := zenity.Question(
		"Keep original color?",
		zenity.Title("Mode"),
		zenity.QuestionIcon,
		zenity.OKLabel("OK"),
		zenity.CancelLabel("No"),
	)

	if ans == nil {
		mode = true
	} else {
		mode = false
	}
}

func main() {
	for _, file := range inputFiles {
		if file.cFile.IsDir() {
			continue
		}

		baseImage, err := os.Open(file.path)
		checkError(err)

		defer baseImage.Close()

		noBackgroundImage, err := utilities.Transform(baseImage, threshold, mode)
		checkError(err)

		utilities.SaveImageToFile(file.cFile, noBackgroundImage, getFileParent(file))
	}
}

// Get parent folder by removing the filename from the absoulte path
func getFileParent(file CustomFile) string {
	return strings.Replace(file.path, file.cFile.Name(), "", 1)
}

// Generate an error window if the error value is not nil
func checkError(e error) {
	if e != nil {
		zenity.Error(
			e.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
			zenity.OKLabel("OK"),
		)

		log.Fatalln(e)
	}
}
