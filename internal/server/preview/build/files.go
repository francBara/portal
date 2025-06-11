package build

import (
	"errors"
	"fmt"
	"io"
	"os"
	"portal/internal/server/github"
)

func copyFile(src string, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// seekExtension returns the file extension of the file in path, given some possible extensions.
func seekExtension(path string, extensions []string) (foundExtension string, err error) {
	if fileExists(path) {
		return "", nil
	}

	for _, ext := range extensions {
		newPath := fmt.Sprintf("%s.%s", path, ext)

		if fileExists(newPath) {
			return ext, nil
		}
	}

	return "", errors.New("no file found for " + path)
}

// seekFiles looks for the given filepaths inside cloned project folder, returns the absolute path relative to the first found.
func seekFiles(paths []string) (path string) {
	for _, path := range paths {
		prefixedPath := fmt.Sprintf("%s/%s", github.RepoFolderName, path)

		if fileExists(prefixedPath) {
			return prefixedPath
		}
	}

	return ""
}
