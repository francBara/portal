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
func PatchFile(content string, newVariables shared.FileVariables, language shared.LanguageRegex) (patchedContent string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var newContent []string
	lineCounter := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCounter++

		if matches := shared.AnnotationRegex.FindStringSubmatch(line); matches != nil {
			newContent = append(newContent, line)

			scanner.Scan()
			line = scanner.Text()
			lineCounter++

			//TODO: Does this nested variable lookup account for @portal all annotations?
			if matches := language.Variable.FindStringSubmatch(line); matches != nil {
				indentation := getIndentation(line)
				declarationType := matches[1]
				varName := matches[2]
				value := matches[3]

				value = strings.Trim(value, ";")

				varType := parser.GetVariableType(value)

				varId := shared.GetId(varName, lineCounter)

				if varType == "integer" {
					newVar, ok := newVariables.Integer[varId]
					if !ok {
						//TODO: Should return error
						fmt.Printf("Variable %s at line %d not found in new variables", varName, lineCounter)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = %d;", indentation, declarationType, varName, newVar.Value)

					//TODO: Move append newLine at the end of the loop, once for every if condition
					newContent = append(newContent, newLine)
				} else if varType == "float" {
					newVar, ok := newVariables.Float[varId]
					if !ok {
						fmt.Printf("Variable %s at line %d not found in new variables", varName, lineCounter)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = %f;", indentation, declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				} else if varType == "string" {
					newVar, ok := newVariables.String[varId]
					if !ok {
						fmt.Printf("Variable %s at line %d not found in new variables", varName, lineCounter)
						newContent = append(newContent, line)
						continue
					}

					newLine := fmt.Sprintf("%s%s %s = \"%s\";", indentation, declarationType, varName, newVar.Value)

					newContent = append(newContent, newLine)
				} else {
					newContent = append(newContent, line)
				}
			} else {
				newContent = append(newContent, line)
			}
		} else {
			newContent = append(newContent, line)
		}
	}

	return strings.Join(newContent, "\n"), nil
}
