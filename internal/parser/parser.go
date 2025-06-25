package parser

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
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

// ParseProject crawls a directory and parses all files looking for @portal annotations.
func ParseProject(rootPath string, options ParseOptions) (shared.PortalVariables, shared.PortalMocks, error) {
	variables := make(shared.PortalVariables)
	mocks := make(shared.PortalMocks)

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
			currentVariables, currentMocks, err := ParseFile(rootPath, relativePath, options)
			if err != nil {
				return fmt.Errorf("parse file %s: %w", relativePath, err)
			}

			if currentVariables.Length() > 0 {
				variables[relativePath] = currentVariables
			}
			if len(currentMocks) > 0 {
				mocks[relativePath] = currentMocks
			}
		}

		return nil
	})

	if err != nil {
		return shared.PortalVariables{}, shared.PortalMocks{}, err
	}
	return variables, mocks, nil
}

// ParseFile takes in a file path and outputs the FileVariables relative to all the file @portal annotations.
func ParseFile(basePath string, filePath string, options ParseOptions) (shared.FileVariables, shared.FileMocks, error) {
	file, err := os.Open(filepath.Join(basePath, filePath))
	if err != nil {
		return shared.FileVariables{}, shared.FileMocks{}, err
	}
	defer file.Close()

	var variables shared.FileVariables
	variables.Init()

	var mocks shared.FileMocks = make(shared.FileMocks)

	scanAll := false
	defaultAnn := annotation.PortalAnnotation{
		Group: "Default",
	}

	var ann annotation.PortalAnnotation
	hasAnnotation := false

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	for scanner.Scan() {
		line := scanner.Text()

		// Annotation match
		if annotationMatches := shared.AnnotationRegex.FindStringSubmatch(line); annotationMatches != nil {
			ann, err = annotation.ParseAnnotation(annotationMatches[1])
			if err != nil {
				return shared.FileVariables{}, shared.FileMocks{}, fmt.Errorf("error parsing annotation %s: %w", annotationMatches[1], err)
			}

			if options.Verbose {
				fmt.Printf("Annotation: %s\n", line)
				if len(ann.Mocks) > 0 {
					fmt.Printf("Mock: %s\n", ann.Mocks[0])
				}
			}

			if ann.UI {
				// UI variables parsing is outsourced to generateTree tool
				variables.UI, err = uiVariablesFactory(basePath, filePath)
				if err != nil {
					return shared.FileVariables{}, shared.FileMocks{}, fmt.Errorf("parsing ui variables: %w", err)
				}

				slog.Info("parsed UI root", "basePath", basePath, "filePath", filePath)

				continue
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
			// Variable declaration match. Name, type and value are parsed and added to file PortalVariables
			if varMatches := shared.VariableRegex.FindStringSubmatch(line); varMatches != nil {
				if options.Verbose {
					fmt.Printf("Variable: %s\n", line)
				}

				varName := varMatches[2]
				value := varMatches[3]

				if len(ann.Mocks) > 0 {
					mocks[varName] = ann.Mocks[0]
					hasAnnotation = false
					continue
				}

				value = strings.Trim(value, ";")

				varType := GetVariableType(value)

				if varType == "integer" {
					variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, ann)
					if err != nil {
						return shared.FileVariables{}, shared.FileMocks{}, fmt.Errorf("parsing integer %s: %w", varName, err)
					}
				} else if varType == "float" {
					variables.Float[varName], err = floatVariableFactory(varName, value, filePath, ann)
					if err != nil {
						return shared.FileVariables{}, shared.FileMocks{}, fmt.Errorf("parsing float %s: %w", varName, err)
					}
				} else if varType == "string" {
					variables.String[varName] = stringVariableFactory(varName, value, filePath, ann)
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

				variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, ann)
				if err != nil {
					return shared.FileVariables{}, shared.FileMocks{}, err
				}
			}
			hasAnnotation = false
		}
	}

	return variables, mocks, nil
}
