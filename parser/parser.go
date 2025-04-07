package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ParseOptions struct {
	Verbose bool
}

var annotationRegex = regexp.MustCompile(`//\s*@portal\s+(.*)`)
var variableRegex = regexp.MustCompile(`(let|const|var)\s+(\w+)\s*=\s*(.+?);`)

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

			variables = variables.Concat(parseFile(rootPath, strings.TrimPrefix(path, rootPath)))
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

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for annotation
		if matches := annotationRegex.FindStringSubmatch(line); matches != nil {
			currentType := matches[1]

			scanner.Scan()

			line = scanner.Text()

			if matches := variableRegex.FindStringSubmatch(line); matches != nil {
				varName := matches[2]
				value := matches[3]

				value = strings.Trim(value, "\"'")

				if currentType == "number" {
					parsedValue, err := strconv.Atoi(value)
					if err != nil {
						panic(err)
					}

					variables.Number = append(variables.Number, NumberVariable{
						Name:  varName,
						Value: parsedValue,
						Max:   100,
						Min:   0,
						Step:  1,
					})
				} else if currentType == "string" {
					variables.String = append(variables.String, StringVariable{
						Name:  varName,
						Value: value,
					})
				}
			}
		}
	}

	return variables
}
