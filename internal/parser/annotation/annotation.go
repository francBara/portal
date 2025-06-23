package annotation

import (
	"unicode"
)

type PortalAnnotation struct {
	UI    bool
	All   bool
	Name  string
	Group string
	View  string
	Mock  []string
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

func parseTokens(tokens []string) (ann PortalAnnotation) {
	isMock := false

	for i := 0; i < len(tokens); i++ {
		if isMock {
			ann.Mock = append(ann.Mock, tokens[i])
			continue
		}

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

		if i < len(tokens)-2 && tokens[i+1] == "=" {
			if tokens[i] == "group" {
				ann.Group = tokens[i+2]
			} else if tokens[i] == "view" {
				ann.View = tokens[i+2]
			} else if tokens[i] == "name" {
				ann.Name = tokens[i+2]
			}
			i += 2
		}
	}

	return ann
}

func ParseAnnotation(annotationStr string) PortalAnnotation {
	tokens := tokenizeAnnotation(annotationStr)
	return parseTokens(tokens)
}
