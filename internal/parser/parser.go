package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"portal/internal/parser/annotation"
	"portal/shared"
	"strings"
)

type ParseOptions struct {
	Verbose bool
}

var acceptedExtensions = [4]string{".js", ".ts", ".jsx", ".tsx"}

func ParseAll(repos []string, options ParseOptions) (shared.AllVariables, error) {
	allVars := make(shared.AllVariables)
	var err error

	for _, repo := range repos {
		allVars[repo], err = ParseProject("repos/"+repo, options)
		if err != nil {
			return shared.AllVariables{}, err
		}
	}

	return allVars, nil
}

// ParseProject crawls a directory and parses all files looking for @portal annotations.
func ParseProject(rootPath string, options ParseOptions) (shared.RepoVariables, error) {
	variables := make(shared.RepoVariables)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		//TODO: Improve rootPath and path parsing
		relativePath := strings.Trim(strings.TrimPrefix(path, strings.Trim(rootPath, "./")), "/")

		if err != nil {
			fmt.Printf("Error walking the path %v: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			//TODO: Replace with map lookup and select regex here
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
				return fmt.Errorf("parse file %s: %w", relativePath, err)
			}

			if currentVariables.Length() > 0 {
				variables[relativePath] = currentVariables
			}
		}

		return nil
	})

	if err != nil {
		return shared.RepoVariables{}, err
	}
	return variables, nil
}

// ParseFile takes in a file path and outputs the FileVariables relative to all the file @portal annotations.
func ParseFile(basePath string, filePath string, options ParseOptions) (shared.FileVariables, error) {
	file, err := os.Open(filepath.Join(basePath, filePath))
	if err != nil {
		return shared.FileVariables{}, err
	}
	defer file.Close()

	var variables shared.FileVariables
	variables.Init()

	scanAll := false
	defaultAnn := annotation.PortalAnnotation{
		Group: "Default",
	}

	var ann annotation.PortalAnnotation
	hasAnnotation := false

	lineCounter := 0

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	language, err := GetLanguageRegex(filePath)
	if err != nil {
		return shared.FileVariables{}, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineCounter++

		// Annotation match
		if annotationMatches := shared.AnnotationRegex.FindStringSubmatch(line); annotationMatches != nil {
			ann, err = annotation.ParseAnnotation(annotationMatches[1])
			if err != nil {
				return shared.FileVariables{}, fmt.Errorf("error parsing annotation %s: %w", annotationMatches[1], err)
			}

			if options.Verbose {
				fmt.Printf("Annotation: %s\n", line)
			}

			// The "all" positional argument implies scanning of all subsequent variables, its arguments are applied globally
			if ann.All {
				scanAll = true
				defaultAnn = ann
			}

			hasAnnotation = true

			continue
		}

		if scanAll && !hasAnnotation {
			ann = defaultAnn
		}

		if hasAnnotation || scanAll {
			// Variable declaration match. Name, type and value are parsed and added to file RepoVariables
			if varMatches := shared.VariableRegex.FindStringSubmatch(line); varMatches != nil {
				if options.Verbose {
					fmt.Printf("Variable: %s\n", line)
				}

				varName := varMatches[2]
				value := varMatches[3]

				value = strings.Trim(value, ";")

				varType := GetVariableType(value)

				varId := shared.GetId(varName, lineCounter)

				if varType == "integer" {
					variables.Integer[varId], err = numberVariableFactory(varName, value, filePath, lineCounter, ann)
					if err != nil {
						return shared.FileVariables{}, fmt.Errorf("parsing integer %s: %w", varName, err)
					}
				} else if varType == "float" {
					variables.Float[varId], err = floatVariableFactory(varName, value, filePath, lineCounter, ann)
					if err != nil {
						return shared.FileVariables{}, fmt.Errorf("parsing float %s: %w", varName, err)
					}
				} else if varType == "string" {
					variables.String[varId] = stringVariableFactory(varName, value, filePath, lineCounter, ann)
				}
			}
			hasAnnotation = false
		}
	}

	return variables, nil
}
