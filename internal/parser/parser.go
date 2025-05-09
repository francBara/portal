package parser

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
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

func ParseProject(rootPath string, options ParseOptions) (shared.PortalVariables, error) {
	var variables shared.PortalVariables

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

			currentVariables, err := parseFile(rootPath, strings.TrimPrefix(path, strings.Trim(rootPath, "./")), options)
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
		return shared.PortalVariables{}, err
	}
	return variables, nil
}

func parseFile(basePath string, filePath string, options ParseOptions) (shared.PortalVariables, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", basePath, filePath))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var variables shared.PortalVariables

	variables.Init()
	variables.FileHashes[filePath] = getFileHash(file)

	scanAll := false
	defaultArguments := portalArguments{
		"group": "Default",
	}

	var currentArguments portalArguments

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	for scanner.Scan() {
		line := scanner.Text()

		if annotationMatches := shared.AnnotationRegex.FindStringSubmatch(line); annotationMatches != nil {
			currentArguments = parseAnnotationArguments(annotationMatches[1])

			if options.Verbose {
				fmt.Printf("Annotation: %s\n", line)
			}

			// The "all" positional argument implies scanning of all subsequent variables
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
						return shared.PortalVariables{}, err
					}
				} else if varType == "float" {
					variables.Float[varName], err = floatVariableFactory(varName, value, filePath, currentArguments)
					if err != nil {
						return shared.PortalVariables{}, err
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
					return shared.PortalVariables{}, err
				}
			}
			currentArguments = nil
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
