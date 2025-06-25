package annotation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"portal/shared"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type PortalAnnotation struct {
	UI          bool
	All         bool
	DisplayName string
	Group       string
	View        string
	Max         int
	Min         int
	Step        int
	Mocks       []string
	Box         struct {
		Height int
		Width  int
	}
}

// GetPortalVariable generates PortalVariable base data basing on arguments, name and filePath.
func (ann PortalAnnotation) GetPortalVariable(name string, filePath string) shared.PortalVariable {
	group := ann.Group
	if group == "" {
		group = "Default"
	}

	displayName := ann.DisplayName
	if displayName == "" {
		displayName = name
	}

	view := ann.View
	if view == "" {
		view = strings.Split(filepath.Base(filePath), ".")[0]
	}

	return shared.PortalVariable{
		Name:        name,
		DisplayName: displayName,
		View:        view,
		Group:       group,
	}
}

func tokenizeAnnotation(annotationStr string) (tokens []string) {
	currentToken := ""

	inString := false
	escaped := false
	openBrackets := 0

	for _, char := range annotationStr {
		if char == '"' {
			if inString && !escaped {
				inString = false
			} else {
				inString = true
				escaped = false
			}
		} else if inString && char == '\\' {
			escaped = true
			continue
		} else if !inString {
			if char == '{' {
				openBrackets += 1
			} else if char == '}' {
				openBrackets -= 1
			}
		}

		if openBrackets == 0 && !inString {
			if unicode.IsSpace(char) {
				if len(currentToken) > 0 {
					tokens = append(tokens, currentToken)
					currentToken = ""
				}
				continue
			} else if char == '=' {
				if len(currentToken) > 0 {
					tokens = append(tokens, currentToken)
				}
				tokens = append(tokens, "=")
				currentToken = ""
				continue
			}
		}

		currentToken += string(char)
	}

	if len(currentToken) > 0 {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

// TODO: Handle bad = assignments
func parseTokens(tokens []string) (ann PortalAnnotation, err error) {
	isMock := false
	boxRe := regexp.MustCompile(`^(\d+|full)x(\d+|full)$`)

	var savedMocks map[string]any

	for i := 0; i < len(tokens); i++ {
		if isMock {
			mockedValue := tokens[i]

			if !unicode.IsDigit(rune(tokens[i][0])) && tokens[i][0] != '"' && tokens[i][0] != '{' {
				if savedMocks == nil {
					fileContent, err := os.ReadFile("mocks.json")
					if err != nil {
						return PortalAnnotation{}, fmt.Errorf("mocks.json does not exist, trying to find %s: %w", tokens[i], err)
					}
					if err = json.Unmarshal(fileContent, &savedMocks); err != nil {
						return PortalAnnotation{}, err
					}
				}
				jsonMock, ok := savedMocks[tokens[i]]
				if !ok {
					return PortalAnnotation{}, fmt.Errorf("could not find mock %s", tokens[i])
				}

				bytesValue, err := json.Marshal(jsonMock)
				if err != nil {
					return PortalAnnotation{}, err
				}
				mockedValue = string(bytesValue)
			}
			ann.Mocks = append(ann.Mocks, mockedValue)
			continue
		}

		tokens[i] = strings.Trim(tokens[i], "\"")

		if tokens[i] == "all" {
			ann.All = true
			continue
		}
		if tokens[i] == "ui" {
			ann.UI = true
			continue
		}
		if tokens[i] == "mock" {
			isMock = true
			continue
		}

		if boxRe.MatchString(tokens[i]) {
			split := strings.Split(tokens[i], "x")

			if split[0] != "full" {
				ann.Box.Height, _ = strconv.Atoi(split[0])
			}
			if split[1] != "full" {
				ann.Box.Width, _ = strconv.Atoi(split[1])
			}
			continue
		}

		if i < len(tokens)-2 && tokens[i+1] == "=" {
			tokens[i+2] = strings.Trim(tokens[i+2], "\"")

			if tokens[i] == "group" {
				ann.Group = tokens[i+2]
			} else if tokens[i] == "view" {
				ann.View = tokens[i+2]
			} else if tokens[i] == "name" {
				ann.DisplayName = tokens[i+2]
			} else if tokens[i] == "max" {
				ann.Max, err = strconv.Atoi(tokens[i+2])
				if err != nil {
					return PortalAnnotation{}, fmt.Errorf("scanning annotation token max %s: %w", tokens[i+2], err)
				}
			} else if tokens[i] == "min" {
				ann.Min, err = strconv.Atoi(tokens[i+2])
				if err != nil {
					return PortalAnnotation{}, fmt.Errorf("scanning annotation token min %s: %w", tokens[i+2], err)
				}
			}
			i += 2
		}
	}

	return ann, nil
}

func ParseAnnotation(annotationStr string) (ann PortalAnnotation, err error) {
	tokens := tokenizeAnnotation(annotationStr)
	ann, err = parseTokens(tokens)
	if err != nil {
		return PortalAnnotation{}, err
	}

	if ann.All && len(ann.Mocks) > 0 {
		return PortalAnnotation{}, fmt.Errorf("annotation error: cannot have both all and mocks defined")
	}
	return ann, nil
}
