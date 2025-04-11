package parser

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ParseOptions struct {
	Verbose bool
}

var acceptedExtensions = [4]string{".js", ".ts", ".jsx", ".tsx"}

func ParseProject(rootPath string, options ParseOptions) PortalVariables {
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

			currentVariables := parseFile(rootPath, strings.TrimPrefix(path, rootPath))

			if currentVariables.HasVariables() {
				variables = variables.Merge(currentVariables)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error during directory walk: %v\n", err)
	}
	return variables
}

func parseFile(basePath string, filePath string) PortalVariables {
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
			currentType := matches[1]

			scanner.Scan()

			line = scanner.Text()

			if matches := VariableRegex.FindStringSubmatch(line); matches != nil {
				varName := matches[2]
				value := matches[3]

				value = strings.Trim(value, "\"'")

				if currentType == "number" {
					parsedValue, err := strconv.Atoi(value)
					if err != nil {
						panic(err)
					}

					variables.Number[varName] = NumberVariable{
						Name:     varName,
						Value:    parsedValue,
						Max:      0,
						Min:      0,
						Step:     1,
						FilePath: filePath,
					}
				} else if currentType == "string" {
					variables.String[varName] = StringVariable{
						Name:     varName,
						Value:    value,
						FilePath: filePath,
					}
				}
			}
		}
	}

	return variables
}

func getFileHash(file *os.File) string {
	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		panic(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
