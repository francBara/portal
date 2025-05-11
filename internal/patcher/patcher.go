// Package patcher provides functions to apply portal changes to code.
package patcher

import (
	"bufio"
	"fmt"
	"portal/internal/parser"
	"portal/shared"
	"regexp"
	"strings"
)

func getIndentation(line string) string {
	return regexp.MustCompile(`^\s+`).FindString(line)
}

// PatchFile returns a modified copy of content, where its annotated variables are updated with newVariables values.
func PatchFile(content string, newVariables shared.PortalVariables) (patchedContent string, err error) {
	content, err = patchUI(content, newVariables.UI)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(content))

	var newContent []string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := shared.AnnotationRegex.FindStringSubmatch(line); matches != nil {
			newContent = append(newContent, line)

			scanner.Scan()
			line = scanner.Text()

			if matches := shared.VariableRegex.FindStringSubmatch(line); matches != nil {
				indentation := getIndentation(line)
				declarationType := matches[1]
				varName := matches[2]
				value := matches[3]

				value = strings.Trim(value, ";")

				varType := parser.GetVariableType(value)

				if varType == "integer" {
					newVar, ok := newVariables.Integer[varName]
					if !ok {
						fmt.Printf("Variable %s not found in new variables", varName)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = %d;", indentation, declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				} else if varType == "float" {
					newVar, ok := newVariables.Float[varName]
					if !ok {
						fmt.Printf("Variable %s not found in new variables", varName)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = %f;", indentation, declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				} else if varType == "string" {
					newVar, ok := newVariables.String[varName]
					if !ok {
						fmt.Printf("Variable %s not found in new variables", varName)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = \"%s\";", indentation, declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				}
			} else if shared.TailwindRegex.MatchString(line) {
				varName, _ := parser.ParseTailwindLine(line)

				newVar, ok := newVariables.Integer[varName]
				if !ok {
					fmt.Printf("Tailwind variable %s not found in new variables", varName)
					newContent = append(newContent, line)
					continue
				}

				newLine := UpdateTailwindLine(line, newVar.Value)

				newContent = append(newContent, newLine)
			}
		} else {
			newContent = append(newContent, line)
		}
	}

	return strings.Join(newContent, "\n"), nil
}
