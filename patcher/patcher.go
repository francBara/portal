package patcher

import (
	"bufio"
	"fmt"
	"portal/parser"
	"strings"
)

func PatchFile(content string, newVariables parser.PortalVariables) string {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var newContent []string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := parser.AnnotationRegex.FindStringSubmatch(line); matches != nil {
			currentType := matches[1]

			newContent = append(newContent, line)

			scanner.Scan()
			line = scanner.Text()

			if matches := parser.VariableRegex.FindStringSubmatch(line); matches != nil {
				declarationType := matches[1]
				varName := matches[2]

				if currentType == "number" {
					newVar, ok := newVariables.Number[varName]
					if !ok {
						fmt.Printf("Variable %s not found in new variables", varName)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s %s = %d;", declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				} else if currentType == "string" {
					newVar, ok := newVariables.String[varName]
					if !ok {
						fmt.Printf("Variable %s not found in new variables", varName)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s %s = \"%s\";", declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				}
			}
		} else {
			newContent = append(newContent, line)
		}
	}

	return strings.Join(newContent, "\n")
}
