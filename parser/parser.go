package parser

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ParseOptions struct {
	Verbose bool
}

var acceptedExtensions = [4]string{".js", ".ts", ".jsx", ".tsx"}

func ParseProject(rootPath string, options ParseOptions) (PortalVariables, error) {
	var variables PortalVariables

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error walking the path %v: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			isValidExtension := false
			for _, extension := range acceptedExtensions {
				if strings.HasSuffix(info.Name(), extension) {
					isValidExtension = true
					break
				}
			}
			if !isValidExtension {
				return nil
			}

			if options.Verbose {
				fmt.Printf("Visiting %s\n", path)
			}

			currentVariables, err := parseFile(rootPath, strings.TrimPrefix(path, rootPath))
			if err != nil {
				return err
			}

			if currentVariables.HasVariables() {
				variables = variables.Merge(currentVariables)
			}
		}

		return nil
	})

	if err != nil {
		return PortalVariables{}, err
	}
	return variables, nil
}

func parseFile(basePath string, filePath string) (PortalVariables, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", basePath, filePath))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var variables PortalVariables

	variables.FileHashes = make(map[string]string)
	variables.FileHashes[filePath] = getFileHash(file)

	variables.Number = make(map[string]NumberVariable)
	variables.String = make(map[string]StringVariable)

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := AnnotationRegex.FindStringSubmatch(line); matches != nil {
			arguments := parseAnnotationArguments(matches[1])

			scanner.Scan()
			line = scanner.Text()

			if matches := VariableRegex.FindStringSubmatch(line); matches != nil {
				varName := matches[2]
				value := matches[3]

				varType := getVariableType(value)

				if varType == "number" {
					variables.Number[varName], err = numberVariableFactory(varName, value, filePath, arguments)
					if err != nil {
						return PortalVariables{}, err
					}
				} else if varType == "string" {
					variables.String[varName] = stringVariableFactory(varName, value, arguments)
				}
			}
		}
	}

	return variables, nil
}

func getFileHash(file *os.File) string {
	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		panic(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
