package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Extracurricular: change this regex to be useful i.e picture_005, picture_006 etc
var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [()]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	walkDir := "sample" // Extracurricular:  Add flag for path
	var toRename []string

	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}

		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFilename, _ := match(filename)
		newPath := filepath.Join(dir, newFilename)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println("Error Renaming: ", oldPath, newPath, err.Error())
		}
	}
}

func match(filename string) (string, error) {

	if !re.MatchString(filename) {
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}

	return re.ReplaceAllString(filename, replaceString), nil
}
