package build

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"portal/internal/server/github"
)

func CopyFile(src string, dst string) error {
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

func CopyDir(src string, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return CopyFile(path, targetPath)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func seekFiles(paths []string) (path string, err error) {
	for _, path := range paths {
		prefixedPath := fmt.Sprintf("%s/%s", github.RepoFolderName, path)

		if fileExists(prefixedPath) {
			return prefixedPath, nil
		}
	}

	return "", errors.New("no file found")
}
