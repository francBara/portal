package parser

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"portal/shared"
	"strings"
)

type ParseOptions struct {
	Verbose bool
}

var acceptedExtensions = [4]string{".js", ".ts", ".jsx", ".tsx"}

// ParseProject crawls a directory and parses all files looking for @portal annotations.
func ParseProject(rootPath string, options ParseOptions) (shared.PortalVariables, error) {
	variables := make(map[string]shared.FileVariables)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		//TODO: Improve rootPath and path parsing
		relativePath := strings.Trim(strings.TrimPrefix(path, strings.Trim(rootPath, "./")), "/")

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

			// Parses the current file and merges it with the total variables
			currentVariables, err := ParseFile(rootPath, relativePath, options)
			if err != nil {
				return err
			}

			if currentVariables.Length() > 0 {
				variables[relativePath] = currentVariables
			}
		}

		return nil
	})

	if err != nil {
		return shared.PortalVariables{}, err
	}
	return variables, nil
}

// ParseFile takes in a file path and outputs the FileVariables relative to all the file @portal annotations.
func ParseFile(basePath string, filePath string, options ParseOptions) (shared.FileVariables, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", basePath, filePath))
	if err != nil {
		return shared.FileVariables{}, err
	}
	defer file.Close()

	var variables shared.FileVariables

	variables.Init()

	scanAll := false
	defaultArguments := portalArguments{
		"group": "Default",
	}

	var currentArguments portalArguments

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	for scanner.Scan() {
		line := scanner.Text()

		// Annotation match
		if annotationMatches := shared.AnnotationRegex.FindStringSubmatch(line); annotationMatches != nil {
			currentArguments = parseAnnotationArguments(annotationMatches[1])

			if options.Verbose {
				fmt.Printf("Annotation: %s\n", line)
			}

			// The "ui" positional argument implies autoscanning of all html trees in the file
			if currentArguments.getString("ui") != "" {
				variables.UI, err = uiVariablesFactory(basePath, filePath, currentArguments)
				if err != nil {
					return shared.FileVariables{}, err
				}

				slog.Info("parsed UI root", "basePath", basePath, "filePath", filePath)
			}

			// The "all" positional argument implies scanning of all subsequent variables, its arguments are applied globally
			if currentArguments.getString("all") != "" {
				scanAll = true
				defaultArguments = currentArguments
			}
			continue
		}

		if scanAll {
			currentArguments = defaultArguments
		}

		if currentArguments != nil {
			// Variable declaration match. Name, type and value are parsed and added to file PortalVariables
			if varMatches := shared.VariableRegex.FindStringSubmatch(line); varMatches != nil {
				varName := varMatches[2]
				value := varMatches[3]

				value = strings.Trim(value, ";")

				if options.Verbose {
					fmt.Printf("Variable: %s\n", line)
				}

				varType := GetVariableType(value)

				if varType == "integer" {
					variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, currentArguments)
					if err != nil {
						return shared.FileVariables{}, err
					}
				} else if varType == "float" {
					variables.Float[varName], err = floatVariableFactory(varName, value, filePath, currentArguments)
					if err != nil {
						return shared.FileVariables{}, err
					}
				} else if varType == "string" {
					variables.String[varName] = stringVariableFactory(varName, value, filePath, currentArguments)
				}
			} else if shared.TailwindRegex.MatchString(line) {
				varName, value := ParseTailwindLine(line)

				if options.Verbose {
					slog.Info(fmt.Sprintf("parsed tailwind line %s", line))
				}

				_, ok := variables.Integer[varName]

				if ok {
					varName += shared.GetRandomString(4)
				}

				variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, currentArguments)
				if err != nil {
					return shared.FileVariables{}, err
				}
			}
			currentArguments = nil
		}
	}

	return variables, nil
}
