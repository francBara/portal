package parser

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

	variables.FileHashes = make(map[string]string)
	variables.FileHashes[filePath] = getFileHash(file)

	variables.Integer = make(map[string]shared.IntVariable)
	variables.Float = make(map[string]shared.FloatVariable)
	variables.String = make(map[string]shared.StringVariable)

	scanner := bufio.NewScanner(file)

	file.Seek(0, io.SeekStart)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := shared.AnnotationRegex.FindStringSubmatch(line); matches != nil {
			arguments := parseAnnotationArguments(matches[1])

			if options.Verbose {
				fmt.Printf("Annotation: %s\n", line)
			}

			scanner.Scan()
			line = scanner.Text()

			if matches := shared.VariableRegex.FindStringSubmatch(line); matches != nil {
				varName := matches[2]
				value := matches[3]

				value = strings.Trim(value, ";")

				if options.Verbose {
					fmt.Printf("Variable: %s\n", line)
				}

				varType := GetVariableType(value)

				if varType == "integer" {
					variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, arguments)
					if err != nil {
						return shared.PortalVariables{}, err
					}
				} else if varType == "float" {
					variables.Float[varName], err = floatVariableFactory(varName, value, filePath, arguments)
					if err != nil {
						return shared.PortalVariables{}, err
					}
				} else if varType == "string" {
					variables.String[varName] = stringVariableFactory(varName, value, filePath, arguments)
				}
			} else if shared.TailwindRegex.MatchString(line) {
				varName, value := ParseTailwindLine(line)

				_, ok := variables.Integer[varName]

				if ok {
					varName += shared.GetRandomString(4)
				}

				variables.Integer[varName], err = numberVariableFactory(varName, value, filePath, arguments)
				if err != nil {
					return shared.PortalVariables{}, err
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
